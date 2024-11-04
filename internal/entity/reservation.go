package entity

type Reservation struct {
	ID              int64  `db:"id"`
	UserID          int64  `db:"reserved_by"`      // User  User   `db:"reserved_by"`  --> Quizas es mejor utilizar el email del user en vez de su ID.
	TableNumber     int64  `db:"table_number"`     // Table Table  `db:"table_number"` quizas conviene usar las estructuras Table y User dentro de Reservation.
	ReservationDate string `db:"reservation_date"` // El formato de este date es ISO 8601
}
