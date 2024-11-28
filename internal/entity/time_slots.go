package entity

import "time"

type TimeSlot struct {
	ID   int64     `db:"id"`
	Time time.Time `db:"time"`
}
