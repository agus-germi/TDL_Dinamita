package repository

import (
	context "context"
	"log"
	"time"

	"errors"

	"github.com/agus-germi/TDL_Dinamita/internal/entity"
)

var ErrReservationNotFound = errors.New("reservation not found")

const (
	qryInsertReservation = `INSERT INTO reservations (reserved_by, table_number, reservation_date) 
							VALUES ($1, $2, $3)`

	qryGetReservation = `SELECT reserved_by, table_number, reservation_date
						 FROM reservations
						 WHERE reserved_by=$1 AND table_number=$2 AND reservation_date=$3`

	qryGetReservationByID = `SELECT id, reserved_by, table_number, reservation_date
							FROM reservations
							WHERE id=$1`

	qryRemoveReservation = `DELETE FROM reservations
							WHERE id=$1`

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

func (r *repo) RemoveReservation(ctx context.Context, reservationID int64) error {
	result, err := r.db.ExecContext(ctx, qryRemoveReservation, reservationID)
	if err != nil {
		return err // Return the error from the query
	}

	// Check the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error in checking the number of rows affected")
		return err // Error when determining rows affected
	}

	if rowsAffected == 0 {
		log.Println("Rows affected = 0")
		return ErrReservationNotFound // Custom error if no rows were deleted (maybe it could be "reservation doesn't exist")
	}

	return nil
}

func (r *repo) GetReservation(ctx context.Context, userID, tableNumber int64) (*entity.Reservation, error) {
	rsv := &entity.Reservation{}

	err := r.db.GetContext(ctx, rsv, qryGetReservation, userID, tableNumber)
	if err != nil {
		return nil, err
	}

	return rsv, nil
}

<<<<<<< HEAD
// No estamos teniendo en cuenta la zona horaria.
// Quizas estaria bueno incluir esta info --> Zona horaria de America/Argentina/Buenos_Aires
=======
func (r *repo) GetReservationByID(ctx context.Context, reservationID int64) (*entity.Reservation, error) {
	rsv := &entity.Reservation{}

	err := r.db.GetContext(ctx, rsv, qryGetReservationByID, reservationID)
	if err != nil {
		log.Println("Error trying to get a reservation by ID")
		return nil, err
	}

	return rsv, nil
}

>>>>>>> 38f5ac2ef8c7cb52b4beb7ee43c632c3919484c7
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
