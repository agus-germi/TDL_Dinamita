package database

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// This function creates a new connection to the database
func CreateConnection(ctx context.Context) (*sqlx.DB, error) {
	err := godotenv.Load("/usr/src/app/.env")
	if err != nil {
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

	db, err := sqlx.ConnectContext(ctx, "postgres", dbURL)
	if err != nil {
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

/*
import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

func New() {
    // Contexto para gestionar la conexión
    ctx := context.Background()

    // Configurar el pool usando Config
    connStr := "postgres://admin:1234@localhost:5432/mi_base?sslmode=disable"
    config, err := pgxpool.ParseConfig(connStr)
    if err != nil {
        log.Fatalf("Error al parsear la configuración del pool: %v", err)
    }

    // Configuración del pool
    config.MaxConns = 10                      // Máximo de conexiones abiertas
    config.MinConns = 2                       // Mínimo de conexiones abiertas
    config.MaxConnIdleTime = 5 * time.Minute  // Tiempo máximo de inactividad de una conexión
    config.MaxConnLifetime = 30 * time.Minute // Tiempo máximo de vida de una conexión
    config.HealthCheckPeriod = 1 * time.Minute // Frecuencia de las verificaciones de salud

    // Crear el pool de conexiones
    pool, err := pgxpool.NewWithConfig(ctx, config)
    if err != nil {
        log.Fatalf("No se pudo crear el pool: %v", err)
    }
    defer pool.Close() // Asegurarse de cerrar el pool al final

    // Verificar la conexión con un simple query
    if err := testConnection(ctx, pool); err != nil {
        log.Fatalf("Conexión fallida: %v", err)
    }

    fmt.Println("Conexión establecida exitosamente con el pool de pgxpool")
}

// Función para probar la conexión con un query
func testConnection(ctx context.Context, pool *pgxpool.Pool) error {
    var now time.Time
    err := pool.QueryRow(ctx, "SELECT NOW()").Scan(&now)
    if err != nil {
        return err
    }
    fmt.Printf("Hora actual de la base de datos: %s\n", now)
    return nil
}

*/
