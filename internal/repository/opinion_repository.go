package repository

import (
	"context"

	"github.com/agus-germi/TDL_Dinamita/internal/entity"
	"github.com/jackc/pgx/v5"
)

const (
	qryInsertOpinion = `INSERT INTO opinions (user_id, opinion, rating) VALUES ($1, $2, $3)`
	qryGetOpinions   = `SELECT o.id, o.user_id, u.name, o.opinion, o.rating
                            FROM opinions o, users u 
                            WHERE o.user_id = u.id
                            ORDER BY o.id desc`
)

func (r *repo) SaveOpinion(ctx context.Context, userID int64, opinion string, rating int) error {
	r.log.Debugf("Received user_id: %d, opinion: %s, rating: %d", userID, opinion, rating)

	operation := func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, qryInsertOpinion, userID, opinion, rating)
		if err != nil {
			r.log.Errorf("Failed to insert opinion: %v", err)
			return err
		}

		r.log.Infof("Opinion from user %d saved successfully", userID)
		return nil
	}

	return r.executeInTransaction(ctx, operation)
}

func (r *repo) GetAllOpinions(ctx context.Context) (*[]entity.Opinion, error) {
	rows, err := r.db.Query(ctx, qryGetOpinions)
	if err != nil {
		r.log.Errorf("Failed to execute select opinions query: %v", err)
		return nil, err
	}
	defer rows.Close()

	opinions := []entity.Opinion{}
	for rows.Next() {
		var opinion entity.Opinion
		if err := rows.Scan(&opinion.ID, &opinion.UserID, &opinion.Name, &opinion.Opinion, &opinion.Rating); err != nil {
			r.log.Errorf("Failed to scan row: %v", err)
			return nil, err
		}

		opinions = append(opinions, opinion)
	}

	if rows.Err() != nil {
		r.log.Errorf("Error occurred during row iteration: %v", rows.Err())
		return nil, rows.Err()
	}

	r.log.Debug("Opinions retrieved successfully.")
	return &opinions, nil
}
