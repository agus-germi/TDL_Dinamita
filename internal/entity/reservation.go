package entity

type Reservation struct {
	ID      int    `db:"id"`
	TableID int    `db:"tableid"` // Table Table  `db:"table"` quizas conviene usar las estructuras Table y User dentro de Reservation.
	UserID  int    `db:"userid"`  // User  User   `db:"user"`
	Date    string `db:"date"`
	Time    string `db:"time"`
}
