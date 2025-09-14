package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"todo-app/internal/infrastructure/container"

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

	// Initialize dependency injection container
	appContainer := container.NewContainer(db)

	// Setup routes
	router := appContainer.GetRouter()
	mux := router.SetupRoutes()

	// Apply CORS middleware
	corsMiddleware := appContainer.GetCORSMiddleware()
	handler := corsMiddleware.Handler(mux)

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
