package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"time"
	"todo-app/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type TokenBlacklistRepository interface {
	AddToken(ctx context.Context, tokenID string, expiresAt time.Time) error
	IsTokenBlacklisted(ctx context.Context, tokenID string) (bool, error)
	CleanupExpiredTokens(ctx context.Context) error
}

type RefreshTokenRepository interface {
	StoreRefreshToken(ctx context.Context, tokenID string, userID int, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, tokenID string) (int, bool, error)
	RevokeRefreshToken(ctx context.Context, tokenID string) error
	RevokeAllUserTokens(ctx context.Context, userID int) error
	CleanupExpiredTokens(ctx context.Context) error
}

type UserInteractor struct {
	UserRepository           UserRepository
	TokenBlacklistRepository TokenBlacklistRepository
	RefreshTokenRepository   RefreshTokenRepository
}

func NewUserInteractor(userRepo UserRepository, tokenBlacklistRepo TokenBlacklistRepository, refreshTokenRepo RefreshTokenRepository) *UserInteractor {
	return &UserInteractor{
		UserRepository:           userRepo,
		TokenBlacklistRepository: tokenBlacklistRepo,
		RefreshTokenRepository:   refreshTokenRepo,
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

type LoginTokens struct {
	AccessToken  string
	RefreshToken string
}

func (ui *UserInteractor) Login(ctx context.Context, username, password string) (*LoginTokens, error) {
	user, err := ui.UserRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate access token (short-lived)
	accessToken, err := ui.generateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	// Generate refresh token (long-lived)
	refreshToken, err := ui.generateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &LoginTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
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

func (ui *UserInteractor) generateAccessToken(userID int, username string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
	}

	// Generate unique token ID
	tokenID, err := ui.generateTokenID()
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"jti":      tokenID,
		"user_id":  userID,
		"username": username,
		"type":     "access",
		"exp":      time.Now().Add(15 * time.Minute).Unix(), // Short-lived access token
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func (ui *UserInteractor) generateRefreshToken(ctx context.Context, userID int) (string, error) {
	// Generate unique token ID for refresh token
	tokenID, err := ui.generateTokenID()
	if err != nil {
		return "", err
	}

	// Store in database
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days
	err = ui.RefreshTokenRepository.StoreRefreshToken(ctx, tokenID, userID, expiresAt)
	if err != nil {
		return "", err
	}

	return tokenID, nil
}

func (ui *UserInteractor) generateTokenID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
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
		// Check if token is blacklisted
		if jti, exists := claims["jti"]; exists {
			if jtiStr, ok := jti.(string); ok {
				blacklisted, err := ui.TokenBlacklistRepository.IsTokenBlacklisted(context.Background(), jtiStr)
				if err != nil {
					return nil, errors.New("failed to check token blacklist")
				}
				if blacklisted {
					return nil, errors.New("token has been invalidated")
				}
			}
		}
		return &claims, nil
	}

	return nil, errors.New("invalid token")
}

func (ui *UserInteractor) Logout(ctx context.Context, tokenString string) error {
	// Parse token to get claims
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			jwtSecret = "your-secret-key-change-in-production"
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Get token ID and expiration
		jti, jtiExists := claims["jti"].(string)
		exp, expExists := claims["exp"].(float64)

		if !jtiExists || !expExists {
			return errors.New("invalid token claims")
		}

		expiresAt := time.Unix(int64(exp), 0)

		// Add token to blacklist
		return ui.TokenBlacklistRepository.AddToken(ctx, jti, expiresAt)
	}

	return errors.New("invalid token claims")
}

func (ui *UserInteractor) RefreshToken(ctx context.Context, refreshToken string) (*LoginTokens, error) {
	// Get refresh token from database
	userID, valid, err := ui.RefreshTokenRepository.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("failed to validate refresh token")
	}
	if !valid {
		return nil, errors.New("invalid or expired refresh token")
	}

	// Get user information
	user, err := ui.UserRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Revoke the old refresh token
	err = ui.RefreshTokenRepository.RevokeRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("failed to revoke old refresh token")
	}

	// Generate new access token
	accessToken, err := ui.generateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.New("failed to generate new access token")
	}

	// Generate new refresh token
	newRefreshToken, err := ui.generateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, errors.New("failed to generate new refresh token")
	}

	return &LoginTokens{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
