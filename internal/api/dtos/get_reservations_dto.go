package dtos

// Maybe we can change this awful name --> Deberiamos enviar el user_id en la URI
type GetReservationsDTO struct {
	UserID int64 `json:"user_id" validate:"required"`
}
