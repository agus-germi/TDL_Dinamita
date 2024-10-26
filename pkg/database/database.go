package database

import (
    "database/sql"
    "log"
    "os"
    _ "github.com/lib/pq"
    "github.com/joho/godotenv"
)

func ConnectDB() (*sql.DB, error) {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error al cargar el archivo .env")
    }

    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")
    sslmode := os.Getenv("DB_SSLMODE")

    connStr := "user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=" + sslmode
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }

    if err := db.Ping(); err != nil {
        return nil, err
    }

    return db, nil
}
