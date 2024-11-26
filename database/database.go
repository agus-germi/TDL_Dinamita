package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// This function creates a new connection to the database
// CreatePoolConnection
func CreateConnection(ctx context.Context) (*pgxpool.Pool, error) {
	err := godotenv.Load("/usr/src/app/.env")
	if err != nil {
		log.Println("Error: .env file couldn't be loaded.")
		return nil, err
	}

	dbHost, err := getEnv("DB_HOST")
	if err != nil {
		return nil, err
	}

	dbPort, err := getEnv("DB_PORT")
	if err != nil {
		return nil, err
	}

	dbUser, err := getEnv("DB_USER")
	if err != nil {
		return nil, err
	}

	dbPasswd, err := getEnv("DB_PASSWORD")
	if err != nil {
		return nil, err
	}

	dbName, err := getEnv("DB_NAME")
	if err != nil {
		return nil, err
	}

	sslMode, err := getEnv("DB_SSLMODE")
	if err != nil {
		return nil, err
	}

	dbURL, err := getEnv("DB_URL")
	if err != nil {
		return nil, err
	}

	dbURL = strings.ReplaceAll(dbURL, "${DB_USER}", dbUser)
	dbURL = strings.ReplaceAll(dbURL, "${DB_PASS}", dbPasswd)
	dbURL = strings.ReplaceAll(dbURL, "${DB_HOST}", dbHost)
	dbURL = strings.ReplaceAll(dbURL, "${DB_PORT}", dbPort)
	dbURL = strings.ReplaceAll(dbURL, "${DB_NAME}", dbName)
	dbURL = strings.ReplaceAll(dbURL, "${DB_SSLMODE}", sslMode)

	maxConns_str, err := getEnv("MAX_CONNS")
	if err != nil {
		return nil, err
	}

	minConns_str, err := getEnv("MIN_CONNS")
	if err != nil {
		return nil, err
	}

	idleTime_str, err := getEnv("MAX_CONN_IDLE_TIME")
	if err != nil {
		return nil, err
	}

	maxConns, _ := strconv.Atoi(maxConns_str)
	minConns, _ := strconv.Atoi(os.Getenv(minConns_str))
	idleTime, _ := time.ParseDuration(os.Getenv(idleTime_str))

	// Pool Configuration
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalf("Error trying parse dbURL to create pool configuration: %v", err)
	}
	config.MaxConns = int32(maxConns)
	config.MinConns = int32(minConns)
	config.MaxConnIdleTime = idleTime // Tiempo máximo que una conexión puede estar inactiva antes de ser cerrada
	//config.MaxConnLifetime = 30 * time.Minute // Tiempo máximo de vida de una conexión (sin importar si está en uso o inactiva)
	//config.HealthCheckPeriod = 1 * time.Minute // Frecuencia de las verificaciones de salud

	dbPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Error trying to create the connections pool: %v", err)
	}

	if err := verifyDatabaseConnection(ctx, dbPool); err != nil {
		log.Fatalf("Error trying to verify the db connection: %v", err)
	}

	return dbPool, nil
}

func getEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("environment variable %s is not set", key)
	}
	return value, nil
}

func verifyDatabaseConnection(ctx context.Context, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	// Verificar que la conexión esté funcionando
	if err := conn.Conn().Ping(ctx); err != nil {
		return err
	}

	log.Println("Verification of db connection successfully done.")
	return nil
}
