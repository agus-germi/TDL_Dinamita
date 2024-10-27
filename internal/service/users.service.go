package service

import (
	"context"
	"errors"

	"github.com/agus-germi/TDL_Dinamita/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func (s *Serv) RegisterUser(ctx context.Context, name, password, email string) error {

	usr, _ := s.repo.GetUserByEmail(ctx, email)
	if usr != nil {
		return ErrUserAlreadyExists
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	return s.repo.SaveUser(ctx, name, hashedPassword, email)
}

func (s *Serv) LoginUser(ctx context.Context, email, password string) (*models.User, error) {
	usr, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if usr == nil {
		return nil, ErrUserNotFound // Podemos cambiar este error por ErrInvalidCredentials (y borramos ErrUserNotFound)
	}

	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return &models.User{
		ID:    usr.ID,
		Name:  usr.Name,
		Email: usr.Email}, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
