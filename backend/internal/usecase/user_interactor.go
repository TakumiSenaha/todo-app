package usecase

import (
	"context"
	"errors"
	"os"
	"time"
	"todo-app/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserInteractor struct {
	UserRepository UserRepository
}

func NewUserInteractor(userRepo UserRepository) *UserInteractor {
	return &UserInteractor{
		UserRepository: userRepo,
	}
}

func (ui *UserInteractor) Register(ctx context.Context, username, email, password string) (*domain.User, error) {
	// Check if user already exists
	existingUser, _ := ui.UserRepository.GetUserByUsername(ctx, username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	existingUser, _ = ui.UserRepository.GetUserByEmail(ctx, email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := &domain.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	err = ui.UserRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ui *UserInteractor) Login(ctx context.Context, username, password string) (string, error) {
	user, err := ui.UserRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := ui.generateJWTToken(user.ID, user.Username)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func (ui *UserInteractor) GetUserByID(ctx context.Context, userID int) (*domain.User, error) {
	return ui.UserRepository.GetUserByID(ctx, userID)
}

func (ui *UserInteractor) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	return ui.UserRepository.GetUserByUsername(ctx, username)
}

func (ui *UserInteractor) UpdateProfile(ctx context.Context, userID int, username, email, currentPassword, newPassword string) (*domain.User, error) {
	// Get current user
	user, err := ui.UserRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Check if username is taken by another user
	if username != user.Username {
		existingUser, _ := ui.UserRepository.GetUserByUsername(ctx, username)
		if existingUser != nil && existingUser.ID != userID {
			return nil, errors.New("username already exists")
		}
	}

	// Check if email is taken by another user
	if email != user.Email {
		existingUser, _ := ui.UserRepository.GetUserByEmail(ctx, email)
		if existingUser != nil && existingUser.ID != userID {
			return nil, errors.New("email already exists")
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
			return nil, errors.New("current password is incorrect")
		}

		// Hash new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("failed to hash new password")
		}
		user.PasswordHash = string(hashedPassword)
	}

	// Save updates
	err = ui.UserRepository.UpdateUser(ctx, user)
	if err != nil {
		return nil, errors.New("failed to update user")
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
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, errors.New("invalid token")
}

func (ui *UserInteractor) Logout(ctx context.Context, tokenString string) error {
	// Simple logout - just return success
	// In a stateless JWT system, logout is handled client-side by removing the token
	return nil
}
