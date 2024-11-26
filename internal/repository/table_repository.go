package repository

import (
	context "context"
	"errors"
	"log"

	entity "github.com/agus-germi/TDL_Dinamita/internal/entity"
	"github.com/jackc/pgx/v5"
)

var ErrTableNotFound = errors.New("table not found")

const (
	qryInsertTable = `INSERT INTO tables (number, seats, location, is_available)
					VALUES ($1, $2, $3, $4)`

	qryDeleteTable = `DELETE FROM tables
					WHERE number=$1`

	qryGetTable = `SELECT number, seats, location, is_available
				   FROM tables
				   WHERE number=$1`
)

func (r *repo) SaveTable(ctx context.Context, tableNumber, seats int64, location string, isAvailable bool) error {
	operation := func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, qryInsertTable, tableNumber, seats, location, isAvailable)
		if err != nil {
			log.Printf("Failed to insert table: %v", err)
			return err
		}

		log.Printf("Table with number %d saved successfully", tableNumber)
		return nil
	}

	return r.executeInTransaction(ctx, operation)
}

func (r *repo) RemoveTable(ctx context.Context, tableNumber int64) error {
	operation := func(tx pgx.Tx) error {
		result, err := r.db.Exec(ctx, qryDeleteTable, tableNumber)
		if err != nil {
			return err
		}

		if result.RowsAffected() == 0 {
			log.Println("Rows affected = 0")
			return ErrTableNotFound
		}

		log.Println("Table removed successfully.")
		return nil
	}

	return r.executeInTransaction(ctx, operation)
}

// We really need this?
func (r *repo) GetTableByNumber(ctx context.Context, tableNumber int64) (*entity.Table, error) {
	table := entity.Table{}

	err := r.db.QueryRow(ctx, qryGetTable, tableNumber).Scan(&table.Number, &table.Seats, &table.Location, &table.IsAvailable)
	if err != nil {
		log.Printf("Failed to execute select query: %v", err)
		return nil, err
	}

	log.Printf("Table retrieved successfully by number: %d", tableNumber)
	return &table, nil
}
