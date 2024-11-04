package dtos

type Reservation struct {
	UserID          int64  `json:"reserved_by"`      // User  User   `json:"reserved_by"`  --> Quizas es mejor utilizar el email del user en vez de su ID.
	TableNumber     int64  `json:"table_number"`     // Table Table  `json:"table_number"` quizas conviene usar las estructuras Table y User dentro de Reservation.
	ReservationDate string `json:"reservation_date"` // El formato de este date es ISO 8601
}
