package repository

import (
	"context"
	"errors"
	"slices"
	"sync"

	"github.com/Sarthak1722/todo_app/internal/models"
)

// InMemoryTodoRepository implements Store using in-memory storage
type InMemoryTodoRepository struct {
	todos []models.Todo
	mu    sync.RWMutex
}

// NewInMemoryTodoRepository creates a new in-memory repository instance
func NewInMemoryTodoRepository() *InMemoryTodoRepository {
	return &InMemoryTodoRepository{
		todos: []models.Todo{},
	}
}

func (s *InMemoryTodoRepository) GetAllTodos(ctx context.Context, userID string) ([]models.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todos := make([]models.Todo, len(s.todos))
	copy(todos, s.todos)
	todos = slices.DeleteFunc(todos, func(todo models.Todo) bool {
		return todo.UserID != userID
	})

	return todos, nil
}

func (s *InMemoryTodoRepository) GetTodoByID(ctx context.Context, id int, userID string) (*models.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i, todo := range s.todos {
		if todo.ID == id && todo.UserID == userID {
			return &s.todos[i], nil
		}
	}

	return nil, nil
}

func (s *InMemoryTodoRepository) CreateTodo(ctx context.Context, todo models.Todo, userID string) (models.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := 1
	if len(s.todos) > 0 {
		id = s.todos[len(s.todos)-1].ID + 1
	}

	todo.ID = id
	todo.UserID = userID
	s.todos = append(s.todos, todo)

	return todo, nil
}

func (s *InMemoryTodoRepository) DeleteTodoByID(ctx context.Context, id int, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	originalLen := len(s.todos)

	s.todos = slices.DeleteFunc(s.todos, func(e models.Todo) bool {
		return e.ID == id && e.UserID == userID
	})

	if originalLen == len(s.todos) {
		return errors.New("todo not found")
	}

	return nil
}

func (s *InMemoryTodoRepository) PatchTodoByID(ctx context.Context, id int, todo models.Todo, userID string) (*models.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, existingTodo := range s.todos {
		if existingTodo.ID == id && existingTodo.UserID == userID {
			todo.ID = existingTodo.ID
			todo.UserID = userID

			if todo.Body == "" {
				todo.Body = existingTodo.Body
			}

			s.todos[i] = todo
			return &s.todos[i], nil
		}
	}

	return nil, nil
}
