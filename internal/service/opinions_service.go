package service

import (
    "context"
    "errors"
    "log"
    
    "github.com/agus-germi/TDL_Dinamita/internal/models"
)

var (
    ErrOpinionNotFound = errors.New("opinion not found")
    ErrAddingOpinion   = errors.New("something went wrong trying to add an opinion")
)

func (s *serv) CreateOpinion(ctx context.Context, opinion models.Opinion) error {
    log.Printf("Creating opinion for user_id: %d with opinion: %s", opinion.UserID, opinion.Opinion)
    err := s.repo.SaveOpinion(ctx, int64(opinion.ID), opinion.Opinion)

	if err != nil {
        return ErrAddingOpinion
    }
    return nil
}


func (s *serv) GetOpinions(ctx context.Context) (*[]models.Opinion, error) {
    opinions, err := s.repo.GetAllOpinions(ctx)
    if err != nil {
        return nil, err
    }

    modelOpinions := make([]models.Opinion, len(*opinions))
    for i, opinion := range *opinions {
        modelOpinions[i] = models.Opinion{
            ID:     opinion.ID,
            UserID: opinion.UserID,
            Opinion: opinion.Opinion, 
        }
    }

    return &modelOpinions, nil
}

