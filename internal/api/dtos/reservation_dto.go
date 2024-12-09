package dtos

import "time"

type ReservationDTO struct {
	ID              int64     `json:"id"`
	TableNumber     int64     `json:"table_number"`
	ReservationDate time.Time `json:"reservation_date" validate:"required,datetime=2006-01-02T15:04:05Z"` // ISO 8601 format
	Promotion       string    `json:"promotion"`
}
