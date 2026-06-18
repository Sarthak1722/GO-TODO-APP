package handlers

import (
	"strconv"

	"github.com/Sarthak1722/todo_app/database"
	"github.com/Sarthak1722/todo_app/dto"
	"github.com/Sarthak1722/todo_app/models"
	"github.com/gofiber/fiber/v3"
)

// TodoHandler holds dependencies for todo-related handlers
type TodoHandler struct {
	db database.DB
}

// NewTodoHandler creates a new todo handler with injected database
func NewTodoHandler(db database.DB) *TodoHandler {
	return &TodoHandler{db: db}
}

func (h *TodoHandler) CreateTodo(c fiber.Ctx) error {
	todoReq := dto.CreateTodoRequest{}

	if err := c.Bind().Body(&todoReq); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Convert DTO to model
	todo := models.Todo{
		Body:      todoReq.Body,
		Completed: todoReq.Completed,
	}

	// Use injected database
	createdTodo := h.db.CreateTodo(todo)

	return c.Status(201).JSON(fiber.Map{
		"msg":  "todo created",
		"data": createdTodo,
	})
}

func (h *TodoHandler) GetTodoByID(c fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	// Use injected database
	todo := h.db.GetTodoByID(id)
	if todo == nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "todo not found",
		})
	}
	return c.Status(200).JSON(todo)
}


func (h *TodoHandler) GetAllTodos(c fiber.Ctx) error{
	allTodos := h.db.GetAllTodos()
	return c.Status(200).JSON(allTodos)
}