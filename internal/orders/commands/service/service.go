package service

import (
	"context"
	"errors"
	"time"

	commonapi "go-gin-ecommerce/internal/common/api"
	"go-gin-ecommerce/internal/orders/commands/repository"
	"go-gin-ecommerce/internal/orders/model"
)

type Service interface {
	Cancel(ctx context.Context, id uint) (model.Order, error)
}

type OrderService struct {
	repository repository.Repository
}

func New(repo repository.Repository) *OrderService {
	return &OrderService{repository: repo}
}

func (s *OrderService) Cancel(ctx context.Context, id uint) (model.Order, error) {
	order, err := s.repository.Cancel(ctx, id, time.Now().UTC())
	if err != nil {
		return model.Order{}, mapRepositoryError(err)
	}

	return order, nil
}

func mapRepositoryError(err error) error {
	if errors.Is(err, repository.ErrNotFound) {
		return commonapi.NewNotFoundError()
	}
	if errors.Is(err, repository.ErrInvalidTransition) {
		return commonapi.NewConflictError("Order cannot be cancelled", []commonapi.ErrorDetail{{
			Field:       "status",
			Constraints: []string{"only pending orders can be cancelled"},
		}})
	}

	return err
}
