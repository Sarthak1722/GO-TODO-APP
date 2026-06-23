package service

import (
	"context"

	"github.com/Sarthak1722/todo_app/internal/dto"
	"github.com/Sarthak1722/todo_app/internal/errors"
	"github.com/Sarthak1722/todo_app/internal/models"
	"github.com/Sarthak1722/todo_app/internal/repository"
	"github.com/Sarthak1722/todo_app/internal/validator"
)

// TodoService handles business logic for todos
type TodoService struct {
	store repository.Store
}

// NewTodoService creates a new todo service
func NewTodoService(store repository.Store) *TodoService {
	return &TodoService{
		store: store,
	}
}

// GetAllTodos retrieves all todos for a specific user
func (s *TodoService) GetAllTodos(ctx context.Context, userID string) ([]models.Todo, error) {
	return s.store.GetAllTodos(ctx, userID)
}

// GetTodoByID retrieves a single todo by ID for a specific user
func (s *TodoService) GetTodoByID(ctx context.Context, id int, userID string) (*models.Todo, error) {
	return s.store.GetTodoByID(ctx, id, userID)
}

// CreateTodo creates a new todo for a specific user
func (s *TodoService) CreateTodo(ctx context.Context, req dto.CreateTodoRequest, userID string) (*models.Todo, error, map[string]string) {
	if err := validator.Validate.Struct(req); err != nil {
		return nil, nil, errors.FormatValidationErrors(err)
	}

	todo := models.Todo{
		Body:      req.Body,
		Completed: req.Completed,
		UserID:    userID, // Bind the identity!
	}

	createdTodo, err := s.store.CreateTodo(ctx, todo, userID)
	if err != nil {
		return nil, err, nil
	}

	return &createdTodo, nil, nil
}

// DeleteTodoByID deletes a todo by ID, ensuring ownership
func (s *TodoService) DeleteTodoByID(ctx context.Context, id int, userID string) error {
	return s.store.DeleteTodoByID(ctx, id, userID)
}

// PatchTodoByID updates a todo by ID, ensuring ownership
func (s *TodoService) PatchTodoByID(ctx context.Context, id int, req dto.PatchTodoRequest, userID string) (*models.Todo, error, map[string]string) {
	if err := validator.Validate.Struct(req); err != nil {
		return nil, nil, errors.FormatValidationErrors(err)
	}

	todo := models.Todo{
		Body:      req.Body,
		Completed: req.Completed,
		UserID:    userID, // Bind the identity!
	}

	updatedTodo, err := s.store.PatchTodoByID(ctx, id, todo, userID)
	if err != nil {
		return nil, err, nil
	}

	return updatedTodo, nil, nil
}
