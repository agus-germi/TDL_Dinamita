package service

import (
	context "context"
	"errors"

	models "github.com/agus-germi/TDL_Dinamita/internal/models"
)

var (
	ErrDishAlreadyExists = errors.New("dish already exists")
	ErrDishNotFound      = errors.New("dish not found")
	ErrRemovingDish      = errors.New("something went wrong trying to remove a dish")
)

func (s *serv) AddDishToMenu(ctx context.Context, name string, price int64, description string) error {
	dish, _ := s.repo.GetDishByName(ctx, name)
	if dish != nil {
		return ErrDishAlreadyExists
	}

	return s.repo.SaveDish(ctx, name, price, description)
}

func (s *serv) RemoveDish(ctx context.Context, dishID int64) error {

	dish, _ := s.repo.GetDishByID(ctx, dishID)
	if dish == nil {
		return ErrDishNotFound
	}

	err := s.repo.RemoveDish(ctx, dishID)
	if err != nil {
		return ErrRemovingDish
	}

	return nil
}

func (s *serv) UpdateDish(ctx context.Context, dishID int64, name string, price int64, description string) error {
	dish, _ := s.repo.GetDishByID(ctx, dishID)
	if dish == nil {
		return ErrDishNotFound
	}

	return s.repo.UpdateDish(ctx, dishID, name, price, description)
}

// get all dishes from the database
func (s *serv) GetDishes(ctx context.Context) (*[]models.Dish, error) {
	dishes, err := s.repo.GetAllDishes(ctx)
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
