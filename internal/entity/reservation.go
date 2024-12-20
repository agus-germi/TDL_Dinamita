package entity

import "time"

type Reservation struct {
	ID          int64     `db:"id"`
	UserID      int64     `db:"reserved_by"`
	TableNumber int64     `db:"table_number"`
	Date        time.Time `db:"date"`
	Time        string    `db:"time"`
	Promotion   string 	  `db:"promotion"`
}
