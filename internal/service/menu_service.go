package service

import (
	context "context"
	"errors"
)

var (
	ErrDishAlreadyExists = errors.New("dish already exists")
)

func (s *serv) AddDishToMenu(ctx context.Context, name string, price int64, description string) error {
	dish, _ := s.repo.GetDishByName(ctx, name)
	if dish != nil {
		return ErrDishAlreadyExists
	}

	return s.repo.SaveDish(ctx, name, price, description)
}
