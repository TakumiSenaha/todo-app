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
	err := ti.todoRepo.CreateTodo(ctx, userID, todo)
	if err != nil {
		return domain.WrapError(err, "DATABASE_ERROR", "Todoの作成に失敗しました", 500)
	}
	return nil
}

func (ti *TodoInteractor) GetTodo(ctx context.Context, userID int, todoID int) (*domain.Todo, error) {
	todo, err := ti.todoRepo.GetTodo(ctx, userID, todoID)
	if err != nil {
		return nil, domain.ErrTodoNotFound
	}
	return todo, nil
}

func (ti *TodoInteractor) GetTodos(ctx context.Context, userID int, sortBy string) ([]*domain.Todo, error) {
	todos, err := ti.todoRepo.GetTodos(ctx, userID, sortBy)
	if err != nil {
		return nil, domain.WrapError(err, "DATABASE_ERROR", "Todo一覧の取得に失敗しました", 500)
	}
	return todos, nil
}

func (ti *TodoInteractor) UpdateTodo(ctx context.Context, userID int, todo *domain.Todo) error {
	err := ti.todoRepo.UpdateTodo(ctx, userID, todo)
	if err != nil {
		return domain.WrapError(err, "DATABASE_ERROR", "Todoの更新に失敗しました", 500)
	}
	return nil
}

func (ti *TodoInteractor) DeleteTodo(ctx context.Context, userID int, todoID int) error {
	err := ti.todoRepo.DeleteTodo(ctx, userID, todoID)
	if err != nil {
		return domain.WrapError(err, "DATABASE_ERROR", "Todoの削除に失敗しました", 500)
	}
	return nil
}

func (ti *TodoInteractor) ToggleTodoComplete(ctx context.Context, userID int, todoID int) (*domain.Todo, error) {
	todo, err := ti.todoRepo.ToggleTodoComplete(ctx, userID, todoID)
	if err != nil {
		return nil, domain.WrapError(err, "DATABASE_ERROR", "Todoの状態変更に失敗しました", 500)
	}
	return todo, nil
}
