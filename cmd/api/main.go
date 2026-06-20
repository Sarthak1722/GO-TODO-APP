package main

import (
	"log"
	"os"

	"github.com/Sarthak1722/todo_app/internal/config"
	"github.com/Sarthak1722/todo_app/internal/database"
	"github.com/Sarthak1722/todo_app/internal/handlers"
	"github.com/Sarthak1722/todo_app/internal/logger"
	"github.com/Sarthak1722/todo_app/internal/middleware"
	"github.com/Sarthak1722/todo_app/internal/repository"
	"github.com/Sarthak1722/todo_app/internal/service"
	"github.com/Sarthak1722/todo_app/internal/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize logger
	logger.Init("dev")

	// Initialize validator
	validator.Init()

	// Load configuration
	cfg := config.Load()

	var todoStore repository.Store

	// Initialize store based on config
	if cfg.DBType == "postgres" {
		// Initialize the PostgreSQL database connection pool
		dbPool, err := database.NewPostgresDB(cfg.PostgresDSN)
		if err != nil {
			log.Fatalf("Fatal: Could not connect to the database: %v", err)
		}

		// Ensure the pool is closed gracefully when the app shuts down
		defer dbPool.Close()

		// Initialize PostgreSQL repository
		todoStore = repository.NewPostgresTodoRepository(dbPool)
		log.Println("Using PostgreSQL database")
	} else {
		// Initialize in-memory repository
		todoStore = repository.NewInMemoryTodoRepository()
		log.Println("Using in-memory database")
	}

	// Initialize service (business logic layer)
	todoService := service.NewTodoService(todoStore)

	// Initialize handlers with injected service
	todoHandler := handlers.NewTodoHandler(todoService)

	// Backend initialization complete
	log.Println("Backend initialization complete. Server is ready!")

	// (Your HTTP server startup code like `http.ListenAndServe` will go here later)

	// Initialize Fiber app
	app := fiber.New()

	errs := godotenv.Load(".env")
	if errs != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	// Middleware
	app.Use(middleware.GetRequestID())
	app.Use(middleware.RequestLogger())
	app.Use(middleware.RecoverPanic())

	// Health check endpoint
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
		})
	})

	// Todo routes
	app.Get("/api/todos", todoHandler.GetAllTodos)
	app.Get("/api/todos/:id", todoHandler.GetTodoByID)
	app.Post("/api/todos", todoHandler.CreateTodo)
	app.Delete("/api/todos/:id", todoHandler.DeleteTodoByID)
	app.Patch("/api/todos/:id", todoHandler.PatchTodoByID)

	// Start server
	log.Fatal(app.Listen(":" + PORT))
}
