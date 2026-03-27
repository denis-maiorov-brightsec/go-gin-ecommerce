package service

import (
	"context"
	"errors"
	"time"

	commonapi "go-gin-ecommerce/internal/common/api"
	"go-gin-ecommerce/internal/promotions/dto"
	"go-gin-ecommerce/internal/promotions/model"
	"go-gin-ecommerce/internal/promotions/repository"

	"github.com/jackc/pgx/v5/pgconn"
)

type Service interface {
	List(ctx context.Context) ([]model.Promotion, error)
	GetByID(ctx context.Context, id uint) (model.Promotion, error)
	Create(ctx context.Context, request dto.CreatePromotionRequest) (model.Promotion, error)
	Update(ctx context.Context, id uint, request dto.UpdatePromotionRequest) (model.Promotion, error)
	Delete(ctx context.Context, id uint) error
}

type PromotionService struct {
	repository repository.Repository
}

func New(repo repository.Repository) *PromotionService {
	return &PromotionService{repository: repo}
}

func (s *PromotionService) List(ctx context.Context) ([]model.Promotion, error) {
	return s.repository.List(ctx)
}

func (s *PromotionService) GetByID(ctx context.Context, id uint) (model.Promotion, error) {
	promotion, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return model.Promotion{}, mapRepositoryError(err)
	}

	return promotion, nil
}

func (s *PromotionService) Create(ctx context.Context, request dto.CreatePromotionRequest) (model.Promotion, error) {
	if err := validateDateWindow(request.StartsAt, request.EndsAt); err != nil {
		return model.Promotion{}, err
	}

	promotion := model.Promotion{
		Name:          request.Name,
		Code:          request.Code,
		DiscountType:  request.DiscountType,
		DiscountValue: request.DiscountValue,
		StartsAt:      request.StartsAt,
		EndsAt:        request.EndsAt,
		Status:        request.Status,
	}

	if err := s.repository.Create(ctx, &promotion); err != nil {
		return model.Promotion{}, mapRepositoryError(err)
	}

	return promotion, nil
}

func (s *PromotionService) Update(ctx context.Context, id uint, request dto.UpdatePromotionRequest) (model.Promotion, error) {
	promotion, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return model.Promotion{}, mapRepositoryError(err)
	}

	if request.Name != nil {
		promotion.Name = *request.Name
	}
	if request.Code != nil {
		promotion.Code = *request.Code
	}
	if request.DiscountType != nil {
		promotion.DiscountType = *request.DiscountType
	}
	if request.DiscountValue != nil {
		promotion.DiscountValue = *request.DiscountValue
	}
	if request.StartsAt.Set {
		if request.StartsAt.Null {
			promotion.StartsAt = nil
		} else {
			startsAt := request.StartsAt.Value
			promotion.StartsAt = &startsAt
		}
	}
	if request.EndsAt.Set {
		if request.EndsAt.Null {
			promotion.EndsAt = nil
		} else {
			endsAt := request.EndsAt.Value
			promotion.EndsAt = &endsAt
		}
	}
	if request.Status != nil {
		promotion.Status = *request.Status
	}

	if err := validateDateWindow(promotion.StartsAt, promotion.EndsAt); err != nil {
		return model.Promotion{}, err
	}

	if err := s.repository.Update(ctx, &promotion); err != nil {
		return model.Promotion{}, mapRepositoryError(err)
	}

	updatedPromotion, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return model.Promotion{}, mapRepositoryError(err)
	}

	return updatedPromotion, nil
}

func (s *PromotionService) Delete(ctx context.Context, id uint) error {
	return mapRepositoryError(s.repository.Delete(ctx, id))
}

func validateDateWindow(startsAt *time.Time, endsAt *time.Time) error {
	if startsAt != nil && endsAt != nil && startsAt.After(*endsAt) {
		return commonapi.NewValidationError([]commonapi.ErrorDetail{{
			Field:       "startsAt",
			Constraints: []string{"startsAt must be less than or equal to endsAt"},
		}})
	}

	return nil
}

func mapRepositoryError(err error) error {
	if errors.Is(err, repository.ErrNotFound) {
		return commonapi.NewNotFoundError()
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return commonapi.NewConflictError("Resource conflict", []commonapi.ErrorDetail{{
			Field:       "code",
			Constraints: []string{"code must be unique"},
		}})
	}

	return err
}
