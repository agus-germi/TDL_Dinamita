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
	RemoveUser(ctx context.Context, userID int64) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, userID int64) (*entity.User, error)

	// UserRole
	SaveUserRole(ctx context.Context, userID, roleID int64) error
	RemoveUserRole(ctx context.Context, userID int64) error
	GetUserRole(ctx context.Context, userID int64) (*entity.UserRole, error)
	HasPermission(ctx context.Context, userID int64) bool

	// Table
	SaveTable(ctx context.Context, tableNumber, seats int64, location string, isAvailable bool) error
	RemoveTable(ctx context.Context, tableNumber int64) error
	GetTableByNumber(ctx context.Context, tableNumber int64) (*entity.Table, error)

	// Reservation
	SaveReservation(ctx context.Context, userID, tableNumber int64, date time.Time) error
	RemoveReservation(ctx context.Context, reservationID int64) error
	GetReservationsByUserID(ctx context.Context, userID int64) (*[]entity.Reservation, error)
	GetReservationByID(ctx context.Context, reservationID int64) (*entity.Reservation, error)
	GetReservationByTableNumberAndDate(ctx context.Context, tableNumber int64, date time.Time) (*entity.Reservation, error)
}

type repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &repo{db: db}
}
