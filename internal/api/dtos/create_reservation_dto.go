package dtos

type CreateReservationDTO struct {
	TableNumber     int64  `json:"table_number" validate:"required"`
	ReservationDate string `json:"reservation_date" validate:"required,datetime=2006-01-02T15:04:05Z"` // ISO 8601 format
	PromotionID     int64  `json:"promotion_id" validate:"required"`
}
