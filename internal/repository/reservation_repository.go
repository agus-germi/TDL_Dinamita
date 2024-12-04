package repository

import (
	context "context"
	"time"

	"errors"

	"github.com/agus-germi/TDL_Dinamita/internal/entity"
	"github.com/jackc/pgx/v5"
)

var ErrReservationNotFound = errors.New("reservation not found")

const (
	invalidTimeSlotID    = -1
	qryInsertReservation = `INSERT INTO reservations (reserved_by, table_number, date, time_slot_id) 
							VALUES ($1, $2, $3, $4)`

	qryGetReservationsByUserID = `SELECT *
						 FROM reservations
						 WHERE reserved_by=$1`

	qryGetReservationByID = `SELECT *
							FROM reservations
							WHERE id=$1`

	qryRemoveReservation = `DELETE FROM reservations
							WHERE id=$1`

	qryGetReservationByTableNumberAndDate = `SELECT *
											FROM reservations
											WHERE table_number=$1 AND date=$2 AND time_slot_id=$3`

	qryGetTimeSlotID = `SELECT id
						FROM time_slots
						WHERE time=$1`
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
		time_slot_id, err := r.getTimeSlotID(ctx, date)
		if err != nil {
			return err
		}

		formattedDate := date.Format("2006-01-02")

		_, err = tx.Exec(ctx, qryInsertReservation, userID, tableNumber, formattedDate, time_slot_id)
		if err != nil {
			r.log.Errorf("Failed to execute insert reservation query: %v", err)
			return err
		}

		r.log.Infof("Reservation saved successfully.")
		return nil
	}

	return r.executeInTransaction(ctx, operation)
}

func (r *repo) RemoveReservation(ctx context.Context, reservationID int64) error {
	operation := func(tx pgx.Tx) error {
		result, err := tx.Exec(ctx, qryRemoveReservation, reservationID)
		if err != nil {
			r.log.Errorf("Failed to execute delete reservation query: %v", err)
			return err
		}

		if result.RowsAffected() == 0 {
			r.log.Errorf("No rows were affected by the delete query: %v", ErrReservationNotFound)
			return ErrReservationNotFound // Custom error if no rows were deleted (maybe it could be "reservation doesn't exist")
		}

		r.log.Infof("Reservation (%d) removed successfully.", reservationID)
		return nil
	}

	return r.executeInTransaction(ctx, operation)
}

func (r *repo) GetReservationsByUserID(ctx context.Context, userID int64) (*[]entity.Reservation, error) {
	rows, err := r.db.Query(ctx, qryGetReservationsByUserID, userID)
	if err != nil {
		r.log.Errorf("Failed to execute select reservation (by user id) query: %v", err)
		return nil, err
	}
	defer rows.Close()

	reservations := []entity.Reservation{}
	for rows.Next() {
		var rsv entity.Reservation
		if err := rows.Scan(&rsv.ID, &rsv.UserID, &rsv.TableNumber, &rsv.Date, &rsv.TimeSlotID); err != nil {
			r.log.Errorf("Failed to scan row: %v", err)
			return nil, err
		}
		reservations = append(reservations, rsv)
	}

	if rows.Err() != nil {
		r.log.Errorf("Error occurred during row iteration: %v", rows.Err())
		return nil, rows.Err()
	}

	// TODO: buscar el time slot correspondiente a cada reserva (pensar si hay que hacerlo aca o en la capa de servicio)

	r.log.Debugf("Reservations (of usr_id=%d) retrieved successfully.", userID)
	return &reservations, nil
}

func (r *repo) GetReservationByID(ctx context.Context, reservationID int64) (*entity.Reservation, error) {
	rsv := entity.Reservation{}

	err := r.db.QueryRow(ctx, qryGetReservationByID, reservationID).Scan(&rsv.ID, &rsv.UserID, &rsv.TableNumber, &rsv.Date)
	if err != nil {
		r.log.Errorf("Failed to execute select reservation (by ID) query: %v", err)
		return nil, err
	}

	r.log.Debugf("Reservation retrieved successfully by ID: %d", reservationID)
	return &rsv, nil
}

// Este metodo deberia devolver todas las reservas hechas de una mesa en el dia determinado (deberia llamarse GetReservationsByTableNumberAndDate)
// Este metodo hay que modificarlo para que se adecue a la nueva estructura de la tabla de reservas (las tablas que estan en el informe de la semana que le entregamos al profe).
// Basicamente hay que extraer la fecha de Date.
func (r *repo) GetReservationByTableNumberAndDate(ctx context.Context, tableNumber int64, date time.Time) (*entity.Reservation, error) {
	time_slot_id, err := r.getTimeSlotID(ctx, date)
	if err != nil {
		return nil, err
	}

	formattedDate := date.Format("2006-01-02")
	reservation := entity.Reservation{}
	err = r.db.QueryRow(ctx, qryGetReservationByTableNumberAndDate, tableNumber, formattedDate, time_slot_id).Scan(&reservation.ID, &reservation.UserID, &reservation.TableNumber, &reservation.Date, &reservation.TimeSlotID)
	if err != nil {
		r.log.Errorf("Failed to execute select reservation (by table number, date and time slot id): %v", err)
		return nil, err
	}

	r.log.Debugf("Reservation retrieved successfully for table number %d on date %s and time %s.", tableNumber, formattedDate, date.Format("15:04"))
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
		if err := rows.Scan(&rsv.UserID, &rsv.TableNumber, &rsv.Date); err != nil {
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

func (r *repo) getTimeSlotID(ctx context.Context, date time.Time) (int64, error) {
	var time_slot_id int64
	formattedTime := date.Format("15:04")
	err := r.db.QueryRow(ctx, qryGetTimeSlotID, formattedTime).Scan(&time_slot_id)
	if err != nil {
		r.log.Errorf("Failed to execute select time slot id: %v", err)
		return invalidTimeSlotID, err
	}

	r.log.Debugf("Time slot ID (%d) retrieved successfully for time %s.", time_slot_id, formattedTime)
	return time_slot_id, nil
}
