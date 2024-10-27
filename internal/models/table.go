package models

type Table struct {
	ID          int    `json:"id"`
	Seats       int    `json:"seats"`
	Location    string `json:"location"`
	IsAvailable bool   `json:"is_available"`
}
