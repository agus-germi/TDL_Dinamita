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

func (s *serv) CreateOpinion(ctx context.Context, userID int64, opinionText string, opinionRating int) error {
    log.Printf("Creating opinion for user_id: %d with opinion: %s and rating: %d", userID, opinionText, opinionRating)
    
    err := s.repo.SaveOpinion(ctx, userID, opinionText, opinionRating)
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

    log.Printf("Retrieved opinions: %+v", opinions)
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

