package handlers

import (
	"strconv"

	"github.com/Sarthak1722/todo_app/internal/dto"
	"github.com/Sarthak1722/todo_app/internal/logger"
	"github.com/Sarthak1722/todo_app/internal/service"
	"github.com/Sarthak1722/todo_app/internal/utils"
	"github.com/gofiber/fiber/v3"
)

// TodoHandler holds dependencies for todo-related handlers
type TodoHandler struct {
	service *service.TodoService
}

// NewTodoHandler creates a new todo handler with injected service
func NewTodoHandler(svc *service.TodoService) *TodoHandler {
	return &TodoHandler{service: svc}
}

func (h *TodoHandler) CreateTodo(c fiber.Ctx) error {
	requestID := c.Locals("request_id").(string)
	todoReq := dto.CreateTodoRequest{}

	if err := c.Bind().Body(&todoReq); err != nil {
		logger.Log.Error().
			Str("requestID", requestID).
			Err(err).
			Msg("bad request body")

		return utils.RespondError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	createdTodo, validationErrs := h.service.CreateTodo(c.Context(), todoReq)

	if validationErrs != nil {
		logger.Log.Error().
			Str("requestID", requestID).
			Interface("validation_errors", validationErrs).
			Msg("validation failed")

		return utils.RespondValidationError(c, validationErrs)
	}

	logger.Log.Info().
		Str("requestID", requestID).
		Int("todo_id", createdTodo.ID).
		Msg("todo created successfully")

	return utils.RespondSuccess(c, fiber.StatusCreated, createdTodo, "todo created")
}

func (h *TodoHandler) GetTodoByID(c fiber.Ctx) error {
	requestID := c.Locals("request_id").(string)
	idStr := c.Params("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log.Error().
			Str("requestID", requestID).
			Str("id_param", idStr).
			Msg("invalid todo id")

		return utils.RespondError(c, fiber.StatusBadRequest, "invalid todo id", nil)
	}

	todo, _ := h.service.GetTodoByID(c.Context(), id)

	if todo == nil {
		logger.Log.Warn().
			Str("requestID", requestID).
			Int("todo_id", id).
			Msg("todo not found")

		return utils.RespondError(c, fiber.StatusNotFound, "todo not found", nil)
	}

	return utils.RespondSuccess(c, fiber.StatusOK, todo, "")
}

func (h *TodoHandler) GetAllTodos(c fiber.Ctx) error {
	requestID := c.Locals("request_id").(string)

	allTodos, err := h.service.GetAllTodos(c.Context())
	if err != nil {
		return utils.RespondError(c, fiber.StatusInternalServerError, "could not fetch todos", nil)
	}

	logger.Log.Info().
		Str("requestID", requestID).
		Int("count", len(allTodos)).
		Msg("fetched all todos")

	return utils.RespondSuccess(c, fiber.StatusOK, allTodos, "")
}

func (h *TodoHandler) DeleteTodoByID(c fiber.Ctx) error {
	requestID := c.Locals("request_id").(string)
	idStr := c.Params("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log.Error().
			Str("requestID", requestID).
			Str("id_param", idStr).
			Msg("invalid todo id")

		return utils.RespondError(c, fiber.StatusBadRequest, "invalid todo id", nil)
	}

	err = h.service.DeleteTodo(c.Context(), id)

	if err != nil {
		logger.Log.Warn().
			Str("requestID", requestID).
			Int("todo_id", id).
			Msg("todo not found for deletion")

		return utils.RespondError(c, fiber.StatusNotFound, "todo not found", nil)
	}

	logger.Log.Info().
		Str("requestID", requestID).
		Int("todo_id", id).
		Msg("todo deleted successfully")

	return utils.RespondSuccess(c, fiber.StatusOK, nil, "todo deleted")
}

func (h *TodoHandler) PatchTodoByID(c fiber.Ctx) error {
	requestID := c.Locals("request_id").(string)
	idStr := c.Params("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log.Error().
			Str("requestID", requestID).
			Str("id_param", idStr).
			Msg("invalid todo id")

		return utils.RespondError(c, fiber.StatusBadRequest, "invalid todo id", nil)
	}

	todoReq := dto.PatchTodoRequest{}

	if err := c.Bind().Body(&todoReq); err != nil {
		logger.Log.Error().
			Str("requestID", requestID).
			Err(err).
			Msg("bad request body")

		return utils.RespondError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	updatedTodo, validationErrs := h.service.PatchTodo(c.Context(), id, todoReq)

	if validationErrs != nil {
		logger.Log.Error().
			Str("requestID", requestID).
			Interface("validation_errors", validationErrs).
			Msg("validation failed")

		return utils.RespondValidationError(c, validationErrs)
	}

	if updatedTodo == nil {
		logger.Log.Warn().
			Str("requestID", requestID).
			Int("todo_id", id).
			Msg("todo not found for update")

		return utils.RespondError(c, fiber.StatusNotFound, "todo not found", nil)
	}

	logger.Log.Info().
		Str("requestID", requestID).
		Int("todo_id", id).
		Msg("todo updated successfully")

	return utils.RespondSuccess(c, fiber.StatusOK, updatedTodo, "todo updated")
}