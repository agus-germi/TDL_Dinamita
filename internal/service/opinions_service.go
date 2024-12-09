package service

import (
	"context"
	"errors"

	"github.com/agus-germi/TDL_Dinamita/internal/entity"
	"github.com/agus-germi/TDL_Dinamita/internal/models"
)

var (
	ErrOpinionNotFound = errors.New("opinion not found")
	ErrAddingOpinion   = errors.New("something went wrong trying to add an opinion")
)

func (s *serv) CreateOpinion(ctx context.Context, userID int64, opinionText string, opinionRating int) error {
	s.log.Debugf("Creating opinion for user_id: %d with opinion: %s and rating: %d", userID, opinionText, opinionRating)

	err := s.executeWithTimeout(ctx, s.config.MaxDBOperationDuration, func(ctx context.Context) error {
		return s.repo.SaveOpinion(ctx, userID, opinionText, opinionRating)
	})

	if err != nil {
		return ErrAddingOpinion
	}

	return nil
}

func (s *serv) GetOpinions(ctx context.Context) (*[]models.Opinion, error) {
	var opinions *[]entity.Opinion
	err := s.executeWithTimeout(ctx, s.config.MaxDBOperationDuration, func(ctx context.Context) error {
		var err error
		opinions, err = s.repo.GetAllOpinions(ctx)
		return err
	})

	if err != nil {
		return nil, err
	}

	s.log.Debugf("Retrieved opinions: %+v", opinions)
	modelOpinions := make([]models.Opinion, len(*opinions))
	for i, opinion := range *opinions {
		modelOpinions[i] = models.Opinion{
			ID:      opinion.ID,
			Name:    opinion.Name,
			UserID:  opinion.UserID,
			Opinion: opinion.Opinion,
			Rating:  opinion.Rating,
		}
	}

	return &modelOpinions, nil
}
