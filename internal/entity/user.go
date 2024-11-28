package entity

type User struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Password string `db:"password"`
	Email    string `db:"email"`
	Role     int32  `db:"role"`
}
