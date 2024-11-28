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
	ErrInvalidPermission  = errors.New("user does not have permission to execute")
)

func (s *serv) AddTable(ctx context.Context, tableNumber, seats int64, location string, email string) error {

	if !s.repo.HasPermission(ctx, email) {
		return ErrInvalidPermission
	}

	table, _ := s.repo.GetTableByNumber(ctx, tableNumber)
	if table != nil {
		return ErrTableAlreadyExists
	}

	return s.repo.SaveTable(ctx, tableNumber, seats, location, true) // All tables are added being available
}

func (s *serv) RemoveTable(ctx context.Context, tableNumber int64, email string) error {

	if !s.repo.HasPermission(ctx, email) {
		return ErrInvalidPermission
	}

	table, _ := s.repo.GetTableByNumber(ctx, tableNumber)
	if table == nil {
		return ErrTableNotFound
	}

	err := s.repo.RemoveTable(ctx, tableNumber)
	if err != nil {
		return ErrRemovingTable
	}

	return nil
}
