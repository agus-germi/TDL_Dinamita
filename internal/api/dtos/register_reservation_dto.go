package dtos

import "time"

type RegisterReservationDTO struct {
	UserID          int64     `json:"user_id" validate:"required"`
	Name            string    `json:"name" validate:"required"`
	Password        string    `json:"password" validate:"required,min=8,max=15"` // Cambiar por Token (investigar)
	Email           string    `json:"email" validate:"required,email"`
	TableNumber     int64     `json:"table_number" validate:"required"`
	ReservationDate time.Time `json:"reservation_date" validate:"required"`
}