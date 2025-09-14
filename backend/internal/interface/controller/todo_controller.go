package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"todo-app/internal/domain"
	"todo-app/internal/interface/middleware"
	"todo-app/internal/usecase"
)

type TodoController struct {
	todoUseCase usecase.TodoUseCase
	validate    *validator.Validate
}

type CreateTodoRequest struct {
	Title    string `json:"title" validate:"required,min=1,max=100"`
	DueDate  string `json:"due_date,omitempty"`
	Priority int    `json:"priority" validate:"min=0,max=2"`
}

type UpdateTodoRequest struct {
	Title       string `json:"title,omitempty" validate:"omitempty,min=1,max=100"`
	DueDate     string `json:"due_date,omitempty"`
	Priority    int    `json:"priority,omitempty" validate:"omitempty,min=0,max=2"`
	IsCompleted bool   `json:"is_completed,omitempty"`
}

type TodoResponse struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Title       string `json:"title"`
	DueDate     string `json:"due_date,omitempty"`
	Priority    int    `json:"priority"`
	IsCompleted bool   `json:"is_completed"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func NewTodoController(todoUseCase usecase.TodoUseCase) *TodoController {
	return &TodoController{
		todoUseCase: todoUseCase,
		validate:    validator.New(),
	}
}

func extractIDFromPath(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "todos" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

func (tc *TodoController) CreateTodo(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		tc.handleErrorResponse(w, domain.ErrUnauthorized)
		return
	}

	var req CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		tc.handleErrorResponse(w, domain.ErrInvalidJSON)
		return
	}

	if err := tc.validate.Struct(req); err != nil {
		validationErr := domain.NewAppError("VALIDATION_FAILED", "バリデーションエラーです: "+err.Error(), http.StatusBadRequest)
		tc.handleErrorResponse(w, validationErr)
		return
	}

	todo := &domain.Todo{
		Title:       req.Title,
		Priority:    req.Priority,
		IsCompleted: false,
	}

	if req.DueDate != "" {
		dueDate, err := time.Parse("2006-01-02", req.DueDate)
		if err != nil {
			dateErr := domain.NewAppError("INVALID_DATE_FORMAT", "日付の形式が正しくありません。YYYY-MM-DD形式で入力してください", http.StatusBadRequest)
			tc.handleErrorResponse(w, dateErr)
			return
		}
		todo.DueDate = &dueDate
	}

	if err := tc.todoUseCase.CreateTodo(r.Context(), userID, todo); err != nil {
		tc.handleErrorResponse(w, err)
		return
	}

	response := tc.todoToResponse(todo)
	tc.writeJSONResponse(w, response, http.StatusCreated)
}

func (tc *TodoController) GetTodos(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		tc.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sortBy := r.URL.Query().Get("sort")

	todos, err := tc.todoUseCase.GetTodos(r.Context(), userID, sortBy)
	if err != nil {
		tc.writeErrorResponse(w, "Failed to get todos", http.StatusInternalServerError)
		return
	}

	responses := make([]TodoResponse, len(todos))
	for i, todo := range todos {
		responses[i] = tc.todoToResponse(todo)
	}

	tc.writeJSONResponse(w, responses, http.StatusOK)
}

func (tc *TodoController) GetTodo(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		tc.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	todoIDStr := extractIDFromPath(r.URL.Path)
	if todoIDStr == "" {
		tc.writeErrorResponse(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}
	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil {
		tc.writeErrorResponse(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	todo, err := tc.todoUseCase.GetTodo(r.Context(), userID, todoID)
	if err != nil {
		tc.writeErrorResponse(w, "Todo not found", http.StatusNotFound)
		return
	}

	response := tc.todoToResponse(todo)
	tc.writeJSONResponse(w, response, http.StatusOK)
}

func (tc *TodoController) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		tc.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	todoIDStr := extractIDFromPath(r.URL.Path)
	if todoIDStr == "" {
		tc.writeErrorResponse(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}
	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil {
		tc.writeErrorResponse(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var req UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		tc.writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := tc.validate.Struct(req); err != nil {
		tc.writeErrorResponse(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	existingTodo, err := tc.todoUseCase.GetTodo(r.Context(), userID, todoID)
	if err != nil {
		tc.writeErrorResponse(w, "Todo not found", http.StatusNotFound)
		return
	}

	if req.Title != "" {
		existingTodo.Title = req.Title
	}
	if req.DueDate != "" {
		dueDate, err := time.Parse("2006-01-02", req.DueDate)
		if err != nil {
			tc.writeErrorResponse(w, "Invalid due_date format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		existingTodo.DueDate = &dueDate
	}
	if req.Priority != 0 {
		existingTodo.Priority = req.Priority
	}
	existingTodo.IsCompleted = req.IsCompleted

	if err := tc.todoUseCase.UpdateTodo(r.Context(), userID, existingTodo); err != nil {
		tc.writeErrorResponse(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}

	response := tc.todoToResponse(existingTodo)
	tc.writeJSONResponse(w, response, http.StatusOK)
}

func (tc *TodoController) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		tc.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	todoIDStr := extractIDFromPath(r.URL.Path)
	if todoIDStr == "" {
		tc.writeErrorResponse(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}
	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil {
		tc.writeErrorResponse(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	if err := tc.todoUseCase.DeleteTodo(r.Context(), userID, todoID); err != nil {
		tc.writeErrorResponse(w, "Failed to delete todo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (tc *TodoController) ToggleTodoComplete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		tc.writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	todoIDStr := extractIDFromPath(r.URL.Path)
	if todoIDStr == "" {
		tc.writeErrorResponse(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}
	todoID, err := strconv.Atoi(todoIDStr)
	if err != nil {
		tc.writeErrorResponse(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	todo, err := tc.todoUseCase.ToggleTodoComplete(r.Context(), userID, todoID)
	if err != nil {
		tc.writeErrorResponse(w, "Failed to toggle todo completion", http.StatusInternalServerError)
		return
	}

	response := tc.todoToResponse(todo)
	tc.writeJSONResponse(w, response, http.StatusOK)
}

func (tc *TodoController) todoToResponse(todo *domain.Todo) TodoResponse {
	response := TodoResponse{
		ID:          todo.ID,
		UserID:      todo.UserID,
		Title:       todo.Title,
		Priority:    todo.Priority,
		IsCompleted: todo.IsCompleted,
		CreatedAt:   todo.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   todo.UpdatedAt.Format(time.RFC3339),
	}

	if todo.DueDate != nil {
		response.DueDate = todo.DueDate.Format("2006-01-02")
	}

	return response
}

func (tc *TodoController) writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (tc *TodoController) writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(ErrorResponse{
		Message: message,
	}); err != nil {
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
	}
}

// handleErrorResponse handles domain errors appropriately
func (tc *TodoController) handleErrorResponse(w http.ResponseWriter, err error) {
	if appErr, ok := domain.IsAppError(err); ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appErr.HTTPCode)

		if encodeErr := json.NewEncoder(w).Encode(appErr); encodeErr != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	// Fallback for non-AppError types
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	fallbackErr := domain.NewAppError("INTERNAL_ERROR", "内部エラーが発生しました", http.StatusInternalServerError)
	if encodeErr := json.NewEncoder(w).Encode(fallbackErr); encodeErr != nil {
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
	}
}
