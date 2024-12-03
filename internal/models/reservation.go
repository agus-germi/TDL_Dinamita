package models

import "time"

type Reservation struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"reserved_by"` // Is it necessary to have this field?
	TableNumber     int64     `json:"table_number"`
	ReservationDate time.Time `json:"reservation_date" validate:"required,datetime=2006-01-02T15:04:05Z"` // ISO 8601 format
	// remember that we represent the timeSlot as an integer (timeSlotID) --> But when is sent to the front, it must be a string (or be included in reservation_date so it can be parsed)
	//ReservationDate string    `json:"reservation_date"` // El formato de este date es ISO 8601
	//TimeSlot        time.Time `json:"time_slot"`
}
