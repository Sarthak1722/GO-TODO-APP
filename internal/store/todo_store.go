package store

import (
	"slices"

	"github.com/Sarthak1722/todo_app/internal/models"
)

// Store interface defines all data access operations
type Store interface {
	GetAllTodos() []models.Todo
	GetTodoByID(id int) *models.Todo
	CreateTodo(todo models.Todo) models.Todo
	DeleteTodoByID(id int) bool
	PatchTodoByID(id int, todo models.Todo) *models.Todo
}

// InMemoryStore implements the Store interface using in-memory storage
type InMemoryStore struct {
	todos []models.Todo
}

// NewInMemoryStore creates a new in-memory store instance
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		todos: []models.Todo{},
	}
}

func (s *InMemoryStore) GetAllTodos() []models.Todo {
	return s.todos
}

func (s *InMemoryStore) GetTodoByID(id int) *models.Todo {
	for i, todo := range s.todos {
		if todo.ID == id {
			return &s.todos[i]
		}
	}
	return nil
}

func (s *InMemoryStore) CreateTodo(todo models.Todo) models.Todo {
	id := 1
	if len(s.todos) > 0 {
		id = s.todos[len(s.todos)-1].ID + 1
	}
	todo.ID = id
	s.todos = append(s.todos, todo)
	return todo
}

func (s *InMemoryStore) DeleteTodoByID(id int) bool {
	originalLen := len(s.todos)
	s.todos = slices.DeleteFunc(s.todos, func(e models.Todo) bool {
		return e.ID == id
	})
	return originalLen > len(s.todos)
}

func (s *InMemoryStore) PatchTodoByID(id int, todo models.Todo) *models.Todo {
	for i, existingTodo := range s.todos {
		if existingTodo.ID == id {
			todo.ID = existingTodo.ID
			if todo.Body==""{
				todo.Body=existingTodo.Body
			}
			s.todos[i] = todo
			return &s.todos[i]
		}
	}
	return nil
}
