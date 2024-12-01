package dtos

import "time"

type CreateReservationDTO struct {
	TableNumber     int64     `json:"table_number" validate:"required"`
	ReservationDate time.Time `json:"reservation_date" validate:"required,datetime=2006-01-02T15:04:05Z"` // ISO 8601 format
	// remember that we represent the timeSlot as an integer (timeSlotID)
}
