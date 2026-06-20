package repository

import (
	"context"

	"github.com/Sarthak1722/todo_app/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store interface defines all data access operations
type Store interface {
    GetAllTodos(ctx context.Context) ([]models.Todo, error)
    GetTodoByID(ctx context.Context, id int) (*models.Todo, error)
    CreateTodo(ctx context.Context, todo models.Todo) (models.Todo, error)
    DeleteTodoByID(ctx context.Context, id int) error
    PatchTodoByID(ctx context.Context, id int, todo models.Todo) (*models.Todo, error)
}

type PostgresTodoRepository struct {
	db *pgxpool.Pool
}

// NewPostgresTodoRepository injects the DB pool and returns a Store interface
func NewPostgresTodoRepository(db *pgxpool.Pool) *PostgresTodoRepository {
	return &PostgresTodoRepository{
		db: db,
	}
}

func (p *PostgresTodoRepository) GetAllTodos(ctx context.Context) ([]models.Todo, error) {
	return []models.Todo{}, nil
}

func (p *PostgresTodoRepository) GetTodoByID(ctx context.Context, id int) (*models.Todo, error) {
	return nil, nil
}

func (p *PostgresTodoRepository) CreateTodo(ctx context.Context, todo models.Todo) (models.Todo, error) {
	return models.Todo{}, nil
}

func (p *PostgresTodoRepository) DeleteTodoByID(ctx context.Context, id int) error {
	return nil
}

func (p *PostgresTodoRepository) PatchTodoByID(ctx context.Context, id int, todo models.Todo) (*models.Todo, error) {
	return nil, nil
}
