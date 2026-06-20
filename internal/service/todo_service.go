package service

import (
	"github.com/Sarthak1722/todo_app/internal/dto"
	"github.com/Sarthak1722/todo_app/internal/errors"
	"github.com/Sarthak1722/todo_app/internal/models"
	"github.com/Sarthak1722/todo_app/internal/store"
	"github.com/Sarthak1722/todo_app/internal/validator"
)

// TodoService handles business logic for todos
type TodoService struct {
	store store.Store
}

// NewTodoService creates a new todo service
func NewTodoService(store store.Store) *TodoService {
	return &TodoService{
		store: store,
	}
}

// GetAllTodos retrieves all todos
func (s *TodoService) GetAllTodos() []models.Todo {
	return s.store.GetAllTodos()
}

// GetTodoByID retrieves a single todo by ID
func (s *TodoService) GetTodoByID(id int) *models.Todo {
	return s.store.GetTodoByID(id)
}

// CreateTodo creates a new todo (with validation)
func (s *TodoService) CreateTodo(req dto.CreateTodoRequest) (*models.Todo, map[string]string) {
	// Validate request
	if err := validator.Validate.Struct(req); err != nil {
		return nil, errors.FormatValidationErrors(err)
	}

	// Convert DTO to model
	todo := models.Todo{
		Body:      req.Body,
		Completed: req.Completed,
	}

	// Create in store
	createdTodo := s.store.CreateTodo(todo)
	return &createdTodo, nil
}

// DeleteTodo deletes a todo by ID
func (s *TodoService) DeleteTodo(id int) bool {
	return s.store.DeleteTodoByID(id)
}

// PatchTodo updates a todo by ID
func (s *TodoService) PatchTodo(id int, req dto.PatchTodoRequest) (*models.Todo, map[string]string) {
	// Validate request
	if err := validator.Validate.Struct(req); err != nil {
		return nil, errors.FormatValidationErrors(err)
	}

	todo := models.Todo{
		Body:      req.Body,
		Completed: req.Completed,
	}

	return s.store.PatchTodoByID(id, todo), nil
}