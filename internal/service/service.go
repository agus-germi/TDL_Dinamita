package service

import (
	"context"
	time "time"

	models "github.com/agus-germi/TDL_Dinamita/internal/models"
	"github.com/agus-germi/TDL_Dinamita/internal/repository"
	"github.com/agus-germi/TDL_Dinamita/logger"
	"github.com/agus-germi/TDL_Dinamita/utils"
)

// Service is the bussiness logic of the application
//
//go:generate mockery --name=Service --output=service --inpackage
type Service interface {
	// Customer features
	RegisterUser(ctx context.Context, name, password, email string) error
	LoginUser(ctx context.Context, email, password string) (*models.User, error)
	RemoveUser(ctx context.Context, userIDToDelete int64) error
	MakeReservation(ctx context.Context, userID, tableNumber int64, date time.Time) error
	CancelReservation(ctx context.Context, reservationID int64) error
	GetReservationsByUserID(ctx context.Context, userID int64) (*[]models.Reservation, error)
	GetReservationByID(ctx context.Context, reservationID int64) (*models.Reservation, error)
	GetDishes(ctx context.Context) (*[]models.Dish, error)
	GetAvailableTables(ctx context.Context) (*[]models.Table, error)
	GetTimeSlots(ctx context.Context) (*[]models.TimeSlot, error)

	// Admin features
	UpdateUserRole(ctx context.Context, userID, newRoleID int64) error
	AddTable(ctx context.Context, tableNumber, seats int64, location string) error
	RemoveTable(ctx context.Context, tableID int64) error
	RemoveDish(ctx context.Context, dishID int64) error
	AddDishToMenu(ctx context.Context, name string, price int64, description string) error
	UpdateDish(ctx context.Context, dishID int64, name string, price int64, description string) error
}

var maxDBOperationsDuration time.Duration

type serv struct {
	repo repository.Repository
	log  logger.Logger
}

func New(repo repository.Repository, log logger.Logger) Service {
	log.Debugf("Logger has been injected into API")
	return &serv{
		repo: repo,
		log:  log,
	}
}

func init() {
	logger.Log.Info("Initializing reservation service")
	logger.Log.Debug("Executing init() function of 'service' package: Loading MAX_DB_OPERATIONS_DURATION from '.env' file")

	maxDBOperationsDurationStr, err := utils.GetEnv("MAX_DB_OPERATIONS_DURATION")
	if err != nil || maxDBOperationsDurationStr == "" {
		logger.Log.Fatalf("MAX_DB_OPERATIONS_DURATION is not set or invalid: %v", err)
	}
	logger.Log.Debugf("Value read from MAX_DB_OPERATIONS_DURATION: %s", maxDBOperationsDurationStr)

	maxDBOperationsDuration, err = time.ParseDuration(maxDBOperationsDurationStr)
	if err != nil {
		logger.Log.Fatalf("Error trying to parse duration of JWT_EXPIRATION_TIME environment variable: %v", err)
	}

	logger.Log.Infof("Maximun DB operations duration loaded successfully from '.env' file.")
}
