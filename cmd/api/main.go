package main

import (
	"log"
	"os"

	"github.com/Sarthak1722/todo_app/internal/handlers"
	"github.com/Sarthak1722/todo_app/internal/logger"
	"github.com/Sarthak1722/todo_app/internal/middleware"
	"github.com/Sarthak1722/todo_app/internal/service"
	"github.com/Sarthak1722/todo_app/internal/store"
	"github.com/Sarthak1722/todo_app/internal/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize logger
	logger.Init("dev")

	// Initialize validator
	validator.Init()

	// Initialize store (data layer)
	todoStore := store.NewInMemoryStore()

	// Initialize service (business logic layer)
	todoService := service.NewTodoService(todoStore)

	// Initialize handlers with injected service
	todoHandler := handlers.NewTodoHandler(todoService)

	// Initialize Fiber app
	app := fiber.New()

	err := godotenv.Load(".env")
	if err!=nil{
		log.Fatal("Error loading .env file")
	}

	PORT:=os.Getenv("PORT")

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
	log.Fatal(app.Listen(":"+PORT))
}
