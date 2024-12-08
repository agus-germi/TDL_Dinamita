package repository

import (
    "context"
	"errors"
    "log"

    entity "github.com/agus-germi/TDL_Dinamita/internal/entity"
    "github.com/jackc/pgx/v5"
)

var ErrPromotionNotFound = errors.New("promotion not found")

const (
    qryInsertPromotion    = `INSERT INTO promotions (description, start_date, due_date, discount) VALUES ($1, $2, $3, $4)`
    qryGetPromotionByID   = `SELECT id, description, start_date, due_date, discount FROM promotions WHERE id = $1`
    qryDeletePromotionByID = `DELETE FROM promotions WHERE id = $1`
)

func (r *repo) SavePromotion(ctx context.Context, description string, startDate string, dueDate string, discount int) error {
    log.Printf("Received promotion with description: %s, start_date: %s, due_date: %s, discount: %d", 
        description, startDate, dueDate, discount)

    operation := func(tx pgx.Tx) error {
        _, err := tx.Exec(ctx, qryInsertPromotion, description, startDate, dueDate, discount)
        if err != nil {
            log.Printf("Failed to insert promotion: %v", err)
            return err
        }

        log.Printf("Promotion saved successfully")
        return nil
    }

    return r.executeInTransaction(ctx, operation)
}

func (r *repo) GetPromotionbyID(ctx context.Context, promotionID int64) (*entity.Promotion, error) {
    promotion := entity.Promotion{}
    err := r.db.QueryRow(ctx, qryGetPromotionByID, promotionID).Scan(&promotion.ID, &promotion.Description, &promotion.StartDate, &promotion.DueDate, &promotion.Discount)
    
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, nil
        }
        log.Printf("Error fetching promotion by ID: %v", err)
        return nil, err
    }

    return &promotion, nil
}

func (r *repo) DeletePromotion(ctx context.Context, promotionID int64) error {
    operation := func(tx pgx.Tx) error {
        result, err := r.db.Exec(ctx, qryDeletePromotionByID, promotionID)
        if err != nil {
            log.Printf("Failed to delete promotion: %v", err)
            return err
        }

        if result.RowsAffected() == 0 {
            log.Printf("No promotion found with ID %d", promotionID)
            return ErrPromotionNotFound
        }

        log.Printf("Promotion with ID %d deleted successfully", promotionID)
        return nil
    }

    return r.executeInTransaction(ctx, operation)
}
