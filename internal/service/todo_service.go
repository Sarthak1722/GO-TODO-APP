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

// GetAllTodos retrieves all todos
func (s *TodoService) GetAllTodos(ctx context.Context) ([]models.Todo, error) {
	return s.store.GetAllTodos(ctx)
}

// GetTodoByID retrieves a single todo by ID
func (s *TodoService) GetTodoByID(ctx context.Context, id int) (*models.Todo, error) {
	return s.store.GetTodoByID(ctx, id)
}

// CreateTodo creates a new todo
func (s *TodoService) CreateTodo(ctx context.Context, req dto.CreateTodoRequest) (*models.Todo, map[string]string) {
	if err := validator.Validate.Struct(req); err != nil {
		return nil, errors.FormatValidationErrors(err)
	}

	todo := models.Todo{
		Body:      req.Body,
		Completed: req.Completed,
	}

	createdTodo, _ := s.store.CreateTodo(ctx, todo)

	return &createdTodo, nil
}

// DeleteTodo deletes a todo by ID
func (s *TodoService) DeleteTodo(ctx context.Context, id int) error {
	return s.store.DeleteTodoByID(ctx, id)
}

// PatchTodo updates a todo by ID
func (s *TodoService) PatchTodo(ctx context.Context, id int, req dto.PatchTodoRequest) (*models.Todo, map[string]string) {
	if err := validator.Validate.Struct(req); err != nil {
		return nil, errors.FormatValidationErrors(err)
	}

	todo := models.Todo{
		Body:      req.Body,
		Completed: req.Completed,
	}

	updatedTodo, _ := s.store.PatchTodoByID(ctx, id, todo)

	return updatedTodo, nil
}