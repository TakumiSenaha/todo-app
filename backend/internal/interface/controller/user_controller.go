package controller

import (
	"encoding/json"
	"net/http"
	"todo-app/internal/interface/middleware"
	"todo-app/internal/usecase"
)

type UserController struct {
	UserInteractor *usecase.UserInteractor
}

func NewUserController(userInteractor *usecase.UserInteractor) *UserController {
	return &UserController{
		UserInteractor: userInteractor,
	}
}

// 通常のCookie設定を作成
func (uc *UserController) createCookie(name, value string, maxAge int) *http.Cookie {
	// Cookie有効期限の設定（デフォルト24時間）
	if maxAge == 0 {
		maxAge = 24 * 60 * 60 // 24時間
	}

	return &http.Cookie{
		Name:   name,
		Value:  value,
		Path:   "/",
		MaxAge: maxAge,
	}
}

type RegisterUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Message  string `json:"message"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token   string `json:"token"`   // JWTアクセストークン
	User    User   `json:"user"`    // ユーザー情報
	Message string `json:"message"` // レスポンスメッセージ
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UpdateProfileRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	CurrentPassword string `json:"current_password,omitempty"`
	NewPassword     string `json:"new_password,omitempty"`
}

type UpdateProfileResponse struct {
	User    User   `json:"user"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func (uc *UserController) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		uc.writeErrorResponse(w, "Invalid JSON format", nil, http.StatusBadRequest)
		return
	}

	// Validate request
	errors := uc.validateRegisterRequest(req)
	if len(errors) > 0 {
		uc.writeErrorResponse(w, "Validation failed", errors, http.StatusBadRequest)
		return
	}

	user, err := uc.UserInteractor.Register(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		// Handle specific error cases
		switch err.Error() {
		case "username already exists":
			uc.writeErrorResponse(w, "このユーザー名は既に使用されています", map[string]string{"username": "このユーザー名は既に使用されています"}, http.StatusConflict)
			return
		case "email already exists":
			uc.writeErrorResponse(w, "このメールアドレスは既に登録されています", map[string]string{"email": "このメールアドレスは既に登録されています"}, http.StatusConflict)
			return
		default:
			uc.writeErrorResponse(w, err.Error(), nil, http.StatusBadRequest)
			return
		}
	}

	response := RegisterUserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Message:  "User created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		uc.writeErrorResponse(w, "Invalid JSON format", nil, http.StatusBadRequest)
		return
	}

	// Validate request
	errors := uc.validateLoginRequest(req)
	if len(errors) > 0 {
		uc.writeErrorResponse(w, "Validation failed", errors, http.StatusBadRequest)
		return
	}

	// Authenticate user
	token, err := uc.UserInteractor.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		uc.writeErrorResponse(w, "ユーザー名またはパスワードが正しくありません", map[string]string{"credentials": "ユーザー名またはパスワードが正しくありません"}, http.StatusUnauthorized)
		return
	}

	// Get user information
	user, err := uc.UserInteractor.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		uc.writeErrorResponse(w, "ユーザー情報の取得に失敗しました", nil, http.StatusNotFound)
		return
	}

	// Set Cookie
	authCookie := uc.createCookie("auth_token", token, 24*60*60) // 24時間
	http.SetCookie(w, authCookie)

	response := LoginResponse{
		User: User{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		},
		Message: "Login successful",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var token string

	// Try to get token from Authorization header first
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		// Fallback to cookie
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			http.Error(w, `{"error":"Authentication token required"}`, http.StatusBadRequest)
			return
		}
		token = cookie.Value
	}

	// Invalidate the token
	if err := uc.UserInteractor.Logout(r.Context(), token); err != nil {
		http.Error(w, `{"error":"Failed to logout"}`, http.StatusInternalServerError)
		return
	}

	// Cookie削除
	deleteCookie := uc.createCookie("auth_token", "", -1)
	http.SetCookie(w, deleteCookie)

	response := map[string]string{
		"message": "Logout successful",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (uc *UserController) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Get user from database
	user, err := uc.UserInteractor.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (uc *UserController) validateUpdateProfileRequest(req UpdateProfileRequest) map[string]string {
	errors := make(map[string]string)

	// Validate username
	if req.Username == "" {
		errors["username"] = "ユーザー名は必須です"
	} else if len(req.Username) < 3 || len(req.Username) > 20 {
		errors["username"] = "ユーザー名は3-20文字で入力してください"
	} else {
		// Check if username contains only alphanumeric characters and underscores
		for _, char := range req.Username {
			if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
				(char >= '0' && char <= '9') || char == '_') {
				errors["username"] = "ユーザー名は英数字とアンダースコアのみ使用できます"
				break
			}
		}
	}

	// Validate email
	if req.Email == "" {
		errors["email"] = "メールアドレスは必須です"
	} else {
		// Simple email validation
		emailValid := false
		if len(req.Email) > 0 {
			atCount := 0
			dotAfterAt := false
			for i, char := range req.Email {
				if char == '@' {
					atCount++
					if atCount == 1 && i > 0 && i < len(req.Email)-1 {
						// Check for dot after @
						for j := i + 1; j < len(req.Email); j++ {
							if req.Email[j] == '.' && j < len(req.Email)-1 {
								dotAfterAt = true
								break
							}
						}
					}
				}
			}
			emailValid = atCount == 1 && dotAfterAt
		}
		if !emailValid {
			errors["email"] = "有効なメールアドレスを入力してください"
		}
	}

	// Validate password if changing
	if req.NewPassword != "" {
		if req.CurrentPassword == "" {
			errors["current_password"] = "現在のパスワードを入力してください"
		}

		if len(req.NewPassword) < 8 {
			errors["new_password"] = "パスワードは8文字以上で入力してください"
		} else {
			// Check if password contains both letters and numbers
			hasLetter := false
			hasNumber := false
			for _, char := range req.NewPassword {
				if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
					hasLetter = true
				}
				if char >= '0' && char <= '9' {
					hasNumber = true
				}
			}
			if !hasLetter || !hasNumber {
				errors["new_password"] = "パスワードは英数字の両方を含む必要があります"
			}
		}
	}

	return errors
}

