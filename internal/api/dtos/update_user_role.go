package dtos

type UpdateUserRoleDTO struct {
	NewRoleID int64 `json:"new_role_id"  validate:"required"`
}
