package services

import (
    "database/sql"
    "golang.org/x/crypto/bcrypt"
    "TDL_Dinamita/pkg/database"
)

// VerifyUser verifica si el usuario existe y si la contraseña es correcta
func VerifyUser(db *sql.DB, username, password string) (bool, error) {
    user, err := database.GetUserByUsername(db, username)
    if err != nil {
        return false, err
    }
    if user == nil {
        return false, nil // Usuario no encontrado
    }

    // Comparar la contraseña hasheada
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return false, nil // Contraseña incorrecta
    }

    return true, nil // Usuario validado
}

// CreateUser crea un nuevo usuario
func CreateUser(db *sql.DB, username, password, email string) error {
    // Hashear la contraseña antes de almacenarla
    hashedPassword, err := HashPassword(password)
    if err != nil {
        return err
    }

    return database.CreateUser(db, username, hashedPassword, email)
}

// HashPassword hashea la contraseña
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}
