package domain

import (
	"fmt"
	"net/http"
)

// AppError represents application-specific errors
type AppError struct {
	Code     string                 `json:"code"`
	Message  string                 `json:"message"`
	Details  map[string]interface{} `json:"details,omitempty"`
	HTTPCode int                    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new AppError
func NewAppError(code, message string, httpCode int) *AppError {
	return &AppError{
		Code:     code,
		Message:  message,
		HTTPCode: httpCode,
	}
}

// NewAppErrorWithDetails creates a new AppError with details
func NewAppErrorWithDetails(code, message string, httpCode int, details map[string]interface{}) *AppError {
	return &AppError{
		Code:     code,
		Message:  message,
		Details:  details,
		HTTPCode: httpCode,
	}
}

// User-related errors
var (
	ErrUserNotFound       = NewAppError("USER_NOT_FOUND", "ユーザーが見つかりません", http.StatusNotFound)
	ErrUserAlreadyExists  = NewAppError("USER_ALREADY_EXISTS", "ユーザーが既に存在します", http.StatusConflict)
	ErrUsernameExists     = NewAppError("USERNAME_EXISTS", "このユーザー名は既に使用されています", http.StatusConflict)
	ErrEmailExists        = NewAppError("EMAIL_EXISTS", "このメールアドレスは既に登録されています", http.StatusConflict)
	ErrInvalidCredentials = NewAppError("INVALID_CREDENTIALS", "ユーザー名またはパスワードが正しくありません", http.StatusUnauthorized)
	ErrPasswordHashFailed = NewAppError("PASSWORD_HASH_FAILED", "パスワードの暗号化に失敗しました", http.StatusInternalServerError)
)

// Todo-related errors
var (
	ErrTodoNotFound     = NewAppError("TODO_NOT_FOUND", "Todoが見つかりません", http.StatusNotFound)
	ErrTodoUnauthorized = NewAppError("TODO_UNAUTHORIZED", "このTodoにアクセスする権限がありません", http.StatusForbidden)
)

// Authentication errors
var (
	ErrUnauthorized = NewAppError("UNAUTHORIZED", "認証が必要です", http.StatusUnauthorized)
	ErrTokenInvalid = NewAppError("TOKEN_INVALID", "無効なトークンです", http.StatusUnauthorized)
	ErrTokenExpired = NewAppError("TOKEN_EXPIRED", "トークンの有効期限が切れています", http.StatusUnauthorized)
)

// Validation errors
var (
	ErrValidationFailed = NewAppError("VALIDATION_FAILED", "バリデーションエラーです", http.StatusBadRequest)
	ErrInvalidJSON      = NewAppError("INVALID_JSON", "無効なJSON形式です", http.StatusBadRequest)
)

// Database errors
var (
	ErrDatabaseConnection = NewAppError("DATABASE_CONNECTION", "データベース接続エラーです", http.StatusInternalServerError)
	ErrDatabaseQuery      = NewAppError("DATABASE_QUERY", "データベースクエリエラーです", http.StatusInternalServerError)
)

// NewValidationError creates a validation error with field details
func NewValidationError(fieldErrors map[string]string) *AppError {
	details := make(map[string]interface{})
	for k, v := range fieldErrors {
		details[k] = v
	}

	return NewAppErrorWithDetails(
		"VALIDATION_FAILED",
		"バリデーションエラーです",
		http.StatusBadRequest,
		details,
	)
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) (*AppError, bool) {
	if appErr, ok := err.(*AppError); ok {
		return appErr, true
	}
	return nil, false
}

// WrapError wraps a standard error into an AppError
func WrapError(err error, code, message string, httpCode int) *AppError {
	return &AppError{
		Code:     code,
		Message:  fmt.Sprintf("%s: %v", message, err),
		HTTPCode: httpCode,
	}
}
