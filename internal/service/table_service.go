package service

import (
	context "context"
	"errors"

	"github.com/agus-germi/TDL_Dinamita/internal/entity"
	models "github.com/agus-germi/TDL_Dinamita/internal/models"
	"github.com/agus-germi/TDL_Dinamita/internal/repository"
)

var (
	ErrTableAlreadyExists = errors.New("table already exists")
	ErrTableNotFound      = errors.New("table not found")
	ErrAddingTable        = errors.New("something went wrong trying to add a table")
	ErrRemovingTable      = errors.New("something went wrong trying to remove a table")
)

func (s *serv) AddTable(ctx context.Context, tableNumber, seats int64, description string) error {
	err := s.executeWithTimeout(ctx, func(ctx context.Context) error {
		return s.repo.SaveTable(ctx, tableNumber, seats, description)
	})

	if errors.Is(err, repository.ErrTableAlreadyExists) {
		return ErrTableAlreadyExists
	}

	return err
}

func (s *serv) RemoveTable(ctx context.Context, tableID int64) error {
	err := s.executeWithTimeout(ctx, func(ctx context.Context) error {
		return s.repo.RemoveTable(ctx, tableID)
	})

	if err != nil {
		if errors.Is(err, repository.ErrTableNotFound) {
			return ErrTableNotFound
		}

		return ErrRemovingTable
	}

	return nil
}

func (s *serv) GetAvailableTables(ctx context.Context) (*[]models.Table, error) {
	var tables *[]entity.Table
	err := s.executeWithTimeout(ctx, func(ctx context.Context) error {
		var err error
		tables, err = s.repo.GetAvailableTables(ctx)
		return err
	})

	if err != nil {
        return nil, err
    }

    var modelTables []models.Table
    for _, t := range *tables {
        modelTable := models.Table{
            Number:    t.Number,
            Seats:     t.Seats,
            Description:  t.Description,
        }
        modelTables = append(modelTables, modelTable)
    }

    return &modelTables, nil
}