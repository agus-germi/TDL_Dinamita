package entity

type Dish struct {
	ID          int64  `db:"id"`
	Name        string `db:"name"`
	Price       int64  `db:"price"`
	Description string `db:"description"`
}
