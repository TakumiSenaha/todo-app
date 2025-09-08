package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
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
	userRepo := persistence.NewUserPersistence(db)

	// Use case layer
	userInteractor := usecase.NewUserInteractor(userRepo)

	// Interface layer
	userController := controller.NewUserController(userInteractor)
	authMiddleware := middleware.NewAuthMiddleware(userInteractor)

	// Setup routes
	mux := http.NewServeMux()

	// CORS middleware
	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

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
