package dtos

type RemoveReservationDTO struct {
	ID int64 `json:"id" validate:"required"`
}
