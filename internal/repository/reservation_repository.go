package repository

import (
	context "context"
	"fmt"
	"time"

	"errors"

	"github.com/agus-germi/TDL_Dinamita/internal/entity"
	"github.com/jackc/pgx/v5"
)

var (
	ErrReservationNotFound = errors.New("reservation not found")
	ErrTableNotAvailable   = errors.New("table not available")
)

const (
	invalidTimeSlotID    = -1
	qryInsertReservation = `INSERT INTO reservations (reserved_by, table_number, date, time_slot_id, promotion_id)
							VALUES ($1, $2, $3, $4, $5)`

	qryGetReservationsByUserID = `SELECT r.id, r.reserved_by, r.table_number, r.date, ts.time, (case when p.id = 1 then p.description
																								else p.description || '-' || p.discount || '%'
																								end) promotion
									FROM reservations r
									INNER JOIN time_slots ts ON r.time_slot_id = ts.id
									INNER JOIN promotions p ON r.promotion_id = p.id
									WHERE r.reserved_by=$1
									ORDER BY r.date, ts.time`

	qryGetReservationByID = `SELECT r.id, r.reserved_by, r.table_number, r.date, ts.time, (case when p.id = 1 then p.description
									else p.description || '-' || p.discount || '%'
									end) promotion
								FROM reservations r
								INNER JOIN time_slots ts ON r.time_slot_id = ts.id
								INNER JOIN promotions p ON r.promotion_id = p.id
								WHERE r.id=$1
								ORDER BY r.date, ts.time`

	qryGetReservationForUpdateByID = `SELECT id
							FROM reservations
							WHERE id=$1
							FOR UPDATE`

	qryRemoveReservation = `DELETE FROM reservations
							WHERE id=$1`

	qryGetReservationByTableNumberAndDate = `SELECT id
											FROM reservations
											WHERE table_number=$1 AND date=$2 AND time_slot_id=$3
											FOR UPDATE`

	qryGetTimeSlotID = `SELECT id
						FROM time_slots
						WHERE time=$1`

	qryGetTimeSlots = `SELECT id, time
					  FROM time_slots
					  ORDER BY time`
)

func (r *repo) SaveReservation(ctx context.Context, userID, tableNumber int64, date time.Time, promotionID int) error {
	// El formato RFC3339 es una forma estándar de representar
	// fechas y horas, que es casi equivalente al formato ISO 8601.
	// El estándar RFC 3339 es compatible con ISO 8601, y se utiliza
	// ampliamente en la web, por lo que es una representación adecuada
	// para el formato ISO 8601.
	operation := func(tx pgx.Tx) error {
		timeSlotID, err := r.getTimeSlotIDInTx(ctx, tx, date)
		if err != nil {
			return err
		}

		formattedDate := date.Format("2006-01-02")

		// Check availability
		var existingReservationID int64
		err = tx.QueryRow(ctx, qryGetReservationByTableNumberAndDate, tableNumber, formattedDate, timeSlotID).Scan(&existingReservationID)
		if err == nil {
			// If there's no error, it means that the reservation already exists
			r.log.Debugf("Table %d is already reserved for the date %s and the time %s", tableNumber, formattedDate, date.Format("15:04"))
			return fmt.Errorf("table %d is already reserved for the date %s and the time %s", tableNumber, formattedDate, date.Format("15:04"))
		} else if err != pgx.ErrNoRows {
			r.log.Errorf("Failed to execute select reservation (by table number, date and time slot id): %v", err)
			return fmt.Errorf("error checking existing reservations: %w", err)
		}

		_, err = tx.Exec(ctx, qryInsertReservation, userID, tableNumber, formattedDate, timeSlotID, promotionID)
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
		// lock reservation row in the table
		var existingReservationID int64
		err := tx.QueryRow(ctx, qryGetReservationForUpdateByID, reservationID).Scan(&existingReservationID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				r.log.Errorf("No rows were found by get reservation by ID query: %v", err)
				return ErrReservationNotFound
			}

			r.log.Errorf("Failed to execute select reservation (by ID): %v", err)
			return err
		}

		_, err = tx.Exec(ctx, qryRemoveReservation, reservationID)
		if err != nil {
			r.log.Errorf("Failed to execute delete reservation query: %v", err)
			return err
		}

		r.log.Infof("Reservation (ID: %d) removed successfully.", reservationID)
		return nil
	}

	return r.executeInTransaction(ctx, operation)
}

func (r *repo) GetReservationsByUserID(ctx context.Context, userID int64) (*[]entity.Reservation, error) {
	rows, err := r.db.Query(ctx, qryGetReservationsByUserID, userID)
	if err != nil {
		r.log.Errorf("Failed to execute select reservations (by user id) query: %v", err)
		return nil, err
	}
	defer rows.Close()

	reservations := []entity.Reservation{}
	for rows.Next() {
		var rsv entity.Reservation
		if err := rows.Scan(&rsv.ID, &rsv.UserID, &rsv.TableNumber, &rsv.Date, &rsv.Time, &rsv.Promotion); err != nil {
			r.log.Errorf("Failed to scan row: %v", err)
			return nil, err
		}

		reservations = append(reservations, rsv)
	}

	if rows.Err() != nil {
		r.log.Errorf("Error occurred during row iteration: %v", rows.Err())
		return nil, rows.Err()
	}

	r.log.Debugf("Reservations (of usr_id=%d) retrieved successfully.", userID)
	return &reservations, nil
}

func (r *repo) GetReservationByID(ctx context.Context, reservationID int64) (*entity.Reservation, error) {
	rsv := entity.Reservation{}

	err := r.db.QueryRow(ctx, qryGetReservationByID, reservationID).Scan(&rsv.ID, &rsv.UserID, &rsv.TableNumber, &rsv.Date, &rsv.Time, &rsv.Promotion)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.log.Errorf("No rows were found by get reservation by ID query: %v", err)
			return nil, ErrReservationNotFound
		}

		r.log.Errorf("Failed to execute select reservation (by ID) query: %v", err)
		return nil, err
	}

	r.log.Debugf("Reservation retrieved successfully by ID: %d", reservationID)
	return &rsv, nil
}

func (r *repo) getTimeSlotIDInTx(ctx context.Context, tx pgx.Tx, date time.Time) (int64, error) {
	var timeSlotID int64
	formattedTime := date.Format("15:04")
	err := tx.QueryRow(ctx, qryGetTimeSlotID, formattedTime).Scan(&timeSlotID)
	if err != nil {
		r.log.Errorf("Failed to execute select time slot ID: %v", err)
		return invalidTimeSlotID, err
	}

	r.log.Debugf("Time slot ID (%d) retrieved successfully for time %s.", timeSlotID, formattedTime)
	return timeSlotID, nil
}

func (r *repo) GetTimeSlots(ctx context.Context) (*[]entity.TimeSlot, error) {
	rows, err := r.db.Query(ctx, qryGetTimeSlots)
	if err != nil {
		r.log.Errorf("Failed to execute select time slots query: %v", err)
		return nil, err
	}
	defer rows.Close()

	timeSlots := []entity.TimeSlot{}
	for rows.Next() {
		var ts entity.TimeSlot
		if err := rows.Scan(&ts.ID, &ts.Time); err != nil {
			r.log.Errorf("Failed to scan row: %v", err)
			return nil, err
		}

		timeSlots = append(timeSlots, ts)
	}

	if rows.Err() != nil {
		r.log.Errorf("Error occurred during row iteration: %v", rows.Err())
		return nil, rows.Err()
	}

	r.log.Debugf("Time slots retrieved successfully.")
	return &timeSlots, nil
}
