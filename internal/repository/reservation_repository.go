package repository

import (
	context "context"
	"time"
)

// YYYY-MM-DD HH:MI:SS[.MS] [TZ] --> Formato ISO 8601, correspondiente a las fechas en PostgreSQL. Ej: '2024-10-29 21:45:30 UTC'
func (r *repo) SaveReservation(ctx context.Context, tableNumber, userID int64, date time.Time) error {
	return nil
}

func (r *repo) RemoveReservation(ctx context.Context, userID int64) error {
	return nil
}
