package repository

import (
	context "context"
	"errors"

	entity "github.com/agus-germi/TDL_Dinamita/internal/entity"
	"github.com/jackc/pgx/v5"
)

var ErrTableNotFound = errors.New("table not found")

const (
	qryInsertTable = `INSERT INTO tables (number, seats, location, is_available)
					VALUES ($1, $2, $3, $4)`

	qryDeleteTable = `DELETE FROM tables
					WHERE id=$1`

	qryGetTableByNumber = `SELECT number, seats, location, is_available
				   FROM tables
				   WHERE number=$1`

	qryGetTableByID = `SELECT id, number, seats, location, is_available
				   FROM tables
				   WHERE id=$1`
)

func (r *repo) SaveTable(ctx context.Context, tableNumber, seats int64, location string, isAvailable bool) error {
	operation := func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, qryInsertTable, tableNumber, seats, location, isAvailable)
		if err != nil {
			r.log.Debugf("Failed to insert table: %v", err)
			return err
		}

		r.log.Infof("Table with number %d saved successfully", tableNumber)
		return nil
	}

	return r.executeInTransaction(ctx, operation)
}

func (r *repo) RemoveTable(ctx context.Context, tableID int64) error {
	operation := func(tx pgx.Tx) error {
		result, err := r.db.Exec(ctx, qryDeleteTable, tableID)
		if err != nil {
			return err
		}

		if result.RowsAffected() == 0 {
			r.log.Debugf("Rows affected = 0")
			return ErrTableNotFound
		}

		r.log.Infof("Table (id=%d) removed successfully.", tableID)
		return nil
	}

	return r.executeInTransaction(ctx, operation)
}

func (r *repo) GetTableByNumber(ctx context.Context, tableNumber int64) (*entity.Table, error) {
	table := entity.Table{}

	err := r.db.QueryRow(ctx, qryGetTableByNumber, tableNumber).Scan(&table.Number, &table.Seats, &table.Location, &table.IsAvailable)
	if err != nil {
		r.log.Debugf("Failed to execute select table (by number) query: %v", err)
		return nil, err
	}

	r.log.Debugf("Table retrieved successfully by number: %d", tableNumber)
	return &table, nil
}

func (r *repo) GetTableByID(ctx context.Context, tableID int64) (*entity.Table, error) {
	table := entity.Table{}

	err := r.db.QueryRow(ctx, qryGetTableByID, tableID).Scan(&table.ID, &table.Number, &table.Seats, &table.Location, &table.IsAvailable)
	if err != nil {
		r.log.Debugf("Failed to execute select table (by id) query: %v", err)
		return nil, err
	}

	r.log.Debugf("Table retrieved successfully by ID: %d", tableID)
	return &table, nil
}
