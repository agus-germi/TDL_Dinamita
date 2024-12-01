package service

import (
	"context"
	time "time"

	models "github.com/agus-germi/TDL_Dinamita/internal/models"
	"github.com/agus-germi/TDL_Dinamita/internal/repository"
)

// Service is the bussiness logic of the application
//
//go:generate mockery --name=Service --output=service --inpackage
type Service interface {
	// Customer features
	RegisterUser(ctx context.Context, name, password, email string) error
	LoginUser(ctx context.Context, email, password string) (*models.User, error)
	RemoveUser(ctx context.Context, userID int64) error
	RegisterReservation(ctx context.Context, userID, tableNumber int64, date time.Time) error
	RemoveReservation(ctx context.Context, reservationID int64) error
	GetReservationsByUserID(ctx context.Context, userID int64) (*[]models.Reservation, error)

	// Admin features
	UpdateUserRole(ctx context.Context, userID, newRoleID int64, clientRoleID int64) error
	AddTable(ctx context.Context, tableNumber, seats int64, location string, clientRoleID int64) error
	RemoveTable(ctx context.Context, tableNumber int64, clientRoleID int64) error
}

type serv struct {
	repo repository.Repository
}

func New(repo repository.Repository) Service {
	return &serv{repo: repo}
}

var adminRoleID int64 = 1
