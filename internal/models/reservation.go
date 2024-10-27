package models

type Reservation struct {
	ID      int    `json:"id"`
	TableID int    `json:"tableid"` // Table Table  `json:"table"` quizas conviene usar las estructuras Table y User dentro de Reservation.
	UserID  int    `json:"userid"`  // User  User   `json:"user"`
	Date    string `json:"date"`
	Time    string `json:"time"`
}
