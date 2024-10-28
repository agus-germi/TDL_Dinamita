package models

// Aca no aparece la passwd.
// Es la estructura que vamos a usar para representar
// un usuario en la capa de servicio (aplicacion).
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
