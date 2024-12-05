package models

import "time"

type Reservation struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"reserved_by"` // Is it necessary to have this field?
	TableNumber     int64     `json:"table_number"`
	ReservationDate time.Time `json:"reservation_date" validate:"required,datetime=2006-01-02T15:04:05Z"` // ISO 8601 format

	// Date        string `json:"reservation_date"`
	// TimeSlot    string `json:"time_slot"`
}
