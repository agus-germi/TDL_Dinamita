package entity

type Reservation struct {
	ID          int64  `db:"id"`
	TableNumber int64  `json:"table_number"` // Table Table  `json:"table"` quizas conviene usar las estructuras Table y User dentro de Reservation.
	UserID      int64  `json:"user_id"`      // User  User   `json:"user"`  --> Quizas es mejor utilizar el email del user en vez de su ID.
	Date        string `db:"date"`
	Time        string `db:"time"`
}
