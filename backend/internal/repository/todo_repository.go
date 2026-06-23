package repository

import (
	"context"
	"errors"

	"github.com/Sarthak1722/todo_app/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store interface defines all data access operations
type Store interface {
	GetAllTodos(ctx context.Context, userID string) ([]models.Todo, error)
	GetTodoByID(ctx context.Context, id int, userID string) (*models.Todo, error)
	CreateTodo(ctx context.Context, todo models.Todo, userID string) (models.Todo, error)
	DeleteTodoByID(ctx context.Context, id int, userID string) error
	PatchTodoByID(ctx context.Context, id int, todo models.Todo, userID string) (*models.Todo, error)
}

type PostgresTodoRepository struct {
	db *pgxpool.Pool
}

// NewPostgresTodoRepository injects the DB pool and returns a Store interface
func NewPostgresTodoRepository(db *pgxpool.Pool) Store {
	return &PostgresTodoRepository{
		db: db,
	}
}

func (p *PostgresTodoRepository) GetAllTodos(ctx context.Context, userID string) ([]models.Todo, error) {
	query := `SELECT id, body, completed FROM todos WHERE user_id = $1 ORDER BY id ASC`

	rows, err := p.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []models.Todo{}
	for rows.Next() {
		var t models.Todo
		if err := rows.Scan(&t.ID, &t.Body, &t.Completed); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (p *PostgresTodoRepository) GetTodoByID(ctx context.Context, id int, userID string) (*models.Todo, error) {
	query := `SELECT id, body, completed FROM todos WHERE id = $1 AND user_id = $2`

	var t models.Todo
	err := p.db.QueryRow(ctx, query, id, userID).Scan(&t.ID, &t.Body, &t.Completed)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &t, nil
}

func (p *PostgresTodoRepository) CreateTodo(ctx context.Context, todo models.Todo, userID string) (models.Todo, error) {
	query := `
        INSERT INTO todos (user_id, body, completed)
        VALUES ($1, $2, $3)
        RETURNING id
    `

	err := p.db.QueryRow(ctx, query, userID, todo.Body, todo.Completed).Scan(&todo.ID)
	if err != nil {
		return models.Todo{}, err
	}

	return todo, nil
}

func (p *PostgresTodoRepository) DeleteTodoByID(ctx context.Context, id int, userID string) error {
	query := `DELETE FROM todos WHERE id = $1 AND user_id = $2`

	cmdTag, err := p.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("todo not found")
	}

	return nil
}

func (p *PostgresTodoRepository) PatchTodoByID(ctx context.Context, id int, todo models.Todo, userID string) (*models.Todo, error) {
	query := `
        UPDATE todos
        SET 
            body = COALESCE(NULLIF($1, ''), body),
            completed = $2
        WHERE id = $3 AND user_id = $4
        RETURNING id, body, completed
    `

	var t models.Todo
	err := p.db.QueryRow(ctx, query, todo.Body, todo.Completed, id, userID).Scan(&t.ID, &t.Body, &t.Completed)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &t, nil
}
