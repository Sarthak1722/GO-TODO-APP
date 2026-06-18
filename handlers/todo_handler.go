package handlers

import (
	"strconv"

	"github.com/Sarthak1722/todo_app/database"
	"github.com/Sarthak1722/todo_app/dto"
	"github.com/Sarthak1722/todo_app/errors"
	"github.com/Sarthak1722/todo_app/logger"
	"github.com/Sarthak1722/todo_app/models"
	"github.com/Sarthak1722/todo_app/validator"
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
		logger.Log.Error().Err(err).Msg("bad request")
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err := validator.Validate.Struct(todoReq)
	if err != nil {
		logger.Log.Error().Err(err).Msg("invalidate data")
		return c.Status(400).JSON(fiber.Map{
			"errors": errors.FormatValidationErrors(err),
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

func (h *TodoHandler) GetAllTodos(c fiber.Ctx) error {
	allTodos := h.db.GetAllTodos()
	return c.Status(200).JSON(allTodos)
}

func (h *TodoHandler) DeleteTodoByID(c fiber.Ctx) error {
	idStr := c.Params("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	// Use injected database
	done := h.db.DeleteTodoByID(id)
	if done == false {
		return c.Status(404).JSON(fiber.Map{
			"error": "todo not found",
		})
	}
	return c.Status(200).JSON(fiber.Map{"msg": "Todo Deleted"})
}

func (h *TodoHandler) PatchTodoByID(c fiber.Ctx) error {
	idStr := c.Params("id")

	id, errs := strconv.Atoi(idStr)
	if errs != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	todoReq := dto.CreateTodoRequest{}

	if err := c.Bind().Body(&todoReq); err != nil {
		logger.Log.Error().Err(err).Msg("bad request")
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err := validator.Validate.Struct(todoReq)
	if err != nil {
		logger.Log.Error().Err(err).Msg("invalidate data")
		return c.Status(400).JSON(fiber.Map{
			"errors": errors.FormatValidationErrors(err),
		})
	}

	// Convert DTO to model
	todo := models.Todo{
		Body:      todoReq.Body,
		Completed: todoReq.Completed,
	}

	// Use injected database
	updatedTodo := h.db.PatchTodoByID(id, todo)

	if updatedTodo == nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "todo not found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"msg":  "todo updated",
		"data": updatedTodo,
	})

}
