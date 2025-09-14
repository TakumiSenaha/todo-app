package usecase

import (
	"context"
	"todo-app/internal/domain"
)

type TodoUseCase interface {
	CreateTodo(ctx context.Context, userID int, todo *domain.Todo) error
	GetTodo(ctx context.Context, userID int, todoID int) (*domain.Todo, error)
	GetTodos(ctx context.Context, userID int, sortBy string) ([]*domain.Todo, error)
	UpdateTodo(ctx context.Context, userID int, todo *domain.Todo) error
	DeleteTodo(ctx context.Context, userID int, todoID int) error
	ToggleTodoComplete(ctx context.Context, userID int, todoID int) (*domain.Todo, error)
}

type TodoInteractor struct {
	todoRepo TodoRepository
}

func NewTodoInteractor(todoRepo TodoRepository) TodoUseCase {
	return &TodoInteractor{
		todoRepo: todoRepo,
	}
}

func (ti *TodoInteractor) CreateTodo(ctx context.Context, userID int, todo *domain.Todo) error {
	todo.UserID = userID
	return ti.todoRepo.CreateTodo(ctx, userID, todo)
}

func (ti *TodoInteractor) GetTodo(ctx context.Context, userID int, todoID int) (*domain.Todo, error) {
	return ti.todoRepo.GetTodo(ctx, userID, todoID)
}

func (ti *TodoInteractor) GetTodos(ctx context.Context, userID int, sortBy string) ([]*domain.Todo, error) {
	return ti.todoRepo.GetTodos(ctx, userID, sortBy)
}

func (ti *TodoInteractor) UpdateTodo(ctx context.Context, userID int, todo *domain.Todo) error {
	return ti.todoRepo.UpdateTodo(ctx, userID, todo)
}

func (ti *TodoInteractor) DeleteTodo(ctx context.Context, userID int, todoID int) error {
	return ti.todoRepo.DeleteTodo(ctx, userID, todoID)
}

func (ti *TodoInteractor) ToggleTodoComplete(ctx context.Context, userID int, todoID int) (*domain.Todo, error) {
	return ti.todoRepo.ToggleTodoComplete(ctx, userID, todoID)
}
