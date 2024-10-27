package entity

type Table struct {
	ID          int    `db:"id"`
	Seats       int    `db:"seats"`
	Location    string `db:"location"`
	IsAvailable bool   `db:"is_available"`
}
