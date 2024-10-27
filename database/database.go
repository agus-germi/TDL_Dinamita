package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func New(ctx context.Context) (*sqlx.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	host, err := getEnv("DB_HOST")
	if err != nil {
		return nil, err
	}

	port, err := getEnv("DB_PORT")
	if err != nil {
		return nil, err
	}

	user, err := getEnv("DB_USER")
	if err != nil {
		return nil, err
	}

	password, err := getEnv("DB_PASSWORD")
	if err != nil {
		return nil, err
	}

	dbname, err := getEnv("DB_NAME")
	if err != nil {
		return nil, err
	}

	sslmode, err := getEnv("DB_SSLMODE")
	if err != nil {
		return nil, err
	}

	connStr := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=" + sslmode
	db, err := sqlx.ConnectContext(ctx, "postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func getEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("environment variable %s is not set", key)
	}
	return value, nil
}
