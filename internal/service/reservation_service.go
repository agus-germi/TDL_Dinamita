package service

import (
	context "context"
	"errors"
	time "time"

	models "github.com/agus-germi/TDL_Dinamita/internal/models"
)

var (
	ErrTableNotAvailable    = errors.New("table is not available")
	ErrRemovingReservation  = errors.New("something went wrong trying to remove a reservation")
	ErrCheckingAvailability = errors.New("error checking table availability")
	ErrReservationNotFound  = errors.New("reservation was not found")
)

func (s *serv) MakeReservation(ctx context.Context, userID, tableNumber int64, date time.Time) error {
	rsv, _ := s.repo.GetReservationByTableNumberAndDate(ctx, tableNumber, date)
	if rsv != nil {
		return ErrTableNotAvailable
	}

	return s.repo.SaveReservation(ctx, userID, tableNumber, date)
}

func (s *serv) CancelReservation(ctx context.Context, reservationID int64) error {
	return s.repo.RemoveReservation(ctx, reservationID)
}

func (s *serv) GetReservationsByUserID(ctx context.Context, userID int64) (*[]models.Reservation, error) {
	usr, err := s.repo.GetUserByID(ctx, userID)
	if usr == nil {
		s.log.Error(ErrUserNotFound)
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	entityReservations, err := s.repo.GetReservationsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if entityReservations == nil {
		return &[]models.Reservation{}, nil
	}

	// Conversion de *[]entity.Reservation a *[]models.Reservation
	modelReservations := make([]models.Reservation, len(*entityReservations))
	for i, entityReservation := range *entityReservations {
		modelReservations[i] = models.Reservation{
			ID:              entityReservation.ID,
			UserID:          entityReservation.UserID,
			TableNumber:     entityReservation.TableNumber,
			ReservationDate: entityReservation.Date, // TODO: hay que conseguir el time slot correspondiente a cada reserva y crear un time.Time con la fecha y la hora
		}
	}

	return &modelReservations, nil
}

func (s *serv) GetReservationByID(ctx context.Context, reservationID int64) (*models.Reservation, error) {
	entityReservation, err := s.repo.GetReservationByID(ctx, reservationID)
	if err != nil {
		return nil, err
	}
	if entityReservation == nil {
		return nil, ErrReservationNotFound
	}

	modelReservation := models.Reservation{
		ID:              entityReservation.ID,
		UserID:          entityReservation.UserID,
		TableNumber:     entityReservation.TableNumber,
		ReservationDate: entityReservation.Date, // TODO: hay que conseguir el time slot correspondiente a cada reserva y crear un time.Time con la fecha y la hora
	}

	return &modelReservation, nil
}
