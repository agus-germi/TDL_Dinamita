package entity

type Table struct {
	ID          int64  `db:"id"`
	Number      int64  `db:"number"`
	Seats       int64  `db:"seats"`
	Location    string `db:"location"`
	IsAvailable bool   `db:"is_available"`
}
