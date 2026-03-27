package service

import (
	"context"
	"errors"

	"go-gin-ecommerce/internal/categories/dto"
	"go-gin-ecommerce/internal/categories/model"
	"go-gin-ecommerce/internal/categories/repository"
	commonapi "go-gin-ecommerce/internal/common/api"

	"github.com/jackc/pgx/v5/pgconn"
)

type Service interface {
	List(ctx context.Context) ([]model.Category, error)
	GetByID(ctx context.Context, id uint) (model.Category, error)
	Create(ctx context.Context, request dto.CreateCategoryRequest) (model.Category, error)
	Update(ctx context.Context, id uint, request dto.UpdateCategoryRequest) (model.Category, error)
	Delete(ctx context.Context, id uint) error
}

type CategoryService struct {
	repository repository.Repository
}

func New(repo repository.Repository) *CategoryService {
	return &CategoryService{repository: repo}
}

func (s *CategoryService) List(ctx context.Context) ([]model.Category, error) {
	return s.repository.List(ctx)
}

func (s *CategoryService) GetByID(ctx context.Context, id uint) (model.Category, error) {
	category, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return model.Category{}, mapRepositoryError(err)
	}

	return category, nil
}

func (s *CategoryService) Create(ctx context.Context, request dto.CreateCategoryRequest) (model.Category, error) {
	category := model.Category{
		Name:        request.Name,
		Slug:        request.Slug,
		Description: request.Description,
	}

	if err := s.repository.Create(ctx, &category); err != nil {
		return model.Category{}, mapRepositoryError(err)
	}

	return category, nil
}

func (s *CategoryService) Update(ctx context.Context, id uint, request dto.UpdateCategoryRequest) (model.Category, error) {
	category, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return model.Category{}, mapRepositoryError(err)
	}

	if request.Name != nil {
		category.Name = *request.Name
	}
	if request.Slug != nil {
		category.Slug = *request.Slug
	}
	if request.Description.Set {
		if request.Description.Null {
			category.Description = nil
		} else {
			description := request.Description.Value
			category.Description = &description
		}
	}

	if err := s.repository.Update(ctx, &category); err != nil {
		return model.Category{}, mapRepositoryError(err)
	}

	updatedCategory, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return model.Category{}, mapRepositoryError(err)
	}

	return updatedCategory, nil
}

func (s *CategoryService) Delete(ctx context.Context, id uint) error {
	return mapRepositoryError(s.repository.Delete(ctx, id))
}

func mapRepositoryError(err error) error {
	if errors.Is(err, repository.ErrNotFound) {
		return commonapi.NewNotFoundError()
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return commonapi.NewConflictError("Resource conflict", []commonapi.ErrorDetail{{
			Field:       "slug",
			Constraints: []string{"slug must be unique"},
		}})
	}

	return err
}
