package entity

type Opinion struct {
    ID      int64  `json:"id"`
    UserID  int64  `json:"user_id"`
    Name    string `json:"name"`
    Opinion string `json:"opinion"`
    Rating  int    `json:"rating"`
}
