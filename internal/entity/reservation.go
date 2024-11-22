package entity

import "time"

type Reservation struct {
	ID              int64     `db:"id"`
	UserID          int64     `db:"reserved_by"`      // User  User   `db:"reserved_by"`  --> Quizas es mejor utilizar el email del user en vez de su ID.
	TableNumber     int64     `db:"table_number"`     // Table Table  `db:"table_number"` quizas conviene usar las estructuras Table y User dentro de Reservation.
	ReservationDate time.Time `db:"reservation_date"` // El formato de este date es ISO 8601. Para tener mas control sobre la fecha podriamos cambiar el tipo de dato a time.Time (en vez de string)
}
