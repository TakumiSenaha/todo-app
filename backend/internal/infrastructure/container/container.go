package container

import (
	"database/sql"
	"todo-app/internal/infrastructure/persistence"
	"todo-app/internal/interface/controller"
	"todo-app/internal/interface/middleware"
	"todo-app/internal/interface/router"
	"todo-app/internal/usecase"
)

// Container manages dependency injection
type Container struct {
	db *sql.DB

	// Infrastructure layer
	queries  *persistence.Queries
	userRepo usecase.UserRepository
	todoRepo usecase.TodoRepository

	// Use case layer
	userInteractor usecase.UserUseCase
	todoInteractor usecase.TodoUseCase

	// Interface layer
	userController *controller.UserController
	todoController *controller.TodoController
	authMiddleware *middleware.AuthMiddleware
	corsMiddleware *middleware.CORSMiddleware
	router         *router.Router
}

// NewContainer creates a new dependency injection container
func NewContainer(db *sql.DB) *Container {
	container := &Container{
		db: db,
	}

	container.buildDependencies()
	return container
}

// buildDependencies constructs all dependencies in the correct order
func (c *Container) buildDependencies() {
	// Infrastructure layer
	c.queries = persistence.New(c.db)
	c.userRepo = persistence.NewUserPersistence(c.db)
	c.todoRepo = persistence.NewTodoRepository(c.queries)

	// Use case layer
	c.userInteractor = usecase.NewUserInteractor(c.userRepo)
	c.todoInteractor = usecase.NewTodoInteractor(c.todoRepo)

	// Interface layer
	c.userController = controller.NewUserController(c.userInteractor)
	c.todoController = controller.NewTodoController(c.todoInteractor)
	c.authMiddleware = middleware.NewAuthMiddleware(c.userInteractor)
	c.corsMiddleware = middleware.NewCORSMiddleware(nil) // Use default config
	c.router = router.NewRouter(c.userController, c.todoController, c.authMiddleware)
}

// GetRouter returns the configured router
func (c *Container) GetRouter() *router.Router {
	return c.router
}

// GetCORSMiddleware returns the CORS middleware
func (c *Container) GetCORSMiddleware() *middleware.CORSMiddleware {
	return c.corsMiddleware
}

// GetDB returns the database connection
func (c *Container) GetDB() *sql.DB {
	return c.db
}
