package repository

import (
	"context"

	"github.com/agus-germi/TDL_Dinamita/internal/entity"
	"github.com/jmoiron/sqlx"
)

// Repository is the interface that wraps the basic CRUD operations.
//
//go:generate mockery --name=Repository --output=repository --inpackage
type Repository interface {
	SaveUser(ctx context.Context, name, passwd, email string) error
	RemoveUser(ctx context.Context, email string) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)

	SaveUserRole(ctx context.Context, userID, roleID int64) error
	RemoveUserRole(ctx context.Context, userID int64) error
	GetUserRole(ctx context.Context, userID int64) (*entity.UserRole, error)

	SaveTable(ctx context.Context, seats int64, location string, isAvailable bool) error
	RemoveTable(ctx context.Context, tableNumber int64) error
	GetTableByNumber(ctx context.Context, tableNumber int64) (*entity.Table, error)
}

type repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &repo{db: db}
}
