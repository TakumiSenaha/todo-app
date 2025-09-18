package usecase

import (
	"context"
	"os"
	"time"
	"todo-app/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// UserUseCase defines the user use case interface
type UserUseCase interface {
	Register(ctx context.Context, username, email, password string) (*domain.User, error)
	Login(ctx context.Context, username, password string) (string, error)
	GetUserByID(ctx context.Context, userID int) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID int, username, email, currentPassword, newPassword string) (*domain.User, error)
	ValidateJWTToken(tokenString string) (*jwt.MapClaims, error)
	Logout(ctx context.Context, tokenString string) error
}

// UserInteractor implements UserUseCase
type UserInteractor struct {
	UserRepository UserRepository
}

func NewUserInteractor(userRepo UserRepository) UserUseCase {
	return &UserInteractor{
		UserRepository: userRepo,
	}
}

func (ui *UserInteractor) Register(ctx context.Context, username, email, password string) (*domain.User, error) {
	// Check if user already exists
	existingUser, _ := ui.UserRepository.GetUserByUsername(ctx, username)
	if existingUser != nil {
		return nil, domain.ErrUsernameExists
	}

	existingUser, _ = ui.UserRepository.GetUserByEmail(ctx, email)
	if existingUser != nil {
		return nil, domain.ErrEmailExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, domain.ErrPasswordHashFailed
	}

	user := &domain.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	err = ui.UserRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, domain.WrapError(err, "DATABASE_ERROR", "ユーザーの作成に失敗しました", 500)
	}

	return user, nil
}

func (ui *UserInteractor) Login(ctx context.Context, username, password string) (string, error) {
	user, err := ui.UserRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	// Check if user exists
	if user == nil {
		return "", domain.ErrInvalidCredentials
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := ui.generateJWTToken(user.ID, user.Username)
	if err != nil {
		return "", domain.WrapError(err, "TOKEN_GENERATION_FAILED", "トークンの生成に失敗しました", 500)
	}

	return token, nil
}

func (ui *UserInteractor) GetUserByID(ctx context.Context, userID int) (*domain.User, error) {
	user, err := ui.UserRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	
	// Check if user exists
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	
	return user, nil
}

func (ui *UserInteractor) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	return ui.UserRepository.GetUserByUsername(ctx, username)
}

func (ui *UserInteractor) UpdateProfile(ctx context.Context, userID int, username, email, currentPassword, newPassword string) (*domain.User, error) {
	// Get current user
	user, err := ui.UserRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	
	// Check if user exists
	if user == nil {
		return nil, domain.ErrUserNotFound
	}

	// Check if username is taken by another user
	if username != user.Username {
		existingUser, _ := ui.UserRepository.GetUserByUsername(ctx, username)
		if existingUser != nil && existingUser.ID != userID {
			return nil, domain.ErrUsernameExists
		}
	}

	// Check if email is taken by another user
	if email != user.Email {
		existingUser, _ := ui.UserRepository.GetUserByEmail(ctx, email)
		if existingUser != nil && existingUser.ID != userID {
			return nil, domain.ErrEmailExists
		}
	}

	// Update basic info
	user.Username = username
	user.Email = email

	// Update password if provided
	if newPassword != "" {
		// Verify current password
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword))
		if err != nil {
			return nil, domain.ErrInvalidCredentials
		}

		// Hash new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, domain.ErrPasswordHashFailed
		}
		user.PasswordHash = string(hashedPassword)
	}

	// Save updates
	err = ui.UserRepository.UpdateUser(ctx, user)
	if err != nil {
		return nil, domain.WrapError(err, "DATABASE_ERROR", "ユーザーの更新に失敗しました", 500)
	}

	return user, nil
}

func (ui *UserInteractor) generateJWTToken(userID int, username string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
	}

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func (ui *UserInteractor) ValidateJWTToken(tokenString string) (*jwt.MapClaims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrTokenInvalid
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, domain.ErrTokenInvalid
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, domain.ErrTokenInvalid
}

func (ui *UserInteractor) Logout(ctx context.Context, tokenString string) error {
	// Simple logout - just return success
	// In a stateless JWT system, logout is handled client-side by removing the token
	return nil
}
