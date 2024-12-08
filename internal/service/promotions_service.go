package service

import (
    "context"
    "errors"
    "log"
    
    // "github.com/agus-germi/TDL_Dinamita/internal/models"
)

var (
    ErrAddingPromotion  = errors.New("something went wrong trying to add a promotion")
    ErrDeletingPromotion = errors.New("something went wrong trying to delete a promotion")
    ErrPromotionNotFound = errors.New("promotion not found")


)
    
func (s *serv) CreatePromotion(ctx context.Context, description string, startDate string, dueDate string, discount int) error {
    log.Printf("Creating promotion with description: %s, start_date: %s, due_date: %s, discount: %d", 
        description, startDate, dueDate, discount)
    
    err := s.repo.SavePromotion(ctx, description, startDate, dueDate, discount)
    if err != nil {
        log.Printf("Error while adding promotion: %v", err)
        return ErrAddingPromotion
    }

    return nil
}

func (s *serv) DeletePromotion(ctx context.Context, promotionID int64) error {
    log.Printf("Deleting promotion with ID: %d", promotionID)

    promotion, err := s.repo.GetPromotionbyID(ctx, promotionID)
    if err != nil {
        log.Printf("Error while fetching promotion: %v", err)
        return ErrPromotionNotFound
    }

    if promotion == nil {
        return ErrPromotionNotFound
    }

    err = s.repo.DeletePromotion(ctx, promotionID)
    if err != nil {
        log.Printf("Error while deleting promotion: %v", err)
        return ErrDeletingPromotion
    }

    return nil
}

