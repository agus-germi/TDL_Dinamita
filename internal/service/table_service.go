package service

import (
	context "context"
	"errors"
)

var (
	ErrTableAlreadyExists = errors.New("table already exists")
	ErrRemovingTable      = errors.New("Something went wrong trying to remove a table")
)

func (s *serv) AddTable(ctx context.Context, tableNumber, seats int64, location string) error {
	table, _ := s.repo.GetTableByNumber(ctx, tableNumber)
	if table != nil {
		return ErrTableAlreadyExists
	}

	return s.repo.SaveTable(ctx, tableNumber, seats, location, true) // All tables are added being available
}

func (s *serv) RemoveTable(ctx context.Context, tableNumber int64) error {
	err := s.repo.RemoveTable(ctx, tableNumber)
	if err != nil {
		return ErrRemovingTable
	}

	return nil
}
