package models

// Aca no aparece la passwd.
// Es la estructura que vamos a usar para representar
// un usuario en la capa de servicio (aplicacion).
type User struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	RoleID string `json:"role_id"`
}
