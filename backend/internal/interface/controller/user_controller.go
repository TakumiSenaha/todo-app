package controller

import (
	"encoding/json"
	"net/http"
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

func (uc *UserController) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// TODO: Hash password properly with bcrypt
	passwordHash := req.Password // This should be hashed in real implementation

	user, err := uc.UserInteractor.CreateUser(req.Username, req.Email, passwordHash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
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
