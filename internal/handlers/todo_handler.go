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


// CreateTodo creates a new todo
// @Summary Create a todo
// @Description Adds a new todo to the in-memory/postgres database
// @Tags todos
// @Accept json
// @Produce json
// @Param request body dto.CreateTodoRequest true "Todo Data"
// @Success 201 {object} models.Todo
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /todos [post]
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

	createdTodo, infraErr, validationErrs := h.service.CreateTodo(c.Context(), todoReq)

	// Handle Infrastructure/Database Error FIRST
	if infraErr != nil {
		logger.Log.Error().
			Str("requestID", requestID).
			Err(infraErr).
			Msg("database error")

		return utils.RespondError(c, fiber.StatusInternalServerError, "internal server error", nil)
	}

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



// GetTodoByID retrieves a single todo by its ID
// @Summary Get a todo by ID
// @Description Returns a single todo if it exists in the database
// @Tags todos
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} models.Todo
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /todos/{id} [get]
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

	todo, err := h.service.GetTodoByID(c.Context(), id)
	// Handle Connection Failure FIRST
	if err != nil {
		logger.Log.Error().Str("requestID", requestID).Err(err).Msg("database error")
		return utils.RespondError(c, fiber.StatusInternalServerError, "internal server error", nil)
	}

	if todo == nil {
		logger.Log.Warn().
			Str("requestID", requestID).
			Int("todo_id", id).
			Msg("todo not found")

		return utils.RespondError(c, fiber.StatusNotFound, "todo not found", nil)
	}

	return utils.RespondSuccess(c, fiber.StatusOK, todo, "")
}




// GetAllTodos retrieves all todos
// @Summary Get all todos
// @Description Returns a list of all todos in the database
// @Tags todos
// @Produce json
// @Success 200 {array} models.Todo
// @Failure 500 {object} utils.ErrorResponse
// @Router /todos [get]
func (h *TodoHandler) GetAllTodos(c fiber.Ctx) error {
	requestID := c.Locals("request_id").(string)

	// By passing context(), if Fiber cancels the c.Context(), that cancellation travels through your Service, into your Repository, and immediately kills the pgx query in your Postgres server!

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



// DeleteTodoByID deletes a todo
// @Summary Delete a todo
// @Description Removes a todo from the database
// @Tags todos
// @Param id path int true "Todo ID"
// @Success 204 "No Content"
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /todos/{id} [delete]
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

	err = h.service.DeleteTodoByID(c.Context(), id)

	if err != nil {
		// Specifically check if it was our custom "not found" error
		if err.Error() == "todo not found" {
			logger.Log.Warn().Str("requestID", requestID).Int("todo_id", id).Msg("todo not found for deletion")
			return utils.RespondError(c, fiber.StatusNotFound, "todo not found", nil)
		}
		logger.Log.Error().
			Str("requestID", requestID).
			Err(err).
			Msg("database error during deletion")

		return utils.RespondError(c, fiber.StatusInternalServerError, "internal server error", nil)
	}

	logger.Log.Info().
		Str("requestID", requestID).
		Int("todo_id", id).
		Msg("todo deleted successfully")

	return utils.RespondSuccess(c, fiber.StatusOK, nil, "todo deleted")
}




// PatchTodoByID updates an existing todo
// @Summary Update a todo
// @Description Partially updates an existing todo's title or completion status
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param request body models.Todo true "Update Data"
// @Success 200 {object} models.Todo
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /todos/{id} [patch]
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

	updatedTodo, infraErr, validationErrs := h.service.PatchTodoByID(c.Context(), id, todoReq)

	// Handle Infrastructure/Database Error FIRST
	if infraErr != nil {
		logger.Log.Error().
			Str("requestID", requestID).
			Err(infraErr).
			Msg("database error")

		return utils.RespondError(c, fiber.StatusInternalServerError, "internal server error", nil)
	}

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
