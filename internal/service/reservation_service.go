package service

import (
	context "context"
	"time"
)

func (s *serv) ReserveTable(ctx context.Context, tableNumber, userID int64, date time.Time) error {
	return nil
}

func (s *serv) RemoveReservation(ctx context.Context, userID int64) error {
	return nil
}
