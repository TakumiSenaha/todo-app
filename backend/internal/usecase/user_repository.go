package usecase

import (
	"todo-app/internal/domain"
)

type UserRepository interface {
	CreateUser(user *domain.User) error
	GetUserByUsername(username string) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	GetUserByID(id int) (*domain.User, error)
}