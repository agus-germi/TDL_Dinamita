package service

import (
	context "context"
	"errors"
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
