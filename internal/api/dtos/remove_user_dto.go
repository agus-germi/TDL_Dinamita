package dtos

type RemoveUserDTO struct {
	UserID int64 `json:"user_id" validate:"required"`
}
