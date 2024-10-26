package main

import (
    "database/sql"
    "log"
    "net/http"
    "TDL_Dinamita/pkg/database"
    "TDL_Dinamita/internal/handlers" // Cambia esto a la ruta correcta
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file" // Importación necesaria para las migraciones
)

// Middleware para habilitar CORS
func enableCors(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", "*") // Permite todos los orígenes. Cambia "*" por tu dominio en producción.
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // Métodos permitidos
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Encabezados permitidos
}

func main() {
    db, err := database.ConnectDB() // Manejo de ambos valores
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Ejecutar migraciones
    if err := migrateDB(db); err != nil {
        log.Fatalf("Error en las migraciones: %v", err)
    }

    // Ruta de bienvenida
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("¡Bienvenido al sistema de reservas!"))
    })

    // Rutas para el login y registro con CORS
    http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        enableCors(w) // Habilita CORS
        if r.Method == http.MethodOptions {
            return // Responde a las solicitudes OPTIONS
        }
        handlers.LoginHandler(w, r) // Usa el handler desde el paquete handlers
    })

    http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
        enableCors(w) // Habilita CORS
        if r.Method == http.MethodOptions {
            return // Responde a las solicitudes OPTIONS
        }
        handlers.RegisterHandler(w, r) // Usa el handler desde el paquete handlers
    })

    log.Println("Servidor en ejecución en :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func migrateDB(db *sql.DB) error {
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        return err
    }

    // Crear una instancia de migración desde el directorio de archivos
    m, err := migrate.NewWithDatabaseInstance(
        "file://migrations", // Asegúrate de que esta ruta sea correcta
        "postgres", driver,
    )
    if err != nil {
        return err
    }

    // Aplicar todas las migraciones pendientes
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }

    return nil
}
