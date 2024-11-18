package dtos

type RemoveTableDTO struct {
	Number int64 `json:"number" validate:"required"`
}
