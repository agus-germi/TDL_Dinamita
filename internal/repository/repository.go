package repository

import (
	"context"
	time "time"

	"github.com/agus-germi/TDL_Dinamita/internal/entity"
	"github.com/agus-germi/TDL_Dinamita/logger"
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
	GetUserByEmail(ctx context.Context, tx pgx.Tx, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, tx pgx.Tx, userID int64) (*entity.User, error)

	// UserRole
	SaveUpdateUserRole(ctx context.Context, userID, roleID int64) error

	// Table
	SaveTable(ctx context.Context, tableNumber, seats int64, description string) error
	RemoveTable(ctx context.Context, tableID int64) error
	GetAvailableTables(ctx context.Context) (*[]entity.Table, error)

	// Reservation
	SaveReservation(ctx context.Context, userID, tableNumber int64, date time.Time, promotionID int) error
	RemoveReservation(ctx context.Context, reservationID int64) error
	GetReservationsByUserID(ctx context.Context, userID int64) (*[]entity.Reservation, error)
	GetReservationByID(ctx context.Context, reservationID int64) (*entity.Reservation, error)

	//Menu
	SaveDish(ctx context.Context, name string, price int64, description string) error
	GetDishByName(ctx context.Context, name string) (*entity.Dish, error)
	GetDishByID(ctx context.Context, dishID int64) (*entity.Dish, error)
	UpdateDish(ctx context.Context, dishID int64, name string, price int64, description string) error
	RemoveDish(ctx context.Context, dishID int64) error
	GetAllDishes(ctx context.Context) (*[]entity.Dish, error)

	//Time slots
	GetTimeSlots(ctx context.Context) (*[]entity.TimeSlot, error)

	// Opinion
	SaveOpinion(ctx context.Context, userID int64, opinion string, rating int) error
	GetAllOpinions(ctx context.Context) (*[]entity.Opinion, error)

	//Promotion
	SavePromotion(ctx context.Context, description string, startDate string, dueDate string, discount int) error
	GetPromotionbyID(ctx context.Context, promotionID int64) (*entity.Promotion, error)
	DeletePromotion(ctx context.Context, promotionID int64) error
	GetAllPromotionsAvailable(ctx context.Context) (*[]entity.Promotion, error)
}

type repo struct {
	db  *pgxpool.Pool
	log logger.Logger
}

func New(db *pgxpool.Pool, log logger.Logger) Repository {
	log.Debugf("Logger has been injected into API")
	return &repo{
		db:  db,
		log: log,
	}
}

func (r *repo) executeInTransaction(ctx context.Context, operation func(tx pgx.Tx) error) error {
	// Start a transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		r.log.Errorf("Failed to start transaction: %v", err)
		return err
	}

	// Ensure the transaction is committed or rolled back appropriately
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			r.log.Panicf("Transaction rollback due to panic: %v", p)
		} else if err != nil {
			r.log.Errorf("Transaction rollback due to error: %v", err)
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
			if err != nil {
				r.log.Errorf("Failed to commit transaction: %v", err)
			} else {
				r.log.Infof("Transaction committed successfully")
			}
		}
	}()

	// Execute the operation passed in as a function
	return operation(tx)
}
