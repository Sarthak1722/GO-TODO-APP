package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	DBType      string
	PostgresDSN string
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
		switch env := os.Getenv("ENV"); env {
		case "dev":
			postgresDSN = os.Getenv("POSTGRES_DEV_DSN")
		case "prod":
			postgresDSN = os.Getenv("POSTGRES_PROD_DSN")
		default:
			log.Fatalf("config: DB_TYPE is postgres but ENV is %q — must be \"dev\" or \"prod\"", env)
		}

		if postgresDSN == "" {
			log.Fatalf("config: DB_TYPE is postgres but the DSN for ENV=%q is empty", os.Getenv("ENV"))
		}
	}

	return &Config{
		DBType:      dbType,
		PostgresDSN: postgresDSN,
	}
}
