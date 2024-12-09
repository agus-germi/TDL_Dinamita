package service

import (
	context "context"
	"errors"
	time "time"

	"github.com/agus-germi/TDL_Dinamita/internal/entity"
	models "github.com/agus-germi/TDL_Dinamita/internal/models"
	"github.com/agus-germi/TDL_Dinamita/internal/repository"
)

var (
	ErrTableNotAvailable    = errors.New("table is not available")
	ErrRemovingReservation  = errors.New("something went wrong trying to remove a reservation")
	ErrCheckingAvailability = errors.New("error checking table availability")
	ErrReservationNotFound  = errors.New("reservation was not found")
)

func (s *serv) MakeReservation(ctx context.Context, userID, tableNumber int64, date time.Time, promotionID int) error {
	err := s.executeWithTimeout(ctx, s.config.MaxDBOperationDuration, func(ctx context.Context) error {
		return s.repo.SaveReservation(ctx, userID, tableNumber, date, promotionID)
	})

	if errors.Is(err, repository.ErrTableNotAvailable) {
		return ErrTableNotAvailable
	}

	return err
}

func (s *serv) CancelReservation(ctx context.Context, reservationID int64) error {
	err := s.executeWithTimeout(ctx, s.config.MaxDBOperationDuration, func(ctx context.Context) error {
		return s.repo.RemoveReservation(ctx, reservationID)
	})

	if errors.Is(err, repository.ErrReservationNotFound) {
		return ErrReservationNotFound
	}

	return err
}

func (s *serv) GetReservationsByUserID(ctx context.Context, userID int64) (*[]models.Reservation, error) {
	var usr *entity.User
	err := s.executeWithTimeout(ctx, s.config.MaxDBOperationDuration, func(ctx context.Context) error {
		var err error
		usr, err = s.repo.GetUserByID(ctx, userID)
		return err
	})

	if usr == nil {
		s.log.Errorf("Failed to search a user:", ErrUserNotFound)
		return nil, ErrUserNotFound
	}
	if err != nil {
		s.log.Errorf("Failed to search a user:", err)
		return nil, err
	}

	var entityReservations *[]entity.Reservation
	err = s.executeWithTimeout(ctx, s.config.MaxDBOperationDuration, func(ctx context.Context) error {
		var err error
		entityReservations, err = s.repo.GetReservationsByUserID(ctx, userID)
		return err
	})

	if err != nil {
		return nil, err
	}

	if entityReservations == nil {
		return &[]models.Reservation{}, nil
	}

	// Conversion de *[]entity.Reservation a *[]models.Reservation
	modelReservations := make([]models.Reservation, len(*entityReservations))
	for i, entityReservation := range *entityReservations {
		reservationDateTime, err := s.combineDateTime(entityReservation.Date, entityReservation.Time)
		if err != nil {
			s.log.Errorf("Failed to combine date and time: %v", err)
			return &[]models.Reservation{}, err
		}

		modelReservations[i] = models.Reservation{
			ID:              entityReservation.ID,
			UserID:          entityReservation.UserID,
			TableNumber:     entityReservation.TableNumber,
			ReservationDate: reservationDateTime,
			Promotion:       entityReservation.Promotion,
		}
	}

	return &modelReservations, nil
}

func (s *serv) GetReservationByID(ctx context.Context, reservationID int64) (*models.Reservation, error) {
	var entityReservation *entity.Reservation
	err := s.executeWithTimeout(ctx, s.config.MaxDBOperationDuration, func(ctx context.Context) error {
		var err error
		entityReservation, err = s.repo.GetReservationByID(ctx, reservationID)
		return err
	})

	if err != nil {
		return nil, err
	}
	if entityReservation == nil {
		return nil, ErrReservationNotFound
	}

	reservationDateTime, err := s.combineDateTime(entityReservation.Date, entityReservation.Time)
	if err != nil {
		s.log.Errorf("Failed to combine date and time: %v", err)
		return &models.Reservation{}, err
	}

	modelReservation := models.Reservation{
		ID:              entityReservation.ID,
		UserID:          entityReservation.UserID,
		TableNumber:     entityReservation.TableNumber,
		ReservationDate: reservationDateTime,
	}

	return &modelReservation, nil
}

func (s *serv) GetTimeSlots(ctx context.Context) (*[]models.TimeSlot, error) {
	var entityTimeSlots *[]entity.TimeSlot
	err := s.executeWithTimeout(ctx, s.config.MaxDBOperationDuration, func(ctx context.Context) error {
		var err error
		entityTimeSlots, err = s.repo.GetTimeSlots(ctx)
		return err
	})

	if err != nil {
		s.log.Errorf("Failed to get time slots: %v", err)
		return nil, err
	}

	if entityTimeSlots == nil {
		return &[]models.TimeSlot{}, nil
	}

	modelTimeSlots := make([]models.TimeSlot, len(*entityTimeSlots))
	for i, entityTimeSlot := range *entityTimeSlots {
		modelTimeSlots[i] = models.TimeSlot{
			ID:   entityTimeSlot.ID,
			Time: entityTimeSlot.Time.Format("15:04:05"),
		}
	}

	return &modelTimeSlots, nil
}

// Combine date (format YYYY-MM-DD) and _time (format HH:mm) in a single string that comply ISO 8601
func (s *serv) combineDateTime(date time.Time, _time string) (time.Time, error) {
	// Determinar el formato de _time
	var timeLayout string
	if len(_time) > 5 { // Si hay más de "HH:mm", asumir que hay segundos o milisegundos
		timeLayout = "15:04:05.000000"
	} else {
		timeLayout = "15:04"
	}

	// Parsear la hora
	parsedTime, err := time.Parse(timeLayout, _time)
	if err != nil {
		s.log.Errorf("Failed to parse time: %v", err)
		return time.Time{}, err
	}

	// Combinar la fecha con la hora
	combinedDateTime := time.Date(
		date.Year(), date.Month(), date.Day(),
		parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), 0, time.UTC,
	)

	s.log.Debugf("Combined date and time: %v", combinedDateTime)
	return combinedDateTime, nil
}
