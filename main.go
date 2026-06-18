package main

import (
	"log"

	"github.com/Sarthak1722/todo_app/database"
	"github.com/Sarthak1722/todo_app/handlers"
	"github.com/Sarthak1722/todo_app/logger"
	"github.com/Sarthak1722/todo_app/middleware"
	"github.com/Sarthak1722/todo_app/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	// Initialize database
	db := database.NewInMemoryDB()
	validator.Init()
	logger.Init("dev")

	// Initialize handlers with injected database
	todoHandler := handlers.NewTodoHandler(db)

	// Initialize a new Fiber app
	app := fiber.New()
	app.Use(middleware.GetRequestID())
	app.Use(middleware.RequestLogger())
	app.Use(recover.New())
	// Define routes with injected handlers
	app.Get("/health", func(c fiber.Ctx) error {

		return c.JSON(fiber.Map{
			"status": "healthy",
		})
	})
	app.Get("api/todos", todoHandler.GetAllTodos)
	app.Get("/api/todos/:id", todoHandler.GetTodoByID)
	app.Post("/api/todos", todoHandler.CreateTodo)
	app.Delete("/api/todos/:id", todoHandler.DeleteTodoByID)
	app.Patch("/api/todos/:id", todoHandler.PatchTodoByID)

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
