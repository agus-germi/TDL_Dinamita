package dtos

type RemoveReservationDTO struct {
	UserID int64 `json:"user_id" validate:"required"`
}
