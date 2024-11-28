package repository

import (
	"context"
	"errors"
	"log"

	"github.com/agus-germi/TDL_Dinamita/internal/entity"
)

var ErrUserNotFound = errors.New("user not found")

const (
	qryInsertUser = `INSERT INTO users (name, password, email,role)
					 VALUES ($1, $2, $3, 2)`

	qryRemoveUser = `DELETE FROM users
					WHERE id=$1`

	qryGetUserByEmail = `SELECT id, name, password, email, role
						FROM users
						WHERE email=$1`

	qryGetUserByID = `SELECT id, name, password, email
						FROM users
						WHERE id=$1`

	qryInsertUserRole = `INSERT INTO user_roles (user_id, role_id)
						VALUES ($1, $2)`

	qryRemoveUserRole = `DELETE FROM user_roles
						WHERE user_id=$1`

	qryGetUserRoleByUserID = `SELECT user_id, role_id
							 FROM user_roles
							 WHERE user_id=$1`
)

func (r *repo) SaveUser(ctx context.Context, name, password, email string) error {
	_, err := r.db.ExecContext(ctx, qryInsertUser, name, password, email)
	return err
}

func (r *repo) RemoveUser(ctx context.Context, userID int64) error {
	result, err := r.db.ExecContext(ctx, qryRemoveUser, userID)
	if err != nil {
		return err // Return the error from the query
	}

	// Check the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error in checking the number of rows affected")
		return err // Error when determining rows affected
	}

	if rowsAffected == 0 {
		log.Println("Rows affected = 0")
		return ErrUserNotFound
	}

	return nil
}

func (r *repo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	usr := &entity.User{}

	err := r.db.GetContext(ctx, usr, qryGetUserByEmail, email)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (r *repo) SaveUserRole(ctx context.Context, userID, roleID int64) error {
	usr, _ := r.GetUserByID(ctx, userID)
	if usr == nil {
		return ErrUserNotFound
	}

	_, err := r.db.ExecContext(ctx, qryInsertUserRole, userID, roleID)

	return err
}

func (r *repo) RemoveUserRole(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(ctx, qryRemoveUserRole, userID)
	return err
}

func (r *repo) GetUserRole(ctx context.Context, userID int64) (*entity.UserRole, error) {
	usr_role := &entity.UserRole{}

	err := r.db.GetContext(ctx, usr_role, qryGetUserRoleByUserID, userID)
	if err != nil {
		return nil, err
	}

	return usr_role, nil
}

func (r *repo) GetUserByID(ctx context.Context, userID int64) (*entity.User, error) {
	usr := &entity.User{}

	err := r.db.GetContext(ctx, usr, qryGetUserByID, userID)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (r *repo) HasPermission(ctx context.Context, email string) bool {
	usr, _ := r.GetUserByEmail(ctx, email)
	println("ROLE: ", usr.Role)
	roleID := usr.Role
	return roleID == 1
}
