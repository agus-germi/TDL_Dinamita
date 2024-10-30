package service

import (
	"context"

	"github.com/agus-germi/TDL_Dinamita/internal/dtos"
	"github.com/agus-germi/TDL_Dinamita/internal/repository"
)

// Service is the bussiness logic of the application
//
//go:generate mockery --name=Service --output=service --inpackage
type Service interface {
	// Customer features
	RegisterUser(ctx context.Context, name, password, email string) error
	LoginUser(ctx context.Context, email, password string) (*dtos.User, error)
	RemoveUser(ctx context.Context, email string) error

	// Admin features
	AddUserRole(ctx context.Context, userID, roleID int64) error
	RemoveUserRole(ctx context.Context, userID int64) error
	AddTable(ctx context.Context, tableNumber, seats int64, location string) error //All tables are added being available
	RemoveTable(ctx context.Context, tableNumber int64) error
}

type serv struct {
	repo repository.Repository
}

func New(repo repository.Repository) Service {
	return &serv{repo: repo}
}
