package database

import "github.com/Sarthak1722/todo_app/models"

// DB interface defines all database operations
type DB interface {
	GetAllTodos() []models.Todo
	GetTodoByID(id int) *models.Todo
	CreateTodo(todo models.Todo) models.Todo
}

// InMemoryDB implements the DB interface using in-memory storage
type InMemoryDB struct {
	todos []models.Todo
}

// NewInMemoryDB creates a new in-memory database instance
func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		todos: []models.Todo{},
	}
}

func (db *InMemoryDB) GetAllTodos() []models.Todo {
	return db.todos
}

func (db *InMemoryDB) GetTodoByID(id int) *models.Todo {
	for i, todo := range db.todos {
		if todo.ID == id {
			return &db.todos[i]
		}
	}
	return nil
}

func (db *InMemoryDB) CreateTodo(todo models.Todo) models.Todo {
	id := 1
	if len(db.todos) > 0 {
		id = db.todos[len(db.todos)-1].ID + 1
	}
	todo.ID = id
	db.todos = append(db.todos, todo)
	return todo
}
