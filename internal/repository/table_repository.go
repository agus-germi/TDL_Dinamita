package repository

import (
	context "context"

	entity "github.com/agus-germi/TDL_Dinamita/internal/entity"
)

func (r *repo) SaveTable(ctx context.Context, tableNumber, seats int64, location string, isAvailable bool) error {
	return nil
}

func (r *repo) RemoveTable(ctx context.Context, tableNumber int64) error {
	return nil
}

func (r *repo) GetTableByNumber(ctx context.Context, tableNumber int64) (*entity.Table, error) {
	return nil, nil
}
