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
func NewPostgresTodoRepository(db *pgxpool.Pool) Store {
	return &PostgresTodoRepository{
		db: db,
	}
}

func (p *PostgresTodoRepository) GetAllTodos(ctx context.Context) ([]models.Todo, error) {
	query := `SELECT id, body, completed FROM todos ORDER BY id ASC`

	// 1. Execute the query
	rows, err := p.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	// 2. CRUCIAL: Always defer closing the rows to prevent connection leaks
	defer rows.Close()
	// Notice defer rows.Close(). When you do a SELECT, Postgres holds a connection open to stream the data to you. If you forget to close it, that connection is locked forever, and your pool will eventually run out of connections.

	// 3. Iterate through the rows
	todos := []models.Todo{}
	for rows.Next() {
		var t models.Todo
		// Scan the current row's columns into our struct
		if err := rows.Scan(&t.ID, &t.Body, &t.Completed); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}

	// 4. Check for errors that might have occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (p *PostgresTodoRepository) GetTodoByID(ctx context.Context, id int) (*models.Todo, error) {
	query := `SELECT id, body, completed FROM todos WHERE id = $1`

	var t models.Todo
	err := p.db.QueryRow(ctx, query, id).Scan(&t.ID, &t.Body, &t.Completed)

	if err != nil {
		// If the error is literally "no rows exist", we gracefully return nil just like the in-memory store
		if errors.Is(err, pgx.ErrNoRows) {
			// In Postgres, the driver throws a specific error: pgx.ErrNoRows
			return nil, nil
		}
		// If it's a real error (e.g., database went offline), return it
		return nil, err
	}

	return &t, nil
}

func (p *PostgresTodoRepository) CreateTodo(ctx context.Context, todo models.Todo) (models.Todo, error) {
	// 1. Define the SQL query
	query := `
		INSERT INTO todos (body, completed)
		VALUES ($1, $2)
		RETURNING id
	`

	// 2. Execute the query and scan the returned ID back into our struct
	// QueryRow asks for exactly one row back. We pass the context, the query, and our variables.
	err := p.db.QueryRow(ctx, query, todo.Body, todo.Completed).Scan(&todo.ID)
	if err != nil {
		return models.Todo{}, err
	}

	// 3. Return the completed struct (now containing the DB-generated ID)
	return todo, nil
}

func (p *PostgresTodoRepository) DeleteTodoByID(ctx context.Context, id int) error {
	query := `DELETE FROM todos WHERE id = $1`
	
	// Exec returns a CommandTag, which tells us what happened
	cmdTag, err := p.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	// If no rows were actually deleted, the ID didn't exist
	// A DELETE statement in SQL does not throw an error if the row doesn't exist; it just deletes zero rows and reports success. We use RowsAffected() == 0 to realize nothing happened and return an error to match your previous in-memory logic.
	if cmdTag.RowsAffected() == 0 {
		var ErrTodoNotFound = errors.New("todo not found")
		return ErrTodoNotFound
	}

	return nil
}

func (p *PostgresTodoRepository) PatchTodoByID(ctx context.Context, id int, todo models.Todo) (*models.Todo, error) {

	// Patching is historically tricky in SQL because you only want to update the fields the user actually provided. We can use a clever SQL function called COALESCE, paired with NULLIF, to conditionally update the body only if it isn't an empty string.
    // We also use RETURNING * so Postgres hands us back the completely updated row in one single network trip.

	query := `
		UPDATE todos
		SET 
			-- If $1 (todo.Body) is an empty string, treat it as NULL. 
			-- COALESCE then falls back to keeping the existing 'body' value.
			body = COALESCE(NULLIF($1, ''), body),
			completed = $2
		WHERE id = $3
		RETURNING id, body, completed
	`

	var t models.Todo
	err := p.db.QueryRow(ctx, query, todo.Body, todo.Completed, id).Scan(&t.ID, &t.Body, &t.Completed)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Todo to update was not found
		}
		return nil, err
	}

	return &t, nil
}
