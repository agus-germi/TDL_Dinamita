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
)

func (r *Repo) SaveUser(ctx context.Context, name, password, email string) error {
	_, err := r.db.ExecContext(ctx, qryInsertUser, name, password, email)
	return err
}

func (r *Repo) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	usr := &entity.User{}

	err := r.db.GetContext(ctx, usr, qryGetUserByEmail, email)
	if err != nil {
		return nil, err
	}

	return usr, nil
}
