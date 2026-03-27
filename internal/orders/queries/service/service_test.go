package service

import (
	"context"
	"errors"
	"testing"

	commonapi "go-gin-ecommerce/internal/common/api"
	"go-gin-ecommerce/internal/orders/dto"
	"go-gin-ecommerce/internal/orders/model"
	"go-gin-ecommerce/internal/orders/queries/repository"
)

func TestGetByIDReturnsNotFoundForMissingOrder(t *testing.T) {
	service := New(&stubRepository{
		getByIDFn: func(ctx context.Context, id uint) (model.Order, error) {
			return model.Order{}, repository.ErrNotFound
		},
	})

	_, err := service.GetByID(context.Background(), 999)
	if err == nil {
		t.Fatal("expected missing order to fail")
	}

	var apiErr *commonapi.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected api error, got %T", err)
	}
	if apiErr.Status != 404 || apiErr.Code != "NOT_FOUND" {
		t.Fatalf("expected not found error, got %#v", apiErr)
	}
}

type stubRepository struct {
	getByIDFn func(context.Context, uint) (model.Order, error)
}

func (s *stubRepository) List(ctx context.Context, params dto.ListOrdersParams) ([]model.Order, int64, error) {
	return nil, 0, nil
}

func (s *stubRepository) GetByID(ctx context.Context, id uint) (model.Order, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id)
	}

	return model.Order{}, nil
}
