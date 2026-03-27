package service

import (
	"context"
	"errors"
	"testing"
	"time"

	commonapi "go-gin-ecommerce/internal/common/api"
	"go-gin-ecommerce/internal/orders/commands/repository"
	"go-gin-ecommerce/internal/orders/model"
)

func TestCancelTransitionsPendingOrder(t *testing.T) {
	repo := &stubRepository{
		cancelFn: func(ctx context.Context, id uint, updatedAt time.Time) (model.Order, error) {
			if id != 42 {
				t.Fatalf("expected order id 42, got %d", id)
			}
			if !updatedAt.After(time.Date(2026, time.January, 10, 9, 0, 0, 0, time.UTC)) {
				t.Fatalf("expected updated_at to advance, got %s", updatedAt)
			}

			order := model.Order{
				ID:        id,
				Status:    "cancelled",
				UpdatedAt: updatedAt,
			}
			if order.Status != "cancelled" {
				t.Fatalf("expected cancelled status, got %q", order.Status)
			}
			return order, nil
		},
	}

	service := New(repo)
	order, err := service.Cancel(context.Background(), 42)
	if err != nil {
		t.Fatalf("expected cancel to succeed, got %v", err)
	}

	if order.Status != "cancelled" {
		t.Fatalf("expected cancelled order, got %#v", order)
	}
	if repo.cancelCalls != 1 {
		t.Fatalf("expected one cancel call, got %d", repo.cancelCalls)
	}
}

func TestCancelReturnsNotFoundForMissingOrder(t *testing.T) {
	service := New(&stubRepository{
		cancelFn: func(ctx context.Context, id uint, updatedAt time.Time) (model.Order, error) {
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
		cancelFn: func(ctx context.Context, id uint, updatedAt time.Time) (model.Order, error) {
			return model.Order{}, repository.ErrInvalidTransition
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
	if repo.cancelCalls != 1 {
		t.Fatalf("expected one cancel call, got %d", repo.cancelCalls)
	}
}

type stubRepository struct {
	cancelFn    func(context.Context, uint, time.Time) (model.Order, error)
	cancelCalls int
}

func (s *stubRepository) Cancel(ctx context.Context, id uint, updatedAt time.Time) (model.Order, error) {
	s.cancelCalls++
	if s.cancelFn != nil {
		return s.cancelFn(ctx, id, updatedAt)
	}

	return model.Order{}, nil
}
