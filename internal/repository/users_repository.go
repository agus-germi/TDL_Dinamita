package repository

import (
	"context"
	"errors"
	"log"

	"github.com/agus-germi/TDL_Dinamita/internal/entity"
	"github.com/jackc/pgx/v5"
)

var ErrUserNotFound = errors.New("user not found")

const (
	invalidRole   = 0 // 0 is an invalid role
	qryInsertUser = `INSERT INTO users (name, password, email, role_id)
					 VALUES ($1, $2, $3, $4)`

	qryRemoveUser = `DELETE FROM users
					WHERE id=$1`

	qryGetUserByEmail = `SELECT id, name, password, email, role_id
						FROM users
						WHERE email=$1`

	qryGetUserByID = `SELECT id, name, password, email, role_id
						FROM users
						WHERE id=$1`

	qryUpdateUserRole = `UPDATE users 
						SET role_id=$1
						WHERE id=$2`

	qryGetUserRoleByUserID = `SELECT role_id
							 FROM users
							 WHERE id=$1`
)

func (r *repo) SaveUser(ctx context.Context, name, password, email string, roleID int64) error {
	operation := func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, qryInsertUser, name, password, email, roleID)
		if err != nil {
			log.Printf("Failed to execute insert query: %v", err)
			return err
		}

		log.Println("User saved successfully.")
		return nil
	}

	return r.executeInTransaction(ctx, operation)
}

func (r *repo) RemoveUser(ctx context.Context, userID int64) error {
	operation := func(tx pgx.Tx) error {
		result, err := tx.Exec(ctx, qryRemoveUser, userID)
		if err != nil {
			log.Printf("Failed to execute delete query: %v", err)
			return err
		}

		if result.RowsAffected() == 0 {
			log.Println("No rows were affected by the delete query.")
			return ErrUserNotFound
		}

		log.Println("User removed successfully.")
		return nil
	}

	return r.executeInTransaction(ctx, operation)
}

func (r *repo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	usr := entity.User{}

	err := r.db.QueryRow(ctx, qryGetUserByEmail, email).Scan(&usr.ID, &usr.Name, &usr.Password, &usr.Email, &usr.RoleID)
	if err != nil {
		return nil, err
	}

	return &usr, nil
}

func (r *repo) GetUserByID(ctx context.Context, userID int64) (*entity.User, error) {
	usr := entity.User{}

	err := r.db.QueryRow(ctx, qryGetUserByID, userID).Scan(&usr.ID, &usr.Name, &usr.Password, &usr.Email, &usr.RoleID)
	if err != nil {
		return nil, err
	}

	return &usr, nil
}

// TODO: Implement the following 3 methods after we delete the user_roles table.
func (r *repo) SaveUpdateUserRole(ctx context.Context, userID, roleID int64) error {
	operation := func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, qryUpdateUserRole, roleID, userID)
		if err != nil {
			log.Printf("Failed to execute update user role query: %v", err)
			return err
		}

		log.Println("User updated successfully.")
		return nil
	}

	return r.executeInTransaction(ctx, operation)
}

func (r *repo) GetUserRole(ctx context.Context, userID int64) (int64, error) {
	usr := entity.User{}

	err := r.db.QueryRow(ctx, qryGetUserRoleByUserID, userID).Scan(&usr.RoleID)
	if err != nil {
		return invalidRole, err
	}

	return usr.RoleID, nil
}
