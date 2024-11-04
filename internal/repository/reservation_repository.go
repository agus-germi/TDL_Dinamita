package repository

import (
	context "context"
	"fmt"
	"time"

	"github.com/agus-germi/TDL_Dinamita/internal/entity"
)

const (
	qryInsertReservation = `INSERT INTO reservations (reserved_by, table_number, reservation_date) 
							VALUES ($1, $2, $3)`

	qryGetReservation = `SELECT reserved_by, table_number, reservation_date
						 FROM reservations
						 WHERE reserved_by=$1 AND table_number=$2 AND reservation_date=$3`

	qryRemoveReservation = `DELETE FROM reservations
							WHERE reserved_by=$1`
)

// When we show the reservation date to the user we have to convert
// the date according to the local location.
// Keep in mind that the date saved in the DB is according to UTC location.
func (r *repo) SaveReservation(ctx context.Context, userID, tableNumber int64, date time.Time) error {
	formattedDate := r.FormatDate(date)
	_, err := r.db.ExecContext(ctx, qryInsertReservation, userID, tableNumber, formattedDate)
	return err
}

func (r *repo) RemoveReservation(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(ctx, qryRemoveReservation, userID)
	return err
}

func (r *repo) GetReservation(ctx context.Context, userID, tableNumber int64) (*entity.Reservation, error) {
	rsv := &entity.Reservation{}

	err := r.db.GetContext(ctx, rsv, qryGetReservation, userID, tableNumber)
	if err != nil {
		return nil, err
	}

	return rsv, nil
}

// Format date to ISO 8601 (complying PostgreSQL date format)
// YYYY-MM-DD HH:MI:SS[.MS] [TZ]
// Example: '2024-10-29 21:45:30 UTC'
func (r *repo) FormatDate(date time.Time) string {
	year, month, day := date.UTC().Date()
	hour := date.UTC().Hour()
	min := date.UTC().Minute()
	sec := date.UTC().Second()
	return fmt.Sprintf("%d-%d-%d %d:%d:%d UTC", year, month, day, hour, min, sec)
}
