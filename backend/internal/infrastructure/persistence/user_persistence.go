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

func (up *UserPersistence) CreateUser(ctx context.Context, user *domain.User) error {

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

func (up *UserPersistence) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {

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

func (up *UserPersistence) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {

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

func (up *UserPersistence) GetUserByID(ctx context.Context, id int) (*domain.User, error) {

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

func (up *UserPersistence) UpdateUser(ctx context.Context, user *domain.User) error {

	params := UpdateUserParams{
		ID:           int32(user.ID),
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}

	sqlcUser, err := up.queries.UpdateUser(ctx, params)
	if err != nil {
		return err
	}

	user.UpdatedAt = sqlcUser.UpdatedAt.Time

	return nil
}
