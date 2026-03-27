package service

import (
	"context"
	"errors"
	"testing"
	"time"

	commonapi "go-gin-ecommerce/internal/common/api"
	"go-gin-ecommerce/internal/promotions/dto"
	"go-gin-ecommerce/internal/promotions/model"
	"go-gin-ecommerce/internal/promotions/repository"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestCreateRejectsInvalidDateWindow(t *testing.T) {
	repo := &stubRepository{}
	startsAt := time.Date(2026, time.April, 10, 0, 0, 0, 0, time.UTC)
	endsAt := time.Date(2026, time.April, 9, 0, 0, 0, 0, time.UTC)

	_, err := New(repo).Create(context.Background(), dto.CreatePromotionRequest{
		Name:          "Spring Sale",
		Code:          "SPRING-10",
		DiscountType:  "percentage",
		DiscountValue: 10,
		StartsAt:      &startsAt,
		EndsAt:        &endsAt,
		Status:        "active",
	})
	if err == nil {
		t.Fatal("expected invalid date window to fail")
	}

	var apiErr *commonapi.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected api error, got %T", err)
	}
	if apiErr.Status != 400 || apiErr.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected validation error, got %#v", apiErr)
	}
	if repo.createCalls != 0 {
		t.Fatalf("expected repository create not to be called, got %d", repo.createCalls)
	}
}

func TestUpdateRejectsMergedInvalidDateWindow(t *testing.T) {
	repo := &stubRepository{
		getByIDFn: func(context.Context, uint) (model.Promotion, error) {
			endsAt := time.Date(2026, time.April, 10, 0, 0, 0, 0, time.UTC)
			return model.Promotion{
				ID:            9,
				Name:          "Spring Sale",
				Code:          "SPRING-10",
				DiscountType:  "percentage",
				DiscountValue: 10,
				EndsAt:        &endsAt,
				Status:        "active",
			}, nil
		},
	}
	startsAt := time.Date(2026, time.April, 11, 0, 0, 0, 0, time.UTC)

	_, err := New(repo).Update(context.Background(), 9, dto.UpdatePromotionRequest{
		StartsAt: dto.OptionalTime{Set: true, Value: startsAt},
	})
	if err == nil {
		t.Fatal("expected invalid merged date window to fail")
	}

	var apiErr *commonapi.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected api error, got %T", err)
	}
	if apiErr.Status != 400 || apiErr.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected validation error, got %#v", apiErr)
	}
	if repo.updateCalls != 0 {
		t.Fatalf("expected repository update not to be called, got %d", repo.updateCalls)
	}
}

func TestCreateMapsDuplicateCodeToConflict(t *testing.T) {
	repo := &stubRepository{
		createFn: func(context.Context, *model.Promotion) error {
			return &pgconn.PgError{Code: "23505"}
		},
	}

	_, err := New(repo).Create(context.Background(), dto.CreatePromotionRequest{
		Name:          "Spring Sale",
		Code:          "SPRING-10",
		DiscountType:  "percentage",
		DiscountValue: 10,
		Status:        "active",
	})
	if err == nil {
		t.Fatal("expected duplicate code to fail")
	}

	var apiErr *commonapi.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected api error, got %T", err)
	}
	if apiErr.Status != 409 || apiErr.Code != "CONFLICT" {
		t.Fatalf("expected conflict error, got %#v", apiErr)
	}
	if len(apiErr.Details) != 1 || apiErr.Details[0].Field != "code" {
		t.Fatalf("expected code conflict detail, got %#v", apiErr.Details)
	}
}

type stubRepository struct {
	listFn    func(context.Context) ([]model.Promotion, error)
	getByIDFn func(context.Context, uint) (model.Promotion, error)
	createFn  func(context.Context, *model.Promotion) error
	updateFn  func(context.Context, *model.Promotion) error
	deleteFn  func(context.Context, uint) error

	createCalls int
	updateCalls int
}

func (s *stubRepository) List(ctx context.Context) ([]model.Promotion, error) {
	if s.listFn != nil {
		return s.listFn(ctx)
	}

	return nil, nil
}

func (s *stubRepository) GetByID(ctx context.Context, id uint) (model.Promotion, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, id)
	}

	return model.Promotion{}, repository.ErrNotFound
}

func (s *stubRepository) Create(ctx context.Context, promotion *model.Promotion) error {
	s.createCalls++
	if s.createFn != nil {
		return s.createFn(ctx, promotion)
	}

	return nil
}

func (s *stubRepository) Update(ctx context.Context, promotion *model.Promotion) error {
	s.updateCalls++
	if s.updateFn != nil {
		return s.updateFn(ctx, promotion)
	}

	return nil
}

func (s *stubRepository) Delete(ctx context.Context, id uint) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, id)
	}

	return nil
}
