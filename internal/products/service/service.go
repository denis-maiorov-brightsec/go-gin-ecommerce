package service

import (
	"context"
	"errors"

	commonapi "go-gin-ecommerce/internal/common/api"
	commonpagination "go-gin-ecommerce/internal/common/pagination"
	"go-gin-ecommerce/internal/products/dto"
	"go-gin-ecommerce/internal/products/model"
	"go-gin-ecommerce/internal/products/repository"
)

type Service interface {
	List(ctx context.Context, params commonpagination.Params) ([]model.Product, int64, error)
	GetByID(ctx context.Context, id uint) (model.Product, error)
	Create(ctx context.Context, request dto.CreateProductRequest) (model.Product, error)
	Update(ctx context.Context, id uint, request dto.UpdateProductRequest) (model.Product, error)
	Delete(ctx context.Context, id uint) error
}

type ProductService struct {
	repository repository.Repository
}

func New(repo repository.Repository) *ProductService {
	return &ProductService{repository: repo}
}

func (s *ProductService) List(ctx context.Context, params commonpagination.Params) ([]model.Product, int64, error) {
	return s.repository.List(ctx, params)
}

func (s *ProductService) GetByID(ctx context.Context, id uint) (model.Product, error) {
	product, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return model.Product{}, mapRepositoryError(err)
	}

	return product, nil
}

func (s *ProductService) Create(ctx context.Context, request dto.CreateProductRequest) (model.Product, error) {
	product := model.Product{
		Name:       request.Name,
		SKU:        request.SKU,
		Price:      request.Price,
		Status:     request.Status,
		CategoryID: request.CategoryID,
	}

	if err := s.repository.Create(ctx, &product); err != nil {
		return model.Product{}, err
	}

	return product, nil
}

func (s *ProductService) Update(ctx context.Context, id uint, request dto.UpdateProductRequest) (model.Product, error) {
	product, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return model.Product{}, mapRepositoryError(err)
	}

	if request.Name != nil {
		product.Name = *request.Name
	}
	if request.SKU != nil {
		product.SKU = *request.SKU
	}
	if request.Price != nil {
		product.Price = *request.Price
	}
	if request.Status != nil {
		product.Status = *request.Status
	}
	if request.CategoryID.Set {
		if request.CategoryID.Null {
			product.CategoryID = nil
		} else {
			categoryID := request.CategoryID.Value
			product.CategoryID = &categoryID
		}
	}

	if err := s.repository.Update(ctx, &product); err != nil {
		return model.Product{}, mapRepositoryError(err)
	}

	updatedProduct, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return model.Product{}, mapRepositoryError(err)
	}

	return updatedProduct, nil
}

func (s *ProductService) Delete(ctx context.Context, id uint) error {
	return mapRepositoryError(s.repository.Delete(ctx, id))
}

func mapRepositoryError(err error) error {
	if errors.Is(err, repository.ErrNotFound) {
		return commonapi.NewNotFoundError()
	}

	return err
}
