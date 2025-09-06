package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"todo-app/internal/infrastructure/persistence"
	"todo-app/internal/interface/controller"
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
	defer db.Close()

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

	// Setup routes
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ok","message":"Server is healthy"}`)
	})

	// User endpoints
	mux.HandleFunc("/api/users/register", userController.RegisterUserHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}