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
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
}

type Repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return &Repo{db: db}
}
