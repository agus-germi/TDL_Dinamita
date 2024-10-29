package service

import (
	"context"
	"errors"

	"github.com/agus-germi/TDL_Dinamita/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrUserRoleAlreadyExists = errors.New("user role already exists")
	ErrRemovingUserRole      = errors.New("something went wrong trying to remove a user role")
)

func (s *serv) RegisterUser(ctx context.Context, name, password, email string) error {

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

func (s *serv) LoginUser(ctx context.Context, email, password string) (*models.User, error) {
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

// TODO: Pensar en los requirimientos: deseamos poder cambiar el rol de un usuario que ya posee uno previamente?
// Si la respuesta es sí, tenemos que modificar este código. (aunque en nuestro caso no tiene mucho sentido este feature)
// En realidad, pensandolo bien, ya tenemos el feature implementado--> Removemos el rol y despues le agregamos uno nuevo :)
func (s *serv) AddUserRole(ctx context.Context, userID, roleID int64) error {
	usr_role, _ := s.repo.GetUserRole(ctx, userID)
	if usr_role != nil {
		return ErrUserRoleAlreadyExists
	}

	return s.repo.SaveUserRole(ctx, userID, roleID)
}

func (s *serv) RemoveUserRole(ctx context.Context, userID int64) error {
	err := s.repo.RemoveUserRole(ctx, userID)
	if err != nil {
		return ErrRemovingUserRole
	}

	return nil
}
