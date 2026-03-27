package dto

import (
	"time"

	"go-gin-ecommerce/internal/products/model"
)

type CreateProductRequest struct {
	Name       string  `json:"name" binding:"required,min=1"`
	SKU        string  `json:"sku" binding:"required,min=1"`
	Price      float64 `json:"price" binding:"required,gt=0"`
	Status     string  `json:"status" binding:"required,min=1"`
	CategoryID *uint   `json:"categoryId"`
}

type UpdateProductRequest struct {
	Name       *string  `json:"name" binding:"omitempty,min=1"`
	SKU        *string  `json:"sku" binding:"omitempty,min=1"`
	Price      *float64 `json:"price" binding:"omitempty,gt=0"`
	Status     *string  `json:"status" binding:"omitempty,min=1"`
	CategoryID *uint    `json:"categoryId"`
}

type ProductResponse struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	SKU        string    `json:"sku"`
	Price      float64   `json:"price"`
	Status     string    `json:"status"`
	CategoryID *uint     `json:"categoryId,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func NewProductResponse(product model.Product) ProductResponse {
	return ProductResponse{
		ID:         product.ID,
		Name:       product.Name,
		SKU:        product.SKU,
		Price:      product.Price,
		Status:     product.Status,
		CategoryID: product.CategoryID,
		CreatedAt:  product.CreatedAt,
		UpdatedAt:  product.UpdatedAt,
	}
}

func NewProductResponses(products []model.Product) []ProductResponse {
	responses := make([]ProductResponse, 0, len(products))
	for _, product := range products {
		responses = append(responses, NewProductResponse(product))
	}

	return responses
}
