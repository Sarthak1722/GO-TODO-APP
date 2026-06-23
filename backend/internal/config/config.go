package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	DBType      string
	PostgresDSN string
	ENV         string
	USE_DB      string
}

// Load reads configuration from environment variables.
// Required env vars when DB_TYPE=postgres:
//   - ENV: "dev" or "prod"
//   - POSTGRES_DEV_DSN  (when ENV=dev)
//   - POSTGRES_PROD_DSN (when ENV=prod)
func Load() *Config {
	dbType := strings.ToLower(os.Getenv("DB_TYPE"))
	if dbType == "" {
		dbType = "memory"
	}

	var postgresDSN string

	if dbType == "postgres" {
		switch env := os.Getenv("USE_DB"); env {
		case "local":
			postgresDSN = os.Getenv("POSTGRES_DEV_DSN")
		case "remote":
			postgresDSN = os.Getenv("POSTGRES_PROD_DSN")
		default:
			log.Fatalf("config: DB_TYPE is postgres but USE_DB is %q — must be \"local\" or \"remote\"", env)
		}

		if postgresDSN == "" {
			log.Fatalf("config: DB_TYPE is postgres but the DSN for USE_DB=%q is empty", os.Getenv("USE_DB"))
		}
	}

	var ENVIRONMENT = strings.ToLower(os.Getenv("ENV"))
	if ENVIRONMENT == "" {
		ENVIRONMENT = "dev"
	}
	var USE_DB = strings.ToLower(os.Getenv("USE_DB"))
	if USE_DB == "" {
		USE_DB = "local"
	}
	return &Config{
		DBType:      dbType,
		PostgresDSN: postgresDSN,
		ENV:         ENVIRONMENT,
		USE_DB:      USE_DB,
	}
}
