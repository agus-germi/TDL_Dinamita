package models

import "time"

type Promotion struct {
    ID          int64     `json:"id"`          
    Description string    `json:"description"` 
    StartDate   time.Time `json:"start_date"` 
    DueDate     time.Time `json:"due_date"`  
    Discount    int64       `json:"discount"` 
}
