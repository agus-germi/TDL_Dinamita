package service

import (
	context "context"
	"errors"
	time "time"
)

var (
	ErrReservationAlreadyExists = errors.New("reservation already exists")
	ErrRemovingReservation      = errors.New("something went wrong trying to remove a reservation")
)

func (s *serv) RegisterReservation(ctx context.Context, userID int64, name, password, email string, tableNumber int64, date time.Time) error {
	rsv, _ := s.repo.GetReservation(ctx, userID, tableNumber)
	if rsv != nil {
		return ErrReservationAlreadyExists
	}

	return s.repo.SaveReservation(ctx, userID, tableNumber, date)
}

func (s *serv) RemoveReservation(ctx context.Context, userID, tableNumber int64) error {
	err := s.repo.RemoveReservation(ctx, userID, tableNumber)
	if err != nil {
		return ErrRemovingReservation
	}

	return nil
}
