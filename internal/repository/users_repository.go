package repository

import (
	"context"
	"errors"

	"github.com/agus-germi/TDL_Dinamita/internal/entity"
	"github.com/jackc/pgx/v5"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

const (
	invalidRole   = 0
	qryInsertUser = `INSERT INTO users (name, password, email, role_id)
					 VALUES ($1, $2, $3, $4)`

	qryRemoveUser = `DELETE FROM users
					WHERE id=$1`

	qryGetUserByEmail = `SELECT *
						FROM users
						WHERE email=$1`

	qryGetUserByID = `SELECT *
						FROM users
						WHERE id=$1`

	qryUpdateUserRole = `UPDATE users 
						SET role_id=$1
						WHERE id=$2`
)

func (r *repo) SaveUser(ctx context.Context, name, password, email string, roleID int64) error {
	operation := func(tx pgx.Tx) error {
		usr, err := r.GetUserByEmail(ctx, tx, email)
		if usr != nil {
			r.log.Errorf("User (email=%s) already exists.", email)
			return ErrUserAlreadyExists
		}
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return err
		}

		_, err = tx.Exec(ctx, qryInsertUser, name, password, email, roleID)
		if err != nil {
			r.log.Errorf("Failed to execute insert user query: %v", err)
			return err
		}

		r.log.Infof("User (email=%s) saved successfully.", email)
		return nil
	}

	return r.executeInTransaction(ctx, operation)
}

func (r *repo) RemoveUser(ctx context.Context, userID int64) error {
	operation := func(tx pgx.Tx) error {
		usr, err := r.GetUserByID(ctx, tx, userID)
		if usr == nil {
			r.log.Errorf("User (id=%d) not found.", userID)
			return ErrUserNotFound
		}
		if err != nil {
			r.log.Errorf("Failed to get user by ID: %v", err)
			return err
		}

		result, err := tx.Exec(ctx, qryRemoveUser, userID)
		if err != nil {
			r.log.Errorf("Failed to execute delete user query: %v", err)
			return err
		}

		if result.RowsAffected() == 0 {
			r.log.Errorf("No rows were affected by the delete user query.")
			return ErrUserNotFound
		}

		r.log.Infof("User (id=%d) removed successfully.", userID)
		return nil
	}

	return r.executeInTransaction(ctx, operation)
}

func (r *repo) GetUserByEmail(ctx context.Context, tx pgx.Tx, email string) (*entity.User, error) {
	usr := entity.User{}

	var err error
	if tx != nil {
		err = tx.QueryRow(ctx, qryGetUserByEmail, email).Scan(&usr.ID, &usr.Name, &usr.Password, &usr.Email, &usr.RoleID)
	} else {
		err = r.db.QueryRow(ctx, qryGetUserByEmail, email).Scan(&usr.ID, &usr.Name, &usr.Password, &usr.Email, &usr.RoleID)
	}
	if err != nil {
		r.log.Debugf("Failed to execute get user by email query: %v", err)
		return nil, err
	}

	return &usr, nil
}

func (r *repo) GetUserByID(ctx context.Context, tx pgx.Tx, userID int64) (*entity.User, error) {
	usr := entity.User{}

	var err error
	if tx != nil {
		err = tx.QueryRow(ctx, qryGetUserByID, userID).Scan(&usr.ID, &usr.Name, &usr.Password, &usr.Email, &usr.RoleID)
	} else {
		err = r.db.QueryRow(ctx, qryGetUserByID, userID).Scan(&usr.ID, &usr.Name, &usr.Password, &usr.Email, &usr.RoleID)
	}
	if err != nil {
		r.log.Errorf("Failed to execute get user by ID query: %v", err)
		return nil, err
	}

	r.log.Debugf("User (id=%d) retrieved successfully.", userID)
	return &usr, nil
}

func (r *repo) SaveUpdateUserRole(ctx context.Context, userID, roleID int64) error {
	operation := func(tx pgx.Tx) error {
		usr, err := r.GetUserByID(ctx, tx, userID)
		if usr == nil {
			r.log.Errorf("User (id=%d) not found.", userID)
			return ErrUserNotFound
		}
		if err != nil {
			r.log.Errorf("Failed to get user by ID: %v", err)
			return err
		}

		_, err = tx.Exec(ctx, qryUpdateUserRole, roleID, userID)
		if err != nil {
			r.log.Errorf("Failed to execute update user role query: %v", err)
			return err
		}

		r.log.Infof("User (id=%d) updated successfully.", userID)
		return nil
	}

	return r.executeInTransaction(ctx, operation)
}
