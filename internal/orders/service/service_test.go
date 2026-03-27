package service

import (
	"context"
	"errors"
	"testing"
	"time"

	commonapi "go-gin-ecommerce/internal/common/api"
	"go-gin-ecommerce/internal/orders/dto"
	"go-gin-ecommerce/internal/orders/model"
	"go-gin-ecommerce/internal/orders/repository"
)

func TestCancelTransitionsPendingOrder(t *testing.T) {
	repo := &stubRepository{
		getByIDFn: func(ctx context.Context, id uint) (model.Order, error) {
			return model.Order{
				ID:        id,
				Status:    "pending",
				UpdatedAt: time.Date(2026, time.January, 10, 9, 0, 0, 0, time.UTC),
			}, nil
		},
		updateFn: func(ctx context.Context, order *model.Order) error {
			if order.Status != "cancelled" {
				t.Fatalf("expected cancelled status, got %q", order.Status)
			}
			if !order.UpdatedAt.After(time.Date(2026, time.January, 10, 9, 0, 0, 0, time.UTC)) {
				t.Fatalf("expected updated_at to advance, got %s", order.UpdatedAt)
			}

			return nil
		},
	}
	repo.getByIDSequence = []model.Order{
		{ID: 42, Status: "pending", UpdatedAt: time.Date(2026, time.January, 10, 9, 0, 0, 0, time.UTC)},
		{ID: 42, Status: "cancelled", UpdatedAt: time.Date(2026, time.January, 10, 10, 0, 0, 0, time.UTC)},
	}

	service := New(repo)
	order, err := service.Cancel(context.Background(), 42)
	if err != nil {
		t.Fatalf("expected cancel to succeed, got %v", err)
	}

	if order.Status != "cancelled" {
		t.Fatalf("expected cancelled order, got %#v", order)
	}
	if repo.updateCalls != 1 {
		t.Fatalf("expected one update call, got %d", repo.updateCalls)
	}
}

func TestCancelReturnsNotFoundForMissingOrder(t *testing.T) {
	service := New(&stubRepository{
		getByIDFn: func(ctx context.Context, id uint) (model.Order, error) {
			return model.Order{}, repository.ErrNotFound
		},
	})

	_, err := service.Cancel(context.Background(), 999)
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

func TestCancelRejectsIneligibleStatus(t *testing.T) {
	repo := &stubRepository{
		getByIDFn: func(ctx context.Context, id uint) (model.Order, error) {
			return model.Order{ID: id, Status: "fulfilled"}, nil
		},
	}

	_, err := New(repo).Cancel(context.Background(), 7)
	if err == nil {
		t.Fatal("expected ineligible status to fail")
	}

	var apiErr *commonapi.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected api error, got %T", err)
	}
	if apiErr.Status != 409 || apiErr.Code != "CONFLICT" {
		t.Fatalf("expected conflict error, got %#v", apiErr)
	}
	if repo.updateCalls != 0 {
		t.Fatalf("expected no update call, got %d", repo.updateCalls)
	}
}

type stubRepository struct {
	getByIDFn       func(context.Context, uint) (model.Order, error)
	updateFn        func(context.Context, *model.Order) error
	getByIDSequence []model.Order
	getByIDCalls    int
	updateCalls     int
}

func (s *stubRepository) List(ctx context.Context, params dto.ListOrdersParams) ([]model.Order, int64, error) {
	return nil, 0, nil
}

func (s *stubRepository) GetByID(ctx context.Context, id uint) (model.Order, error) {
	s.getByIDCalls++
	if len(s.getByIDSequence) >= s.getByIDCalls {
		return s.getByIDSequence[s.getByIDCalls-1], nil
	}
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id)
	}

	return model.Order{}, nil
}

func (s *stubRepository) Update(ctx context.Context, order *model.Order) error {
	s.updateCalls++
	if s.updateFn != nil {
		return s.updateFn(ctx, order)
	}

	return nil
}
