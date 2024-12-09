package repository

import (
	context "context"
	"errors"

	entity "github.com/agus-germi/TDL_Dinamita/internal/entity"
	"github.com/jackc/pgx/v5"
)

var (
	ErrTableNotFound      = errors.New("table not found")
	ErrTableAlreadyExists = errors.New("table already exists")
)

const (
	qryInsertTable = `INSERT INTO tables (number, seats, description)
					VALUES ($1, $2, $3)`

	qryDeleteTable = `DELETE FROM tables
					WHERE id=$1`

	qryGetTableByNumber = `SELECT number, seats, description
				   FROM tables
				   WHERE number=$1`

	qryGetTableByID = `SELECT id, number, seats, description
				   FROM tables
				   WHERE id=$1`

	qryGetAvailableTables = `SELECT id, number, seats, description
				   FROM tables`
)

func (r *repo) SaveTable(ctx context.Context, tableNumber, seats int64, description string) error {
	operation := func(tx pgx.Tx) error {
		table, err := r.getTableByNumber(ctx, tx, tableNumber)
		if table != nil {
			return ErrTableAlreadyExists
		}
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}

		_, err = tx.Exec(ctx, qryInsertTable, tableNumber, seats, description)
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
		table, err := r.getTableByID(ctx, tx, tableID)
		if table == nil {
			r.log.Debugf("Table not found by ID (id=%d)", tableID)
			return ErrTableNotFound
		}
		if err != nil {
			return err
		}

		result, err := tx.Exec(ctx, qryDeleteTable, tableID)
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

func (r *repo) GetAvailableTables(ctx context.Context) (*[]entity.Table, error) {
	var tables []entity.Table

	rows, err := r.db.Query(ctx, qryGetAvailableTables)
	if err != nil {
		r.log.Debugf("Failed to execute select query: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var table entity.Table
		if err := rows.Scan(&table.ID, &table.Number, &table.Seats, &table.Description); err != nil {
			r.log.Debugf("Failed to scan table: %v", err)
			return nil, err
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		r.log.Debugf("Error iterating rows: %v", err)
		return nil, err
	}

	r.log.Infof("Successfully fetched available tables")
	return &tables, nil
}

// Private functions
func (r *repo) getTableByNumber(ctx context.Context, tx pgx.Tx, tableNumber int64) (*entity.Table, error) {
	table := entity.Table{}

	err := tx.QueryRow(ctx, qryGetTableByNumber, tableNumber).Scan(&table.Number, &table.Seats, &table.Location, &table.IsAvailable)
	if err != nil {
		r.log.Debugf("Failed to execute select table (by number) query: %v", err)
		return nil, err
	}

	r.log.Debugf("Table retrieved successfully by number: %d", tableNumber)
	return &table, nil
}

func (r *repo) getTableByID(ctx context.Context, tx pgx.Tx, tableID int64) (*entity.Table, error) {
	table := entity.Table{}

	err := tx.QueryRow(ctx, qryGetTableByID, tableID).Scan(&table.ID, &table.Number, &table.Seats, &table.Location, &table.IsAvailable)
	if err != nil {
		r.log.Debugf("Failed to execute select table (by id) query: %v", err)
		return nil, err
	}

	r.log.Debugf("Table retrieved successfully by ID: %d", tableID)
	return &table, nil
}
