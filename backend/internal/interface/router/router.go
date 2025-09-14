package router

import (
	"net/http"
	"todo-app/internal/interface/controller"
	"todo-app/internal/interface/middleware"
)

// Router represents the application router
type Router struct {
	userController *controller.UserController
	todoController *controller.TodoController
	authMiddleware *middleware.AuthMiddleware
}

// NewRouter creates a new router instance
func NewRouter(
	userController *controller.UserController,
	todoController *controller.TodoController,
	authMiddleware *middleware.AuthMiddleware,
) *Router {
	return &Router{
		userController: userController,
		todoController: todoController,
		authMiddleware: authMiddleware,
	}
}

// SetupRoutes configures all application routes
func (r *Router) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", r.healthCheck)

	// Public endpoints (no authentication required)
	mux.HandleFunc("/api/v1/register", r.userController.Register)
	mux.HandleFunc("/api/v1/login", r.userController.Login)
	mux.HandleFunc("/api/v1/logout", r.userController.Logout)

	// Protected endpoints (authentication required)
	mux.Handle("/api/v1/me", r.authMiddleware.RequireAuth(http.HandlerFunc(r.userController.Me)))
	mux.Handle("/api/v1/profile", r.authMiddleware.RequireAuth(http.HandlerFunc(r.userController.UpdateProfile)))

	// Todo endpoints (authentication required)
	mux.Handle("/api/v1/todos", r.authMiddleware.RequireAuth(http.HandlerFunc(r.handleTodos)))
	mux.Handle("/api/v1/todos/", r.authMiddleware.RequireAuth(http.HandlerFunc(r.handleTodoOperations)))

	return mux
}

// healthCheck handles health check requests
func (r *Router) healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`{"status":"ok","message":"Server is healthy"}`)); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// handleTodos handles /api/v1/todos endpoint
func (r *Router) handleTodos(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.todoController.GetTodos(w, req)
	case http.MethodPost:
		r.todoController.CreateTodo(w, req)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleTodoOperations handles /api/v1/todos/* endpoints
func (r *Router) handleTodoOperations(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	// Handle toggle functionality: /api/v1/todos/{id}/toggle
	if len(path) > 16 && path[len(path)-7:] == "/toggle" && req.Method == http.MethodPatch {
		r.todoController.ToggleTodoComplete(w, req)
		return
	}

	// Handle individual todo operations: /api/v1/todos/{id}
	// Path format: /api/v1/todos/{id}
	pathSegments := len(path)
	if pathSegments >= 16 && path[:14] == "/api/v1/todos/" {
		switch req.Method {
		case http.MethodGet:
			r.todoController.GetTodo(w, req)
		case http.MethodPut:
			r.todoController.UpdateTodo(w, req)
		case http.MethodDelete:
			r.todoController.DeleteTodo(w, req)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	http.Error(w, "Not found", http.StatusNotFound)
}
