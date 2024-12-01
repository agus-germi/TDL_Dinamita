package repository

import (
	"context"
	"log"
	time "time"

	"github.com/agus-germi/TDL_Dinamita/internal/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository is the interface that wraps the basic CRUD operations.
//
//go:generate mockery --name=Repository --output=repository --inpackage
type Repository interface {
	// User
	SaveUser(ctx context.Context, name, passwd, email string, roleID int64) error
	RemoveUser(ctx context.Context, userID int64) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, userID int64) (*entity.User, error)

	// UserRole
	SaveUpdateUserRole(ctx context.Context, userID, roleID int64) error
	GetUserRole(ctx context.Context, userID int64) (int64, error)

	// Table
	SaveTable(ctx context.Context, tableNumber, seats int64, location string, isAvailable bool) error
	RemoveTable(ctx context.Context, tableNumber int64) error
	GetTableByNumber(ctx context.Context, tableNumber int64) (*entity.Table, error)

	// Reservation
	SaveReservation(ctx context.Context, userID, tableNumber int64, date time.Time) error
	RemoveReservation(ctx context.Context, reservationID int64) error
	GetReservationsByUserID(ctx context.Context, userID int64) (*[]entity.Reservation, error)
	GetReservationByID(ctx context.Context, reservationID int64) (*entity.Reservation, error)
	GetReservationByTableNumberAndDate(ctx context.Context, tableNumber int64, date time.Time) (*entity.Reservation, error) // Este metodo deberia devolver todas las reservas hechas de una mesa en el dia determinado (deberia llamarse GetReservationsByTableNumberAndDate)
	//GetReservationsByTableNumberAndDate(ctx context.Context, tableNumber int64, date time.Time) (*[]entity.Reservation, error) // Este metodo deberia devolver todas las reservas hechas de una mesa en el dia determinado (deberia llamarse GetReservationsByTableNumberAndDate)
	//CheckTableAvailability(ctx context.Context, tableNumber int64, reservationDate time.Time) (bool, error)
}

type repo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) Repository {
	return &repo{db: db}
}

func (r *repo) executeInTransaction(ctx context.Context, operation func(tx pgx.Tx) error) error {
	// Start a transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Printf("Failed to start transaction: %v", err)
		return err
	}

	// Ensure the transaction is committed or rolled back appropriately
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			log.Println("Transaction rollback due to panic")
			panic(p)
		} else if err != nil {
			tx.Rollback(ctx)
			log.Println("Transaction rollback due to error")
		} else {
			err = tx.Commit(ctx)
			if err != nil {
				log.Printf("Failed to commit transaction: %v", err)
			} else {
				log.Println("Transaction committed successfully")
			}
		}
	}()

	// Execute the operation passed in as a function
	return operation(tx)
}
