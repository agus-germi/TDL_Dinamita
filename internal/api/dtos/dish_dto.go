package dtos

type DishDTO struct {
	Name        string `json:"name" validate:"required"`
	Price       int64  `json:"price"  validate:"required"`
	Description string `json:"description"`
}
