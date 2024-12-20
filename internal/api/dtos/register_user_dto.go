package dtos

type RegisterUserDTO struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=15"`
	Email    string `json:"email" validate:"required,email"`
}
