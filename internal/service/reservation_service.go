package service

import (
	context "context"
	"errors"
	"log"
	time "time"
)

var (
	ErrReservationAlreadyExists = errors.New("reservation already exists")
	ErrRemovingReservation      = errors.New("something went wrong trying to remove a reservation")
	ErrCheckingAvailability     = errors.New("error checking table availability")
	ErrTableNotAvailable        = errors.New("table is not available")
)

func (s *serv) RegisterReservation(ctx context.Context, userID int64, name, password, email string, tableNumber int64, date time.Time) error {
	is_available, err := s.repo.CheckTableAvailability(ctx, tableNumber, date)
	if err != nil {
		return ErrCheckingAvailability
	}

	if !is_available {
		log.Println("Availability: ", is_available)
		return ErrTableNotAvailable
	}

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
