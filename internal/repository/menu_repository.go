package repository

import (
	context "context"
	"log"

	"github.com/agus-germi/TDL_Dinamita/internal/entity"
	"github.com/jackc/pgx/v5"
)

const (
	qryInsertDish = `INSERT INTO dishes (name, price, description) VALUES ($1, $2, $3)`
	qryGetDish    = `SELECT name, price, description FROM dishes WHERE name=$1`
)

func (r *repo) SaveDish(ctx context.Context, name string, price int64, description string) error {
	operation := func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, qryInsertDish, name, price, description)
		if err != nil {
			log.Printf("Failed to insert dish: %v", err)
			return err
		}

		log.Printf("Dish with name %s saved successfully", name)
		return nil
	}

	return r.executeInTransaction(ctx, operation)
}

func (r *repo) GetDishByName(ctx context.Context, name string) (*entity.Dish, error) {
	dish := entity.Dish{}

	err := r.db.QueryRow(ctx, qryGetDish, name).Scan(&dish.Name, &dish.Price, &dish.Description)
	if err != nil {
		log.Printf("Failed to execute select query: %v", err)
		return nil, err
	}

	log.Printf("Dish retrieved successfully by name: %s", name)
	return &dish, nil
}
