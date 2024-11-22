package dtos

// Maybe we can change this awful name
type GetReservationsDTO struct {
	UserID int64 `json:"user_id" validate:"required"`
}
