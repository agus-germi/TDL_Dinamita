package repository

import (
	"context"
	time "time"

	"github.com/agus-germi/TDL_Dinamita/internal/entity"
	"github.com/jmoiron/sqlx"
)

// Repository is the interface that wraps the basic CRUD operations.
//
//go:generate mockery --name=Repository --output=repository --inpackage
type Repository interface {
	// User
	SaveUser(ctx context.Context, name, passwd, email string) error
	RemoveUser(ctx context.Context, email string) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)

	// UserRole
	SaveUserRole(ctx context.Context, userID, roleID int64) error
	RemoveUserRole(ctx context.Context, userID int64) error
	GetUserRole(ctx context.Context, userID int64) (*entity.UserRole, error)

	// Table
	SaveTable(ctx context.Context, tableNumber, seats int64, location string, isAvailable bool) error
	RemoveTable(ctx context.Context, tableNumber int64) error
	GetTableByNumber(ctx context.Context, tableNumber int64) (*entity.Table, error)
	CheckTableAvailability(ctx context.Context, tableNumber int64, reservationDate time.Time) (bool, error)

	// Reservation
	SaveReservation(ctx context.Context, userID, tableNumber int64, date time.Time) error
	RemoveReservation(ctx context.Context, userID, tableNumber int64) error // Considerar remover una reserva usando su ID.
	GetReservation(ctx context.Context, userID, tableNumber int64) (*entity.Reservation, error)
}

type repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &repo{db: db}
}
