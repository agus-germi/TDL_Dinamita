package models

import "time"

type Reservation struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"reserved_by"`
	TableNumber     int64     `json:"table_number"`
	ReservationDate string    `json:"reservation_date"` // El formato de este date es ISO 8601
	TimeSlot        time.Time `json:"time_slot"`
}
