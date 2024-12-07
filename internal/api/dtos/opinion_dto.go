package dtos

type CreateOpinionDTO struct {
    UserID  int64  `json:"user_id" validate:"required"`
    Opinion string `json:"opinion" validate:"required"`
}