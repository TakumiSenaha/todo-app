package persistence

import (
	"context"
	"database/sql"
	"time"
	"todo-app/internal/domain"
	"todo-app/internal/usecase"
)

type TodoRepository struct {
	queries *Queries
}

func NewTodoRepository(queries *Queries) usecase.TodoRepository {
	return &TodoRepository{
		queries: queries,
	}
}

func (tr *TodoRepository) CreateTodo(ctx context.Context, userID int, todo *domain.Todo) error {
	params := CreateTodoParams{
		UserID:      int32(userID),
		Title:       todo.Title,
		DueDate:     toSQLNullTime(todo.DueDate),
		Priority:    int32(todo.Priority),
		IsCompleted: todo.IsCompleted,
	}

	sqlcTodo, err := tr.queries.CreateTodo(ctx, params)
	if err != nil {
		return err
	}

	todo.ID = int(sqlcTodo.ID)
	todo.UserID = int(sqlcTodo.UserID)
	todo.CreatedAt = fromSQLNullTime(sqlcTodo.CreatedAt)
	todo.UpdatedAt = fromSQLNullTime(sqlcTodo.UpdatedAt)

	return nil
}

func (tr *TodoRepository) GetTodo(ctx context.Context, userID int, todoID int) (*domain.Todo, error) {
	sqlcTodo, err := tr.queries.GetTodo(ctx, int32(todoID))
	if err != nil {
		return nil, err
	}

	if int(sqlcTodo.UserID) != userID {
		return nil, sql.ErrNoRows
	}

	return &domain.Todo{
		ID:          int(sqlcTodo.ID),
		UserID:      int(sqlcTodo.UserID),
		Title:       sqlcTodo.Title,
		DueDate:     fromSQLNullTimePtr(sqlcTodo.DueDate),
		Priority:    int(sqlcTodo.Priority),
		IsCompleted: sqlcTodo.IsCompleted,
		CreatedAt:   fromSQLNullTime(sqlcTodo.CreatedAt),
		UpdatedAt:   fromSQLNullTime(sqlcTodo.UpdatedAt),
	}, nil
}

func (tr *TodoRepository) GetTodos(ctx context.Context, userID int, sortBy string) ([]*domain.Todo, error) {
	var sqlcTodos []Todo
	var err error

	if sortBy != "" {
		params := ListTodosWithSortParams{
			UserID:  int32(userID),
			Column2: sortBy,
		}
		sqlcTodos, err = tr.queries.ListTodosWithSort(ctx, params)
	} else {
		sqlcTodos, err = tr.queries.ListTodos(ctx, int32(userID))
	}

	if err != nil {
		return nil, err
	}

	todos := make([]*domain.Todo, len(sqlcTodos))
	for i, sqlcTodo := range sqlcTodos {
		todos[i] = &domain.Todo{
			ID:          int(sqlcTodo.ID),
			UserID:      int(sqlcTodo.UserID),
			Title:       sqlcTodo.Title,
			DueDate:     fromSQLNullTimePtr(sqlcTodo.DueDate),
			Priority:    int(sqlcTodo.Priority),
			IsCompleted: sqlcTodo.IsCompleted,
			CreatedAt:   fromSQLNullTime(sqlcTodo.CreatedAt),
			UpdatedAt:   fromSQLNullTime(sqlcTodo.UpdatedAt),
		}
	}

	return todos, nil
}

func (tr *TodoRepository) UpdateTodo(ctx context.Context, userID int, todo *domain.Todo) error {
	params := UpdateTodoParams{
		ID:          int32(todo.ID),
		Title:       todo.Title,
		DueDate:     toSQLNullTime(todo.DueDate),
		Priority:    int32(todo.Priority),
		IsCompleted: todo.IsCompleted,
		UserID:      int32(userID),
	}

	sqlcTodo, err := tr.queries.UpdateTodo(ctx, params)
	if err != nil {
		return err
	}

	todo.UpdatedAt = fromSQLNullTime(sqlcTodo.UpdatedAt)

	return nil
}

func (tr *TodoRepository) DeleteTodo(ctx context.Context, userID int, todoID int) error {
	params := DeleteTodoParams{
		ID:     int32(todoID),
		UserID: int32(userID),
	}

	return tr.queries.DeleteTodo(ctx, params)
}

func (tr *TodoRepository) ToggleTodoComplete(ctx context.Context, userID int, todoID int) (*domain.Todo, error) {
	params := ToggleTodoCompleteParams{
		ID:     int32(todoID),
		UserID: int32(userID),
	}

	sqlcTodo, err := tr.queries.ToggleTodoComplete(ctx, params)
	if err != nil {
		return nil, err
	}

	return &domain.Todo{
		ID:          int(sqlcTodo.ID),
		UserID:      int(sqlcTodo.UserID),
		Title:       sqlcTodo.Title,
		DueDate:     fromSQLNullTimePtr(sqlcTodo.DueDate),
		Priority:    int(sqlcTodo.Priority),
		IsCompleted: sqlcTodo.IsCompleted,
		CreatedAt:   fromSQLNullTime(sqlcTodo.CreatedAt),
		UpdatedAt:   fromSQLNullTime(sqlcTodo.UpdatedAt),
	}, nil
}

func toSQLNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

func fromSQLNullTime(nt sql.NullTime) time.Time {
	if !nt.Valid {
		return time.Time{}
	}
	return nt.Time
}

func fromSQLNullTimePtr(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}
