package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"todo-app/internal/infrastructure/persistence"
	"todo-app/internal/interface/controller"
	"todo-app/internal/interface/middleware"
	"todo-app/internal/usecase"

	_ "github.com/lib/pq"
)

func main() {
	// Get database connection string from environment
	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		dbSource = "postgresql://user:password@localhost:5432/todo_db?sslmode=disable"
	}

	// Connect to database
	db, err := sql.Open("postgres", dbSource)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Database connected successfully")

	// Dependency Injection
	// Infrastructure layer
	queries := persistence.New(db)
	userRepo := persistence.NewUserPersistence(db)
	todoRepo := persistence.NewTodoRepository(queries)

	// Use case layer
	userInteractor := usecase.NewUserInteractor(userRepo)
	todoInteractor := usecase.NewTodoInteractor(todoRepo)

	// Interface layer
	userController := controller.NewUserController(userInteractor)
	todoController := controller.NewTodoController(todoInteractor)
	authMiddleware := middleware.NewAuthMiddleware(userInteractor)

	// Setup routes
	mux := http.NewServeMux()

	// CORS middleware
	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cookie")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprint(w, `{"status":"ok","message":"Server is healthy"}`); err != nil {
			log.Printf("Error writing health check response: %v", err)
		}
	})

	// Public endpoints (no authentication required)
	mux.HandleFunc("/api/v1/register", userController.Register)
	mux.HandleFunc("/api/v1/login", userController.Login)
	mux.HandleFunc("/api/v1/logout", userController.Logout)

	// Protected endpoints (authentication required)
	mux.Handle("/api/v1/me", authMiddleware.RequireAuth(http.HandlerFunc(userController.Me)))
	mux.Handle("/api/v1/profile", authMiddleware.RequireAuth(http.HandlerFunc(userController.UpdateProfile)))

	// Todo endpoints (authentication required)
	mux.Handle("/api/v1/todos", authMiddleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			todoController.GetTodos(w, r)
		case http.MethodPost:
			todoController.CreateTodo(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Handle all /api/v1/todos/* paths with custom routing
	mux.Handle("/api/v1/todos/", authMiddleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Handle toggle functionality: /api/v1/todos/{id}/toggle
		if strings.Contains(path, "/toggle") && r.Method == http.MethodPatch {
			todoController.ToggleTodoComplete(w, r)
			return
		}

		// Handle individual todo operations: /api/v1/todos/{id}
		if len(strings.Split(path, "/")) == 5 { // /api/v1/todos/{id}
			switch r.Method {
			case http.MethodGet:
				todoController.GetTodo(w, r)
			case http.MethodPut:
				todoController.UpdateTodo(w, r)
			case http.MethodDelete:
				todoController.DeleteTodo(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		http.Error(w, "Not found", http.StatusNotFound)
	})))

	// Apply CORS middleware
	handler := corsHandler(mux)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
