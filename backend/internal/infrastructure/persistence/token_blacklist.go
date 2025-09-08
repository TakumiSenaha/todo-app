package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type TokenBlacklist struct {
	db *sql.DB
}

func NewTokenBlacklist(db *sql.DB) *TokenBlacklist {
	return &TokenBlacklist{db: db}
}

func (tb *TokenBlacklist) AddToken(ctx context.Context, tokenID string, expiresAt time.Time) error {
	query := `
		INSERT INTO token_blacklist (token_id, expires_at) 
		VALUES ($1, $2)
		ON CONFLICT (token_id) DO NOTHING
	`
	_, err := tb.db.ExecContext(ctx, query, tokenID, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}
	return nil
}

func (tb *TokenBlacklist) IsTokenBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM token_blacklist 
			WHERE token_id = $1 AND expires_at > NOW()
		)
	`
	var exists bool
	err := tb.db.QueryRowContext(ctx, query, tokenID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check token blacklist: %w", err)
	}
	return exists, nil
}

func (tb *TokenBlacklist) CleanupExpiredTokens(ctx context.Context) error {
	query := `DELETE FROM token_blacklist WHERE expires_at <= NOW()`
	_, err := tb.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}
	return nil
}
