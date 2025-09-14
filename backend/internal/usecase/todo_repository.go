package usecase

import (
	"context"
	"todo-app/internal/domain"
)

type TodoRepository interface {
	CreateTodo(ctx context.Context, userID int, todo *domain.Todo) error
	GetTodo(ctx context.Context, userID int, todoID int) (*domain.Todo, error)
	GetTodos(ctx context.Context, userID int, sortBy string) ([]*domain.Todo, error)
	UpdateTodo(ctx context.Context, userID int, todo *domain.Todo) error
	DeleteTodo(ctx context.Context, userID int, todoID int) error
	ToggleTodoComplete(ctx context.Context, userID int, todoID int) (*domain.Todo, error)
}
