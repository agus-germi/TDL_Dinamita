package repository

import (
	context "context"
	"errors"
	"log"

	entity "github.com/agus-germi/TDL_Dinamita/internal/entity"
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
	_, err := r.db.ExecContext(ctx, qryInsertTable, tableNumber, seats, location, isAvailable)
	return err
}

func (r *repo) RemoveTable(ctx context.Context, tableNumber int64) error {
	result, err := r.db.ExecContext(ctx, qryDeleteTable, tableNumber)
	if err != nil {
		return err
	}

	// Check the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error in checking the number of rows affected")
		return err // Error when determining rows affected
	}

	if rowsAffected == 0 {
		log.Println("Rows affected = 0")
		return ErrTableNotFound
	}

	return nil
}

func (r *repo) GetTableByNumber(ctx context.Context, tableNumber int64) (*entity.Table, error) {
	table := &entity.Table{}

	err := r.db.GetContext(ctx, table, qryGetTable, tableNumber)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}

	return table, nil
}