func (uc *UserController) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		uc.writeErrorResponse(w, "Invalid JSON", nil, http.StatusBadRequest)
		return
	}

	// Validate request
	errors := uc.validateUpdateProfileRequest(req)
	if len(errors) > 0 {
		uc.writeErrorResponse(w, "Validation failed", errors, http.StatusBadRequest)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Call use case to update profile
	updatedUser, err := uc.UserInteractor.UpdateProfile(r.Context(), userID, req.Username, req.Email, req.CurrentPassword, req.NewPassword)
	if err != nil {
		statusCode := http.StatusBadRequest
		message := err.Error()

		// Handle specific error cases
		switch message {
		case "username already exists":
			uc.writeErrorResponse(w, "このユーザー名は既に使用されています", map[string]string{"username": "このユーザー名は既に使用されています"}, http.StatusConflict)
			return
		case "email already exists":
			uc.writeErrorResponse(w, "このメールアドレスは既に登録されています", map[string]string{"email": "このメールアドレスは既に登録されています"}, http.StatusConflict)
			return
		case "current password is incorrect":
			uc.writeErrorResponse(w, "現在のパスワードが正しくありません", map[string]string{"current_password": "現在のパスワードが正しくありません"}, http.StatusUnauthorized)
			return
		default:
			uc.writeErrorResponse(w, message, nil, statusCode)
			return
		}
	}

	response := UpdateProfileResponse{
		User: User{
			ID:       updatedUser.ID,
			Username: updatedUser.Username,
			Email:    updatedUser.Email,
		},
		Message: "Profile updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (uc *UserController) validateRegisterRequest(req RegisterUserRequest) map[string]string {
	errors := make(map[string]string)

	// Validate username
	if req.Username == "" {
		errors["username"] = "ユーザー名は必須です"
	} else if len(req.Username) < 3 || len(req.Username) > 20 {
		errors["username"] = "ユーザー名は3-20文字で入力してください"
	} else {
		// Check if username contains only alphanumeric characters and underscores
		for _, char := range req.Username {
			if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
				(char >= '0' && char <= '9') || char == '_') {
				errors["username"] = "ユーザー名は英数字とアンダースコアのみ使用できます"
				break
			}
		}
	}

	// Validate email
	if req.Email == "" {
		errors["email"] = "メールアドレスは必須です"
	} else {
		// Simple email validation
		emailValid := false
		if len(req.Email) > 0 {
			atCount := 0
			dotAfterAt := false
			for i, char := range req.Email {
				if char == '@' {
					atCount++
					if atCount == 1 && i > 0 && i < len(req.Email)-1 {
						// Check for dot after @
						for j := i + 1; j < len(req.Email); j++ {
							if req.Email[j] == '.' && j < len(req.Email)-1 {
								dotAfterAt = true
								break
							}
						}
					}
				}
			}
			emailValid = atCount == 1 && dotAfterAt
		}
		if !emailValid {
			errors["email"] = "有効なメールアドレスを入力してください"
		}
	}

	// Validate password
	if req.Password == "" {
		errors["password"] = "パスワードは必須です"
	} else if len(req.Password) < 8 {
		errors["password"] = "パスワードは8文字以上で入力してください"
	} else {
		// Check if password contains both letters and numbers
		hasLetter := false
		hasNumber := false
		for _, char := range req.Password {
			if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
				hasLetter = true
			}
			if char >= '0' && char <= '9' {
				hasNumber = true
			}
		}
		if !hasLetter || !hasNumber {
			errors["password"] = "パスワードは英数字の両方を含む必要があります"
		}
	}

	return errors
}

func (uc *UserController) validateLoginRequest(req LoginRequest) map[string]string {
	errors := make(map[string]string)

	// Validate username
	if req.Username == "" {
		errors["username"] = "ユーザー名は必須です"
	}

	// Validate password
	if req.Password == "" {
		errors["password"] = "パスワードは必須です"
	}

	return errors
}

func (uc *UserController) writeErrorResponse(w http.ResponseWriter, message string, errors map[string]string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := ErrorResponse{
		Message: message,
		Errors:  errors,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
	}
}
