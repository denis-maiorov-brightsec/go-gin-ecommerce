package service

import (
	"context"
	"errors"

	commonapi "go-gin-ecommerce/internal/common/api"
	"go-gin-ecommerce/internal/orders/dto"
	"go-gin-ecommerce/internal/orders/model"
	"go-gin-ecommerce/internal/orders/repository"
)

type Service interface {
	List(ctx context.Context, params dto.ListOrdersParams) ([]model.Order, int64, error)
	GetByID(ctx context.Context, id uint) (model.Order, error)
}

type OrderService struct {
	repository repository.Repository
}

func New(repo repository.Repository) *OrderService {
	return &OrderService{repository: repo}
}

func (s *OrderService) List(ctx context.Context, params dto.ListOrdersParams) ([]model.Order, int64, error) {
	return s.repository.List(ctx, params)
}

func (s *OrderService) GetByID(ctx context.Context, id uint) (model.Order, error) {
	order, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return model.Order{}, mapRepositoryError(err)
	}

	return order, nil
}

func mapRepositoryError(err error) error {
	if errors.Is(err, repository.ErrNotFound) {
		return commonapi.NewNotFoundError()
	}

	return err
}
