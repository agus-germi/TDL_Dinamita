package dtos

type CreateTableDTO struct {
	Number   int64  `json:"number" validate:"required"`
	Seats    int64  `json:"seats"  validate:"required"`
	Description string `json:"description" validate:"required"`
}
