package service

import (
	context "context"
	"errors"

	models "github.com/agus-germi/TDL_Dinamita/internal/models"
	"github.com/agus-germi/TDL_Dinamita/internal/repository"
)

var (
	ErrTableAlreadyExists = errors.New("table already exists")
	ErrTableNotFound      = errors.New("table not found")
	ErrAddingTable        = errors.New("something went wrong trying to add a table")
	ErrRemovingTable      = errors.New("something went wrong trying to remove a table")
)

func (s *serv) AddTable(ctx context.Context, tableNumber, seats int64, location string) error {
	err := s.repo.SaveTable(ctx, tableNumber, seats, location, true) // All tables are added being available
	if errors.Is(err, repository.ErrTableAlreadyExists) {
		return ErrTableAlreadyExists
	}

	return err
}

func (s *serv) RemoveTable(ctx context.Context, tableID int64) error {
	table, _ := s.repo.GetTableByID(ctx, tableID)
	if table == nil {
		return ErrTableNotFound
	}

	err := s.repo.RemoveTable(ctx, tableID)
	if err != nil {
		return ErrRemovingTable
	}

	return nil
}

func (s *serv) GetAvailableTables(ctx context.Context) (*[]models.Table, error) {
	tables, err := s.repo.GetAvailableTables(ctx)
	if err != nil {
		return nil, err
	}

	var modelTables []models.Table
	for _, t := range *tables {
		modelTable := models.Table{
			Number:      t.Number,
			Seats:       t.Seats,
			Location:    t.Location,
			IsAvailable: t.IsAvailable,
		}
		modelTables = append(modelTables, modelTable)
	}

	return &modelTables, nil
}
