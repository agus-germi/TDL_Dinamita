package repository

import (
	"context"

	"github.com/agus-germi/TDL_Dinamita/internal/entity"
)

// Si por algun motivo llegase a fallar las querys, probar poniendo ";" al final de cada query.
const (
	qryInsertUser = `INSERT INTO users (name, password, email)
					  VALUES ($1, $2, $3)`

	qryGetUserByEmail = `SELECT id, name, password, email 
						FROM users
						WHERE email=$1`

	qryInsertUserRole = `INSERT INTO user_roles (user_id, role_id)
						VALUES ($1, $2)`

	qryRemoveUserRole = `DELETE FROM user_roles
						WHERE user_id=$1`

	qryGetUserRoleByUserID = `SELECT userID, roleID
							 FROM user_roles
							 WHERE user_id=$1`
)

func (r *repo) SaveUser(ctx context.Context, name, password, email string) error {
	_, err := r.db.ExecContext(ctx, qryInsertUser, name, password, email)
	return err
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
