package dto

import (
	"time"

	"go-gin-ecommerce/internal/products/model"
)

type ProductResponse struct {
	ID               uint      `json:"id" example:"42"`
	Name             string    `json:"name" example:"Trail Running Shoes"`
	StockKeepingUnit string    `json:"stockKeepingUnit" example:"TRAIL-001"`
	Price            float64   `json:"price" example:"129.99"`
	Status           string    `json:"status" example:"active"`
	CategoryID       *uint     `json:"categoryId,omitempty" example:"3"`
	CreatedAt        time.Time `json:"createdAt" example:"2026-03-01T10:00:00Z"`
	UpdatedAt        time.Time `json:"updatedAt" example:"2026-03-02T11:00:00Z"`
}

type ProductListResponse struct {
	Items      []ProductResponse `json:"items"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	Total      int64             `json:"total"`
	TotalPages int               `json:"totalPages"`
}

func NewProductResponse(product model.Product) ProductResponse {
	return ProductResponse{
		ID:               product.ID,
		Name:             product.Name,
		StockKeepingUnit: product.StockKeepingUnit,
		Price:            product.Price,
		Status:           product.Status,
		CategoryID:       product.CategoryID,
		CreatedAt:        product.CreatedAt,
		UpdatedAt:        product.UpdatedAt,
	}
}

func NewProductResponses(products []model.Product) []ProductResponse {
	responses := make([]ProductResponse, 0, len(products))
	for _, product := range products {
		responses = append(responses, NewProductResponse(product))
	}

	return responses
}
