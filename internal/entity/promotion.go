package entity

import "time"

type Promotion struct {
    ID          int64     `db:"id"`          
    Description string    `db:"description"` 
    StartDate   time.Time `db:"start_date"`  
    DueDate     time.Time `db:"due_date"`   
    Discount    int       `db:"discount"`  
}