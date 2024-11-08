package repository

import (
	context "context"
	"log"
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
							WHERE reserved_by=$1 AND table_number=$2 AND reservation_date=$3`

	qryGetReservationByTableNumberAndDate = `SELECT reserved_by, table_number, reservation_date
											FROM reservations
											WHERE table_number=$1 AND reservation_date=$2`
)

// When we show the reservation date to the user we have to convert
// the date according to the local time zone.
// Keep in mind that the date saved in the DB is according to UTC location.
func (r *repo) SaveReservation(ctx context.Context, userID, tableNumber int64, date time.Time) error {
	// El formato RFC3339 es una forma estándar de representar
	// fechas y horas, que es casi equivalente al formato ISO 8601.
	// El estándar RFC 3339 es compatible con ISO 8601, y se utiliza
	// ampliamente en la web, por lo que es una representación adecuada
	// para el formato ISO 8601.
	formattedDate := date.Format(time.RFC3339)
	_, err := r.db.ExecContext(ctx, qryInsertReservation, userID, tableNumber, formattedDate)
	return err
}

func (r *repo) RemoveReservation(ctx context.Context, userID, tableNumber int64) error {
	_, err := r.db.ExecContext(ctx, qryRemoveReservation, userID, tableNumber)
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

// No estamos teniendo en cuenta la zona horaria.
// Quizas estaria bueno incluir esta info --> Zona horaria de America/Argentina/Buenos_Aires
func (r *repo) GetReservationByTableNumberAndDate(ctx context.Context, tableNumber int64, date time.Time) (*entity.Reservation, error) {
	rsv := &entity.Reservation{}
	formattedDate := date.Format(time.RFC3339)

	err := r.db.GetContext(ctx, rsv, qryGetReservationByTableNumberAndDate, tableNumber, formattedDate)
	if err != nil {
		log.Println("Error trying to get a reservation by number and date")
		return nil, err
	}

	return rsv, nil
}
