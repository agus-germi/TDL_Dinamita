package service

import (
	context "context"
	"errors"
	time "time"
)

var (
	ErrTableNotAvailable    = errors.New("table is not available")
	ErrRemovingReservation  = errors.New("something went wrong trying to remove a reservation")
	ErrCheckingAvailability = errors.New("error checking table availability")
)

func (s *serv) RegisterReservation(ctx context.Context, userID int64, name, password, email string, tableNumber int64, date time.Time) error {
	rsv, _ := s.repo.GetReservationByTableNumberAndDate(ctx, tableNumber, date)
	if rsv != nil {
		return ErrTableNotAvailable
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
