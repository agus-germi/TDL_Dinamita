package models

type Table struct {
	Number      int64  `json:"number"`
	Seats       int64  `json:"seats"`
	Location    string `json:"location"`
	IsAvailable bool   `json:"is_available"`
}
