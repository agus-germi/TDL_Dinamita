package dtos

type RemoveReservationDTO struct {
	ID int64 `json:"id" validate:"required"` // We have to send the reservation ID in the URI. (then this DTO will be useless)
}
