package repository

import (
	context "context"
	"log"

	entity "github.com/agus-germi/TDL_Dinamita/internal/entity"
)

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
	_, err := r.db.ExecContext(ctx, qryDeleteTable, tableNumber)
	return err
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
