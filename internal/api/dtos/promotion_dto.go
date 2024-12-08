package dtos

type CreatePromotionDTO struct {
    Description string `json:"description" validate:"required,min=1,max=255"`
    StartDate   string `json:"start_date" validate:"required,datetime=2006-01-02T15:04:05Z""`
    DueDate     string `json:"due_date" validate:"required,datetime=2006-01-02T15:04:05Z",gtfield=StartDate"`
    Discount    int    `json:"discount" validate:"required,min=0,max=100"`
}
