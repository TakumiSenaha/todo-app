package usecase

import (
	"errors"
	"todo-app/internal/domain"
)

type UserInteractor struct {
	UserRepository UserRepository
}

func NewUserInteractor(userRepo UserRepository) *UserInteractor {
	return &UserInteractor{
		UserRepository: userRepo,
	}
}

func (ui *UserInteractor) CreateUser(username, email, passwordHash string) (*domain.User, error) {
	// Check if user already exists
	existingUser, _ := ui.UserRepository.GetUserByUsername(username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	existingUser, _ = ui.UserRepository.GetUserByEmail(email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	user := &domain.User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
	}

	err := ui.UserRepository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ui *UserInteractor) GetUserByUsername(username string) (*domain.User, error) {
	return ui.UserRepository.GetUserByUsername(username)
}