package service

import (
	context "context"
	"errors"
)

var (
	ErrTableAlreadyExists = errors.New("table already exists")
	ErrTableNotFound      = errors.New("table not found")
	ErrAddingTable        = errors.New("something went wrong trying to add a table")
	ErrRemovingTable      = errors.New("something went wrong trying to remove a table")
)

func (s *serv) AddTable(ctx context.Context, tableNumber, seats int64, location string) error {
	table, _ := s.repo.GetTableByNumber(ctx, tableNumber)
	if table != nil {
		return ErrTableAlreadyExists
	}

	return s.repo.SaveTable(ctx, tableNumber, seats, location, true) // All tables are added being available
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
