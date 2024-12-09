package service

import (
	context "context"
	"errors"

	"github.com/agus-germi/TDL_Dinamita/internal/entity"
	models "github.com/agus-germi/TDL_Dinamita/internal/models"
	"github.com/agus-germi/TDL_Dinamita/internal/repository"
)

var (
	ErrDishAlreadyExists = errors.New("dish already exists")
	ErrDishNotFound      = errors.New("dish not found")
	ErrRemovingDish      = errors.New("something went wrong trying to remove a dish")
	ErrUpdatingDish      = errors.New("something went wrong trying to update a dish")
)

func (s *serv) AddDishToMenu(ctx context.Context, name string, price int64, description string) error {
	err := s.executeWithTimeout(ctx, s.config.MaxDBOperationDuration, func(ctx context.Context) error {
		return s.repo.SaveDish(ctx, name, price, description)
	})

	if errors.Is(err, repository.ErrDishAlreadyExists) {
		return ErrDishAlreadyExists
	}

	return err
}

func (s *serv) RemoveDish(ctx context.Context, dishID int64) error {
	err := s.executeWithTimeout(ctx, s.config.MaxDBOperationDuration, func(ctx context.Context) error {
		return s.repo.RemoveDish(ctx, dishID)
	})

	if err != nil {
		if errors.Is(err, repository.ErrDishNotFound) {
			return ErrDishNotFound
		}

		return ErrRemovingDish
	}

	return nil
}

func (s *serv) UpdateDish(ctx context.Context, dishID int64, name string, price int64, description string) error {
	err := s.executeWithTimeout(ctx, s.config.MaxDBOperationDuration, func(ctx context.Context) error {
		return s.repo.UpdateDish(ctx, dishID, name, price, description)
	})

	if err != nil {
		if errors.Is(err, repository.ErrDishNotFound) {
			return ErrDishNotFound
		}

		return ErrUpdatingDish
	}

	return nil
}

// get all dishes from the database
func (s *serv) GetDishes(ctx context.Context) (*[]models.Dish, error) {
	var dishes *[]entity.Dish
	err := s.executeWithTimeout(ctx, s.config.MaxDBOperationDuration, func(ctx context.Context) error {
		var err error
		dishes, err = s.repo.GetAllDishes(ctx)
		return err
	})

	if err != nil {
		return nil, err
	}
	if dishes == nil {
		return &[]models.Dish{}, nil
	}

	modelDishes := make([]models.Dish, len(*dishes))
	for i, dish := range *dishes {
		modelDishes[i] = models.Dish(dish)
	}
	return &modelDishes, nil
}
