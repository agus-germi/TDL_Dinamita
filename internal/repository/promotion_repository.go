package repository

import (
	"context"
	"errors"

	entity "github.com/agus-germi/TDL_Dinamita/internal/entity"
	"github.com/jackc/pgx/v5"
)

var ErrPromotionNotFound = errors.New("promotion not found")

const (
	qryInsertPromotion     = `INSERT INTO promotions (description, start_date, due_date, discount) VALUES ($1, $2, $3, $4)`
	qryGetPromotionByID    = `SELECT id, description, start_date, due_date, discount FROM promotions WHERE id = $1`
	qryDeletePromotionByID = `DELETE FROM promotions WHERE id = $1`
	qryGetPromotions       = `SELECT id, description, start_date, due_date, discount FROM promotions
                        ORDER BY id`
)

func (r *repo) SavePromotion(ctx context.Context, description string, startDate string, dueDate string, discount int) error {
	r.log.Debugf("Received promotion with description: %s, start_date: %s, due_date: %s, discount: %d",
		description, startDate, dueDate, discount)

	operation := func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, qryInsertPromotion, description, startDate, dueDate, discount)
		if err != nil {
			r.log.Errorf("Failed to insert promotion: %v", err)
			return err
		}

		r.log.Infof("Promotion saved successfully")
		return nil
	}

	return r.executeInTransaction(ctx, operation)
}

func (r *repo) DeletePromotion(ctx context.Context, promotionID int64) error {
	operation := func(tx pgx.Tx) error {
		prom, err := r.getPromotionbyID(ctx, tx, promotionID)
		if prom == nil {
			r.log.Errorf("No promotion found with ID %d", promotionID)
			return ErrPromotionNotFound
		}
		if err != nil {
			r.log.Errorf("Failed to fetch promotion: %v", err)
			return err
		}

		result, err := tx.Exec(ctx, qryDeletePromotionByID, promotionID)
		if err != nil {
			r.log.Errorf("Failed to delete promotion: %v", err)
			return err
		}

		if result.RowsAffected() == 0 {
			r.log.Errorf("No promotion found with ID %d", promotionID)
			return ErrPromotionNotFound
		}

		r.log.Infof("Promotion with ID %d deleted successfully", promotionID)
		return nil
	}

	return r.executeInTransaction(ctx, operation)
}

func (r *repo) GetAllPromotionsAvailable(ctx context.Context) (*[]entity.Promotion, error) {
	rows, err := r.db.Query(ctx, qryGetPromotions)
	if err != nil {
		r.log.Errorf("Failed to execute select promotions query: %v", err)
		return nil, err
	}
	defer rows.Close()

	promotions := []entity.Promotion{}
	for rows.Next() {
		var promotion entity.Promotion
		// Corregir la línea donde se escanean las fechas
		if err := rows.Scan(&promotion.ID, &promotion.Description, &promotion.StartDate, &promotion.DueDate, &promotion.Discount); err != nil {
			r.log.Errorf("Failed to scan row: %v", err)
			return nil, err
		}

		promotions = append(promotions, promotion)
	}

	if rows.Err() != nil {
		r.log.Errorf("Error occurred during row iteration: %v", rows.Err())
		return nil, rows.Err()
	}

	r.log.Info("Promotions retrieved successfully.")
	return &promotions, nil
}

// Private functions
func (r *repo) getPromotionbyID(ctx context.Context, tx pgx.Tx, promotionID int64) (*entity.Promotion, error) {
	promotion := entity.Promotion{}

	err := tx.QueryRow(ctx, qryGetPromotionByID, promotionID).Scan(&promotion.ID, &promotion.Description, &promotion.StartDate, &promotion.DueDate, &promotion.Discount)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		r.log.Errorf("Error fetching promotion by ID: %v", err)
		return nil, err
	}

	r.log.Debugf("Promotion (ID=%d) retrieved succesfully.", promotionID)
	return &promotion, nil
}
