package database

import (
    "database/sql"
    "golang.org/x/crypto/bcrypt" // Para hashear contraseñas
    _ "github.com/lib/pq" // Importar el controlador de PostgreSQL
)

// User representa la estructura de un usuario en la base de datos
type User struct {
    ID       int
    Username string
    Password string
    Email    string
}

// GetUserByUsername busca un usuario por su nombre de usuario
func GetUserByUsername(db *sql.DB, username string) (*User, error) {
    var user User
    query := "SELECT id, username, password, email FROM users WHERE username = $1"
    row := db.QueryRow(query, username)
    
    if err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email); err != nil {
        if err == sql.ErrNoRows {
            return nil, nil // Usuario no encontrado
        }
        return nil, err // Otro error
    }
    
    return &user, nil
}

// CreateUser crea un nuevo usuario en la base de datos
func CreateUser(db *sql.DB, username, password, email string) error {
    query := "INSERT INTO users (username, password, email) VALUES ($1, $2, $3)"
    _, err := db.Exec(query, username, password, email)
    return err
}

// hashPassword hashea la contraseña utilizando bcrypt
func hashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}
