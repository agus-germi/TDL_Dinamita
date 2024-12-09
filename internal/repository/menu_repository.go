package repository

import (
	context "context"
	"errors"

	"github.com/agus-germi/TDL_Dinamita/internal/entity"
	"github.com/jackc/pgx/v5"
)

var (
	ErrDishAlreadyExists = errors.New("dish already exists")
	ErrDishNotFound      = errors.New("dish not found")
)

const (
	qryInsertDish = `INSERT INTO dishes (name, price, description) VALUES ($1, $2, $3)`

	qryGetDish = `SELECT name, price, description FROM dishes WHERE name=$1`

	qryGetDishByID = `SELECT name, price, description FROM dishes WHERE id=$1`

	qryDeleteDish = `DELETE FROM dishes WHERE id=$1`

	qryGetDishes = `SELECT * FROM dishes`

	qryUpdateDish = `UPDATE dishes SET name=$1, price=$2, description=$3 WHERE id=$4`
)

func (r *repo) SaveDish(ctx context.Context, name string, price int64, description string) error {
	operation := func(tx pgx.Tx) error {
		dish, err := r.getDishByName(ctx, tx, name)
		if dish != nil {
			r.log.Errorf("Dish with name %s already exists", name)
			return ErrDishAlreadyExists
		}
		if err != nil && err != pgx.ErrNoRows {
			r.log.Errorf("Failed to execute select dish query: %v", err)
			return err
		}

		_, err = tx.Exec(ctx, qryInsertDish, name, price, description)
		if err != nil {
			r.log.Errorf("Failed to insert dish: %v", err)
			return err
		}

		r.log.Debugf("Dish with name %s saved successfully", name)
		return nil
	}

	return r.executeInTransaction(ctx, operation)
}

func (r *repo) RemoveDish(ctx context.Context, dishID int64) error {
	operation := func(tx pgx.Tx) error {
		dish, err := r.getDishByID(ctx, tx, dishID)
		if dish != nil {
			r.log.Errorf("Dish with ID %d not found", dishID)
			return ErrDishNotFound
		}
		if err != nil {
			return err
		}

		result, err := tx.Exec(ctx, qryDeleteDish, dishID)
		if err != nil {
			return err
		}

		if result.RowsAffected() == 0 {
			r.log.Errorf("Rows affected = 0")
			return ErrDishNotFound
		}

		r.log.Infof("Dish (ID=%d) removed successfully.", dishID)
		return nil
	}

	return r.executeInTransaction(ctx, operation)
}

func (r *repo) UpdateDish(ctx context.Context, dishID int64, name string, price int64, description string) error {
	operation := func(tx pgx.Tx) error {
		dish, err := r.getDishByID(ctx, tx, dishID)
		if dish != nil {
			r.log.Errorf("Dish with ID %d not found", dishID)
			return ErrDishNotFound
		}
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, qryUpdateDish, name, price, description, dishID)
		if err != nil {
			r.log.Errorf("Failed to execute update user role query: %v", err)
			return err
		}

		r.log.Infof("Dish (id=%d) updated successfully.", dishID)
		return nil
	}

	return r.executeInTransaction(ctx, operation)
}

func (r *repo) GetAllDishes(ctx context.Context) (*[]entity.Dish, error) {
	rows, err := r.db.Query(ctx, qryGetDishes)
	if err != nil {
		r.log.Errorf("Failed to execute select reservations (by user id) query: %v", err)
		return nil, err
	}
	defer rows.Close()

	dishes := []entity.Dish{}
	for rows.Next() {
		var dish entity.Dish
		if err := rows.Scan(&dish.ID, &dish.Name, &dish.Price, &dish.Description); err != nil {
			r.log.Errorf("Failed to scan row: %v", err)
			return nil, err
		}

		dishes = append(dishes, dish)
	}

	if rows.Err() != nil {
		r.log.Errorf("Error occurred during row iteration: %v", rows.Err())
		return nil, rows.Err()
	}

	r.log.Debugf("Dishes retrieved successfully.")
	return &dishes, nil
}

// Private funcions
func (r *repo) getDishByID(ctx context.Context, tx pgx.Tx, dishID int64) (*entity.Dish, error) {
	dish := entity.Dish{}

	err := tx.QueryRow(ctx, qryGetDishByID, dishID).Scan(&dish.Name, &dish.Price, &dish.Description)
	if err != nil {
		r.log.Errorf("Failed to execute select query: %v", err)
		return nil, err
	}

	r.log.Debugf("Dish retrieved successfully by ID: %d", dishID)
	return &dish, nil
}

func (r *repo) getDishByName(ctx context.Context, tx pgx.Tx, name string) (*entity.Dish, error) {
	dish := entity.Dish{}

	err := tx.QueryRow(ctx, qryGetDish, name).Scan(&dish.Name, &dish.Price, &dish.Description)
	if err != nil {
		r.log.Errorf("Failed to execute select dish (by name) query: %v", err)
		return nil, err
	}

	r.log.Debugf("Dish retrieved successfully by name: %s", name)
	return &dish, nil
}
