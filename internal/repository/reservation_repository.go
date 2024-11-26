package repository

import (
	context "context"
	"log"
	"time"

	"errors"

	"github.com/agus-germi/TDL_Dinamita/internal/entity"
	"github.com/jackc/pgx/v5"
)

var ErrReservationNotFound = errors.New("reservation not found")

const (
	qryInsertReservation = `INSERT INTO reservations (reserved_by, table_number, reservation_date) 
							VALUES ($1, $2, $3)`

	qryGetReservationsByUserID = `SELECT id, reserved_by, table_number, reservation_date
						 FROM reservations
						 WHERE reserved_by=$1`

	qryGetReservationByID = `SELECT id, reserved_by, table_number, reservation_date
							FROM reservations
							WHERE id=$1`

	qryRemoveReservation = `DELETE FROM reservations
							WHERE id=$1`

	qryGetReservationByTableNumberAndDate = `SELECT reserved_by, table_number, reservation_date
											FROM reservations
											WHERE table_number=$1 AND reservation_date=$2`
)

// TODO: Investigate how to use SELECT ... FOR UPDATE (to lock the row in the table) and UPDATE ... SET ... WHERE ... (to update the row in the table)
// Maybe we can levarage the pgxpool.Tx to do this kind of operations, and we can use the pgxpool.Tx.ExecContext to execute the queries.
// Moreover, we can use SELECT ... FOR UPDATE to lock the corresponding row in the availability table (that we disscussed to create)
// Same observation applies to all repository method that needs this kind of control.

// When we show the reservation date to the user we have to convert
// the date according to the local time zone.
// Keep in mind that the date saved in the DB is according to UTC location.
func (r *repo) SaveReservation(ctx context.Context, userID, tableNumber int64, date time.Time) error {
	// El formato RFC3339 es una forma estándar de representar
	// fechas y horas, que es casi equivalente al formato ISO 8601.
	// El estándar RFC 3339 es compatible con ISO 8601, y se utiliza
	// ampliamente en la web, por lo que es una representación adecuada
	// para el formato ISO 8601.
	operation := func(tx pgx.Tx) error {
		formattedDate := date.Format(time.RFC3339) // TODO: According to the new table structure we need to get only the Date (excluding the time)
		_, err := tx.Exec(ctx, qryInsertReservation, userID, tableNumber, formattedDate)
		if err != nil {
			log.Printf("Failed to execute insert query: %v", err)
			return err
		}
		log.Println("Reservation saved successfully.")
		return nil
	}

	return r.executeInTransaction(ctx, operation)
}

func (r *repo) RemoveReservation(ctx context.Context, reservationID int64) error {
	operation := func(tx pgx.Tx) error {
		result, err := tx.Exec(ctx, qryRemoveReservation, reservationID)
		if err != nil {
			log.Printf("Failed to execute delete query: %v", err)
			return err
		}

		if result.RowsAffected() == 0 {
			log.Println("No rows were affected by the delete query.")
			return ErrReservationNotFound // Custom error if no rows were deleted (maybe it could be "reservation doesn't exist")
		}

		log.Println("Reservation removed successfully.")
		return nil
	}

	return r.executeInTransaction(ctx, operation)
}

func (r *repo) GetReservationsByUserID(ctx context.Context, userID int64) (*[]entity.Reservation, error) {
	rows, err := r.db.Query(ctx, qryGetReservationsByUserID, userID)
	if err != nil {
		log.Printf("Failed to execute select query: %v", err)
		return nil, err
	}
	defer rows.Close()

	reservations := []entity.Reservation{}
	for rows.Next() {
		var rsv entity.Reservation
		if err := rows.Scan(&rsv.ID, &rsv.UserID, &rsv.TableNumber, &rsv.ReservationDate); err != nil {
			log.Printf("Failed to scan row: %v", err)
			return nil, err
		}
		reservations = append(reservations, rsv)
	}

	if rows.Err() != nil {
		log.Printf("Error occurred during row iteration: %v", rows.Err())
		return nil, rows.Err()
	}

	log.Println("Reservations retrieved successfully.")
	return &reservations, nil
}

func (r *repo) GetReservationByID(ctx context.Context, reservationID int64) (*entity.Reservation, error) {
	rsv := entity.Reservation{}

	err := r.db.QueryRow(ctx, qryGetReservationByID, reservationID).Scan(&rsv.ID, &rsv.UserID, &rsv.TableNumber, &rsv.ReservationDate)
	if err != nil {
		log.Printf("Failed to execute select query: %v", err)
		return nil, err
	}

	log.Printf("Reservation retrieved successfully by ID: %d", reservationID)
	return &rsv, nil
}

// Este metodo deberia devolver todas las reservas hechas de una mesa en el dia determinado (deberia llamarse GetReservationsByTableNumberAndDate)
// Este metodo hay que modificarlo para que se adecue a la nueva estructura de la tabla de reservas (las tablas que estan en el informe de la semana que le entregamos al profe).
// Basicamente hay que extraer la fecha de Date.
func (r *repo) GetReservationByTableNumberAndDate(ctx context.Context, tableNumber int64, date time.Time) (*entity.Reservation, error) {
	formattedDate := date.Format(time.RFC3339)
	reservation := entity.Reservation{}
	err := r.db.QueryRow(ctx, qryGetReservationByTableNumberAndDate, tableNumber, formattedDate).Scan(&reservation.UserID, &reservation.TableNumber, &reservation.ReservationDate)
	if err != nil {
		log.Printf("Failed to execute select query: %v", err)
		return nil, err
	}

	log.Printf("Reservation retrieved successfully for table number %d on date %s.", tableNumber, formattedDate)
	return &reservation, nil
}

// Este metodo deberia devolver todas las reservas hechas de una mesa en el dia determinado (deberia llamarse GetReservationsByTableNumberAndDate)
// Este metodo hay que modificarlo para que se adecue a la nueva estructura de la tabla de reservas (las tablas que estan en el informe de la semana que le entregamos al profe).
// Basicamente hay que extraer la fecha de Date.
/*
func (r *repo) GetReservationsByTableNumberAndDate(ctx context.Context, tableNumber int64, date time.Time) (*[]entity.Reservation, error) {
	formattedDate := date.Format(time.RFC3339)
	rows, err := r.db.Query(ctx, qryGetReservationByTableNumberAndDate, tableNumber, formattedDate)
	if err != nil {
		log.Printf("Failed to execute select query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var reservations []entity.Reservation
	for rows.Next() {
		var rsv entity.Reservation
		if err := rows.Scan(&rsv.UserID, &rsv.TableNumber, &rsv.ReservationDate); err != nil {
			log.Printf("Failed to scan row: %v", err)
			return nil, err
		}
		reservations = append(reservations, rsv)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error occurred during row iteration: %v", err)
		return nil, err
	}

	log.Printf("Reservations retrieved successfully for table number %d on date %s.", tableNumber, formattedDate)
	return &reservations, nil
}
*/
