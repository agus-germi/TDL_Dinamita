package database

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/agus-germi/TDL_Dinamita/logger"
	"github.com/agus-germi/TDL_Dinamita/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// This function creates a new connection to the database
// CreatePoolConnection
func CreateConnection(ctx context.Context) (*pgxpool.Pool, error) {
	err := godotenv.Load("/usr/src/app/.env")
	if err != nil {
		logger.Log.Errorf("'.env' file couldn't be loaded: %v", err)
		return nil, err
	}

	dbHost, err := utils.GetEnv("DB_HOST")
	if err != nil {
		return nil, err
	}

	dbPort, err := utils.GetEnv("DB_PORT")
	if err != nil {
		return nil, err
	}

	dbUser, err := utils.GetEnv("DB_USER")
	if err != nil {
		return nil, err
	}

	dbPasswd, err := utils.GetEnv("DB_PASSWORD")
	if err != nil {
		return nil, err
	}

	dbName, err := utils.GetEnv("DB_NAME")
	if err != nil {
		return nil, err
	}

	sslMode, err := utils.GetEnv("DB_SSLMODE")
	if err != nil {
		return nil, err
	}

	dbURL, err := utils.GetEnv("DB_URL")
	if err != nil {
		return nil, err
	}

	dbURL = strings.ReplaceAll(dbURL, "${DB_USER}", dbUser)
	dbURL = strings.ReplaceAll(dbURL, "${DB_PASS}", dbPasswd)
	dbURL = strings.ReplaceAll(dbURL, "${DB_HOST}", dbHost)
	dbURL = strings.ReplaceAll(dbURL, "${DB_PORT}", dbPort)
	dbURL = strings.ReplaceAll(dbURL, "${DB_NAME}", dbName)
	dbURL = strings.ReplaceAll(dbURL, "${DB_SSLMODE}", sslMode)

	maxConnsStr, err := utils.GetEnv("MAX_CONNS")
	if err != nil {
		return nil, err
	}

	minConnsStr, err := utils.GetEnv("MIN_CONNS")
	if err != nil {
		return nil, err
	}

	idleTimeStr, err := utils.GetEnv("MAX_CONN_IDLE_TIME")
	if err != nil {
		return nil, err
	}

	maxConns, err := strconv.Atoi(maxConnsStr)
	if err != nil {
		logger.Log.Fatalf("Error trying to convert MAX_CONNS environment variable to int: %v", err)
		return nil, err
	}
	minConns, err := strconv.Atoi(minConnsStr)
	if err != nil {
		logger.Log.Fatalf("Error trying to convert MIN_CONNS environment variable to int: %v", err)
		return nil, err
	}
	idleTime, err := time.ParseDuration(idleTimeStr)
	if err != nil {
		logger.Log.Fatalf("Error trying to parse duration of MAX_CONN_IDLE_TIME environment variable: %v", err)
		return nil, err
	}

	// Pool Configuration
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		logger.Log.Fatalf("Error trying parse dbURL to create pool configuration: %v", err)
	}
	config.MaxConns = int32(maxConns)
	config.MinConns = int32(minConns)
	config.MaxConnIdleTime = idleTime // Tiempo máximo que una conexión puede estar inactiva antes de ser cerrada
	//config.MaxConnLifetime = 30 * time.Minute // Tiempo máximo de vida de una conexión (sin importar si está en uso o inactiva)
	//config.HealthCheckPeriod = 1 * time.Minute // Frecuencia de las verificaciones de salud

	dbPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		logger.Log.Fatalf("Error trying to create the connections pool: %v", err)
	}

	if err := verifyDatabaseConnection(ctx, dbPool); err != nil {
		logger.Log.Fatalf("Error trying to verify the db connection: %v", err)
	}

	return dbPool, nil
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

	logger.Log.Infof("Verification of db connection successfully done.")
	return nil
}
