package dtos

type RegisterReservationDTO struct {
	UserID      int64  `json:"userID" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Password    string `json:"password" validate:"required,min=8,max=15"`
	Email       string `json:"email" validate:"required,email"`
	TableNumber int64  `json:"tableNumber" validate:"required"`
}
