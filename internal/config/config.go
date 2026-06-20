package config

import (
	"os"
	"strings"
)

type Config struct {
	DBType      string
	PostgresDSN string
}

// Load reads configuration from environment variables
func Load() *Config {
	dbType := strings.ToLower(os.Getenv("DB_TYPE"))
	if dbType == "" {
		dbType = "memory" // default to in-memory
	}

	postgresDSN := os.Getenv("POSTGRES_DSN")
	if postgresDSN == "" {
		postgresDSN = "postgres://sarthak@localhost:5432/todo_app"
	}

	return &Config{
		DBType:      dbType,
		PostgresDSN: postgresDSN,
	}
}
