package models

type Table struct {
    ID           int     `json:"id"`
    Seats        int     `json:"seats"`
    Location     string  `json:"location"`
    IsAvailable  bool    `json:"is_available"`
}

type Reservation struct {
    ID        int    `json:"id"`
    TableID   int    `json:"table_id"`
    UserID    int    `json:"user_id"`
    Date      string `json:"date"` 
    Time      string `json:"time"`
}

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"` 
    Email    string `json:"email"`
}