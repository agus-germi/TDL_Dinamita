package models

type Opinion struct {
    ID      int64  `json:"id"`   
    Name    string `json:"user_name"`    
    UserID  int64  `json:"user_id"` 
    Opinion string `json:"opinion"`
    Rating int     `json:"rating"`
}