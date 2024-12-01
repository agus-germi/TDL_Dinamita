package dtos

type UpdateUserRoleDTO struct {
	UserID    int64 `json:"user_id" validate:"required"`
	NewRoleID int64 `json:"new_role_id"  validate:"required"`
}
