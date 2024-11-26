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
	qryInsertUser = `INSERT INTO users (name, password, email)
					 VALUES ($1, $2, $3)`

	qryRemoveUser = `DELETE FROM users
					WHERE id=$1`

	qryGetUserByEmail = `SELECT id, name, password, email
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
	operation := func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, qryInsertUser, name, password, email)
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

	err := r.db.QueryRow(ctx, qryGetUserByEmail, email).Scan(&usr.ID, &usr.Name, &usr.Password, &usr.Email)
	if err != nil {
		return nil, err
	}

	return &usr, nil
}

func (r *repo) GetUserByID(ctx context.Context, userID int64) (*entity.User, error) {
	usr := entity.User{}

	err := r.db.QueryRow(ctx, qryGetUserByID, userID).Scan(&usr.ID, &usr.Name, &usr.Password, &usr.Email)
	if err != nil {
		return nil, err
	}

	return &usr, nil
}

// TODO: Implement the following 3 methods after we delete the user_roles table.
// func (r *repo) SaveUserRole(ctx context.Context, userID, roleID int64) error {
// 	usr, _ := r.GetUserByID(ctx, userID)
// 	if usr == nil {
// 		return ErrUserNotFound
// 	}

// 	_, err := r.db.ExecContext(ctx, qryInsertUserRole, userID, roleID)

// 	return err
// }

// func (r *repo) RemoveUserRole(ctx context.Context, userID int64) error {
// 	_, err := r.db.ExecContext(ctx, qryRemoveUserRole, userID)
// 	return err
// }

// func (r *repo) GetUserRole(ctx context.Context, userID int64) (*entity.UserRole, error) {
// 	usr_role := &entity.UserRole{}

// 	err := r.db.GetContext(ctx, usr_role, qryGetUserRoleByUserID, userID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return usr_role, nil
// }
