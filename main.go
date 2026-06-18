package main

import (
	"log"

	"github.com/Sarthak1722/todo_app/database"
	"github.com/Sarthak1722/todo_app/handlers"
	"github.com/gofiber/fiber/v3"
)

func main() {
	// Initialize database
	db := database.NewInMemoryDB()

	// Initialize handlers with injected database
	todoHandler := handlers.NewTodoHandler(db)

	// Initialize a new Fiber app
	app := fiber.New()

	// Define routes with injected handlers
	app.Get("api/todos", todoHandler.GetAllTodos)
	app.Get("/api/todos/:id", todoHandler.GetTodoByID)
	app.Post("/api/todos", todoHandler.CreateTodo)

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
