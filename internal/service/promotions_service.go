package service

import (
	"context"
	"errors"

	"github.com/agus-germi/TDL_Dinamita/internal/entity"
	"github.com/agus-germi/TDL_Dinamita/internal/models"
	"github.com/agus-germi/TDL_Dinamita/internal/repository"
)

var (
	ErrAddingPromotion   = errors.New("something went wrong trying to add a promotion")
	ErrDeletingPromotion = errors.New("something went wrong trying to delete a promotion")
	ErrPromotionNotFound = errors.New("promotion not found")
)

func (s *serv) CreatePromotion(ctx context.Context, description string, startDate string, dueDate string, discount int) error {
	s.log.Debugf("Creating promotion with description: %s, start_date: %s, due_date: %s, discount: %d",
		description, startDate, dueDate, discount)

	err := s.executeWithTimeout(ctx, s.config.MaxDBOperationDuration, func(ctx context.Context) error {
		return s.repo.SavePromotion(ctx, description, startDate, dueDate, discount)
	})

	if err != nil {
		s.log.Debugf("Error while adding promotion: %v", err)
		return ErrAddingPromotion
	}

	return nil
}

func (s *serv) DeletePromotion(ctx context.Context, promotionID int64) error {
	s.log.Debugf("Deleting promotion with ID: %d", promotionID)

	err := s.executeWithTimeout(ctx, s.config.MaxDBOperationDuration, func(ctx context.Context) error {
		return s.repo.DeletePromotion(ctx, promotionID)
	})

	if err != nil {
		if errors.Is(err, repository.ErrPromotionNotFound) {
			return ErrPromotionNotFound
		}
		s.log.Errorf("Error while deleting promotion: %v", err)
		return ErrDeletingPromotion
	}

	return nil
}

func (s *serv) GetPromotions(ctx context.Context) (*[]models.Promotion, error) {
	var promotions *[]entity.Promotion
	err := s.executeWithTimeout(ctx, s.config.MaxDBOperationDuration, func(ctx context.Context) error {
		var err error
		promotions, err = s.repo.GetAllPromotionsAvailable(ctx)
		return err
	})

	if err != nil {
		return nil, err
	}

	s.log.Debugf("Retrieved promotions: %+v", promotions)
	modelPromotions := make([]models.Promotion, len(*promotions))
	for i, promotion := range *promotions {
		modelPromotions[i] = models.Promotion{
			ID:          promotion.ID,
			Description: promotion.Description,
			StartDate:   promotion.StartDate,
			DueDate:     promotion.DueDate,
			Discount:    promotion.Discount,
		}
	}

	return &modelPromotions, nil
}
