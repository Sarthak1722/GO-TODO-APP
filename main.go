package main

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

var todos = []Todo{}

func main() {
	// Initialize a new Fiber app
	app := fiber.New()

	// Define a route for the GET method on the root path '/'
	app.Get("/api/todos/:id", func(c fiber.Ctx) error {
		idStr := c.Params("id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "invalid id",
			})
		}
		for _, todo := range todos {
			if todo.ID == id {
				return c.Status(200).JSON(todo)
			}
		}
		return c.Status(404).JSON(fiber.Map{"msg": "todo not found"})
	})

	app.Post("/api/todos", func(c fiber.Ctx) error {

		todo := &Todo{}

		if err := c.Bind().Body(&todo); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		todos = append(todos, *todo)

		return c.Status(201).JSON(fiber.Map{
			"msg":  "todo created",
			"body": *todo,
		})
	})

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
