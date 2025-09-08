package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type RefreshToken struct {
	db *sql.DB
}

func NewRefreshToken(db *sql.DB) *RefreshToken {
	return &RefreshToken{db: db}
}

func (rt *RefreshToken) StoreRefreshToken(ctx context.Context, tokenID string, userID int, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (token_id, user_id, expires_at) 
		VALUES ($1, $2, $3)
	`
	_, err := rt.db.ExecContext(ctx, query, tokenID, userID, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to store refresh token: %w", err)
	}
	return nil
}

func (rt *RefreshToken) GetRefreshToken(ctx context.Context, tokenID string) (int, bool, error) {
	query := `
		SELECT user_id, is_revoked, expires_at 
		FROM refresh_tokens 
		WHERE token_id = $1
	`
	var userID int
	var isRevoked bool
	var expiresAt time.Time

	err := rt.db.QueryRowContext(ctx, query, tokenID).Scan(&userID, &isRevoked, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("failed to get refresh token: %w", err)
	}

	// Check if token is expired or revoked
	if isRevoked || time.Now().After(expiresAt) {
		return 0, false, nil
	}

	return userID, true, nil
}

func (rt *RefreshToken) RevokeRefreshToken(ctx context.Context, tokenID string) error {
	query := `
		UPDATE refresh_tokens 
		SET is_revoked = TRUE 
		WHERE token_id = $1
	`
	_, err := rt.db.ExecContext(ctx, query, tokenID)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}
	return nil
}

func (rt *RefreshToken) RevokeAllUserTokens(ctx context.Context, userID int) error {
	query := `
		UPDATE refresh_tokens 
		SET is_revoked = TRUE 
		WHERE user_id = $1 AND is_revoked = FALSE
	`
	_, err := rt.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke user refresh tokens: %w", err)
	}
	return nil
}

func (rt *RefreshToken) CleanupExpiredTokens(ctx context.Context) error {
	query := `
		DELETE FROM refresh_tokens 
		WHERE expires_at <= NOW() OR is_revoked = TRUE
	`
	_, err := rt.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired refresh tokens: %w", err)
	}
	return nil
}
