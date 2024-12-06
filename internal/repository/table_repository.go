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

	qryGetTableByNumber = `SELECT number, seats, location, is_available
				   FROM tables
				   WHERE number=$1`

	qryGetTableByID = `SELECT id, number, seats, location, is_available
				   FROM tables
				   WHERE id=$1`

	qryGetAvailableTables = `SELECT id, number, seats, location, is_available
				   FROM tables
				   WHERE is_available = TRUE`
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

func (r *repo) GetTableByNumber(ctx context.Context, tableNumber int64) (*entity.Table, error) {
	table := entity.Table{}

	err := r.db.QueryRow(ctx, qryGetTableByNumber, tableNumber).Scan(&table.Number, &table.Seats, &table.Location, &table.IsAvailable)
	if err != nil {
		log.Printf("Failed to execute select query: %v", err)
		return nil, err
	}

	log.Printf("Table retrieved successfully by number: %d", tableNumber)
	return &table, nil
}

func (r *repo) GetTableByID(ctx context.Context, tableID int64) (*entity.Table, error) {
	table := entity.Table{}

	err := r.db.QueryRow(ctx, qryGetTableByID, tableID).Scan(&table.ID, &table.Number, &table.Seats, &table.Location, &table.IsAvailable)
	if err != nil {
		log.Printf("Failed to execute select query: %v", err)
		return nil, err
	}

	log.Printf("Table retrieved successfully by ID: %d", tableID)
	return &table, nil
}

func (r *repo) GetAvailableTables(ctx context.Context) (*[]entity.Table, error) {
    var tables []entity.Table

    rows, err := r.db.Query(ctx, qryGetAvailableTables)
    if err != nil {
        log.Printf("Failed to execute select query: %v", err)
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var table entity.Table
        if err := rows.Scan(&table.ID, &table.Number, &table.Seats, &table.Location, &table.IsAvailable); err != nil {
            log.Printf("Failed to scan table: %v", err)
            return nil, err
        }
        tables = append(tables, table)
    }

    if err := rows.Err(); err != nil {
        log.Printf("Error iterating rows: %v", err)
        return nil, err
    }

    log.Printf("Successfully fetched available tables")
    return &tables, nil
}

