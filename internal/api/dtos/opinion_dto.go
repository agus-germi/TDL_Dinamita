package dtos

type CreateOpinionDTO struct {
    Opinion string `json:"opinion" validate:"required"`
    Rating  int    `json:"rating" validate:"required,min=1,max=5"`
}
