package persistence

import (
	"context"
	"database/sql"
	"todo-app/internal/domain"
	"todo-app/internal/usecase"
)

type UserPersistence struct {
	queries *Queries
}

func NewUserPersistence(db *sql.DB) usecase.UserRepository {
	return &UserPersistence{
		queries: New(db),
	}
}

func (up *UserPersistence) CreateUser(user *domain.User) error {
	ctx := context.Background()

	params := CreateUserParams{
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}

	sqlcUser, err := up.queries.CreateUser(ctx, params)
	if err != nil {
		return err
	}

	user.ID = int(sqlcUser.ID)
	user.CreatedAt = sqlcUser.CreatedAt.Time
	user.UpdatedAt = sqlcUser.UpdatedAt.Time

	return nil
}

func (up *UserPersistence) GetUserByUsername(username string) (*domain.User, error) {
	ctx := context.Background()

	sqlcUser, err := up.queries.GetUserByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	user := &domain.User{
		ID:           int(sqlcUser.ID),
		Username:     sqlcUser.Username,
		Email:        sqlcUser.Email,
		PasswordHash: sqlcUser.PasswordHash,
		CreatedAt:    sqlcUser.CreatedAt.Time,
		UpdatedAt:    sqlcUser.UpdatedAt.Time,
	}

	return user, nil
}

func (up *UserPersistence) GetUserByEmail(email string) (*domain.User, error) {
	ctx := context.Background()

	sqlcUser, err := up.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	user := &domain.User{
		ID:           int(sqlcUser.ID),
		Username:     sqlcUser.Username,
		Email:        sqlcUser.Email,
		PasswordHash: sqlcUser.PasswordHash,
		CreatedAt:    sqlcUser.CreatedAt.Time,
		UpdatedAt:    sqlcUser.UpdatedAt.Time,
	}

	return user, nil
}

func (up *UserPersistence) GetUserByID(id int) (*domain.User, error) {
	ctx := context.Background()

	sqlcUser, err := up.queries.GetUserByID(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	user := &domain.User{
		ID:           int(sqlcUser.ID),
		Username:     sqlcUser.Username,
		Email:        sqlcUser.Email,
		PasswordHash: sqlcUser.PasswordHash,
		CreatedAt:    sqlcUser.CreatedAt.Time,
		UpdatedAt:    sqlcUser.UpdatedAt.Time,
	}

	return user, nil
}
