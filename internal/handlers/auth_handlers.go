package handlers // Cambia esto si es necesario

import (
    "encoding/json"
    "net/http"
    "TDL_Dinamita/pkg/database"
    "TDL_Dinamita/internal/models"
    "TDL_Dinamita/internal/services" // Importar el paquete donde están VerifyUser y CreateUser
)

// LoginHandler maneja las solicitudes de inicio de sesión
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
        return
    }

    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
        return
    }

    db, err := database.ConnectDB()
    if err != nil {
        http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    // Verificar el usuario
    valid, err := services.VerifyUser(db, user.Username, user.Password) // Asegúrate de que esto sea correcto
    if err != nil || !valid {
        http.Error(w, "Usuario o contraseña incorrectos", http.StatusUnauthorized)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Inicio de sesión exitoso"))
}

// RegisterHandler maneja las solicitudes de registro
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
        return
    }

    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
        return
    }

    db, err := database.ConnectDB()
    if err != nil {
        http.Error(w, "Error al conectar a la base de datos", http.StatusInternalServerError)
        return
    }
    defer db.Close()

    // Crear el usuario
    err = services.CreateUser(db, user.Username, user.Password, user.Email) // Asegúrate de que esto sea correcto
    if err != nil {
        http.Error(w, "Error al crear el usuario", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Usuario creado exitosamente"))
}
