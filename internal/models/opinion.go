package models

type Opinion struct {
    ID      int64  `json:"id"`     
    UserID  int64  `json:"user_id"` 
    Opinion string `json:"opinion"`
}