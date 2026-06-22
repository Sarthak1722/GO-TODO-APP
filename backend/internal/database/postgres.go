package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgresDB initializes a connection pool to the database.
func NewPostgresDB(dsn string) (*pgxpool.Pool, error) {
	// 1. Create the pool configuration
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	// 2. Establish the connection pool
	// We use context.Background() here because this is app startup, not a specific web request
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 3. Ping the database to verify the connection is actually alive
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("database didn't respond to ping: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database!")
	return pool, nil
}

// We return a *pgxpool.Pool: This allows your whole app to safely share connections concurrently.

// We use Ping(): Just because pgxpool.New succeeds doesn't mean the database is actually reachable (it initializes the pool lazily). Ping() forces a real network check so your app crashes immediately on startup if the DB is down, rather than failing on the first user request.

// Error Wrapping (%w): We wrap errors so higher levels of your app know exactly what failed.
