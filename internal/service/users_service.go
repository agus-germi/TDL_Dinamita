package service

import (
	"context"
	"errors"
	time "time"

	"github.com/agus-germi/TDL_Dinamita/internal/entity"
	models "github.com/agus-germi/TDL_Dinamita/internal/models"
	"github.com/agus-germi/TDL_Dinamita/internal/repository"
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

func (s *serv) RegisterUser(ctx context.Context, name, password, email string) error {
	start := time.Now()
	var hashedPassword string
	err := s.executeWithTimeout(ctx, s.config.MaxHashOperationDuration, func(ctx context.Context) error {
		var err error
		hashedPassword, err = hashPassword(password)
		return err
	})
	s.log.Debugf("Time to hash password: %v", time.Since(start))
	s.log.Debugf("Hashed password: %v", hashedPassword)

	if err != nil {
		s.log.Errorf("Failed to hash password: %v", err)
		return err
	}

	err = s.executeWithTimeout(ctx, s.config.MaxDBOperationDuration, func(ctx context.Context) error {
		return s.repo.SaveUser(ctx, name, hashedPassword, email, 2) // Every user is created with a role 2 (user)
	})

	if errors.Is(err, repository.ErrUserAlreadyExists) {
		return ErrUserAlreadyExists
	}

	return err
}

func (s *serv) LoginUser(ctx context.Context, email, password string) (*models.User, error) {
	var usr *entity.User
	err := s.executeWithTimeout(ctx, s.config.MaxDBOperationDuration, func(ctx context.Context) error {
		var err error
		usr, err = s.repo.GetUserByEmail(ctx, nil, email)
		return err
	})

	if usr == nil {
		s.log.Debugf("User (email:%s) not found", email)
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	start := time.Now()
	err = s.executeWithTimeout(ctx, s.config.MaxHashOperationDuration, func(ctx context.Context) error {
		return bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(password))
	})
	s.log.Debugf("Time to compare password: %v", time.Since(start))

	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return &models.User{
		ID:     usr.ID,
		Name:   usr.Name,
		Email:  usr.Email,
		RoleID: usr.RoleID}, nil
}

func (s *serv) RemoveUser(ctx context.Context, userID int64) error {
	err := s.executeWithTimeout(ctx, s.config.MaxDBOperationDuration, func(ctx context.Context) error {
		return s.repo.RemoveUser(ctx, userID)
	})

	if errors.Is(err, repository.ErrUserNotFound) {
		return ErrUserNotFound
	}

	return err
}

func (s *serv) UpdateUserRole(ctx context.Context, userID, newRoleID int64) error {
	err := s.executeWithTimeout(ctx, s.config.MaxDBOperationDuration, func(ctx context.Context) error {
		return s.repo.SaveUpdateUserRole(ctx, userID, newRoleID)
	})

	if errors.Is(err, repository.ErrUserNotFound) {
		return ErrUserNotFound
	}

	return err
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
