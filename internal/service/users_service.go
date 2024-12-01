package service

import (
	"context"
	"errors"

	models "github.com/agus-germi/TDL_Dinamita/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	// User messages errors
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrRemovingUser       = errors.New("something went wrong trying to remove a user")
	ErrInvalidCredentials = errors.New("invalid credentials")
	// User role messages errors
	ErrUserRoleAlreadyAdded = errors.New("role was already added for this user")
	ErrRemovingUserRole     = errors.New("something went wrong trying to remove a user role")
	ErrUserRoleNotFound     = errors.New("this user has any role")
)

const adminRoleID = 1

func (s *serv) RegisterUser(ctx context.Context, name, password, email string) error {
	usr, _ := s.repo.GetUserByEmail(ctx, email)
	if usr != nil {
		return ErrUserAlreadyExists
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	return s.repo.SaveUser(ctx, name, hashedPassword, email, 2) // Every user is created with a role 2 (user)
}

func (s *serv) LoginUser(ctx context.Context, email, password string) (*models.User, error) {
	usr, err := s.repo.GetUserByEmail(ctx, email)
	if usr == nil {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	if usr.RoleID != adminRoleID {
		err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password))
		if err != nil {
			return nil, ErrInvalidCredentials
		}
	}

	return &models.User{
		ID:    usr.ID,
		Name:  usr.Name,
		Email: usr.Email}, nil
}

func (s *serv) RemoveUser(ctx context.Context, userID int64) error {
	usr, _ := s.repo.GetUserByID(ctx, userID)
	if usr == nil {
		return ErrUserNotFound
	}

	return s.repo.RemoveUser(ctx, userID)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *serv) UpdateUserRole(ctx context.Context, userID, roleID int64, email string) error {

	hasPermission, err := s.repo.HasPermission(ctx, email)
	if err != nil {
		return err
	}
	if !hasPermission {
		return ErrInvalidPermission
	}

	usr, _ := s.repo.GetUserByID(ctx, userID)
	if usr == nil {
		return ErrUserNotFound
	}

	return s.repo.SaveUpdateUserRole(ctx, userID, roleID)
}
