package dto

import (
	"bytes"
	"encoding/json"
	"time"

	commonapi "go-gin-ecommerce/internal/common/api"
	"go-gin-ecommerce/internal/products/model"
)

type CreateProductRequest struct {
	Name             string  `json:"name" binding:"required,min=1"`
	StockKeepingUnit *string `json:"stockKeepingUnit" binding:"omitempty,min=1"`
	// Deprecated: sku is accepted as a request alias for backward compatibility.
	SKUAlias   *string `json:"sku" binding:"omitempty,min=1"`
	Price      float64 `json:"price" binding:"required,gt=0"`
	Status     string  `json:"status" binding:"required,min=1"`
	CategoryID *uint   `json:"categoryId"`
}

type UpdateProductRequest struct {
	Name             *string `json:"name" binding:"omitempty,min=1"`
	StockKeepingUnit *string `json:"stockKeepingUnit" binding:"omitempty,min=1"`
	// Deprecated: sku is accepted as a request alias for backward compatibility.
	SKUAlias   *string      `json:"sku" binding:"omitempty,min=1"`
	Price      *float64     `json:"price" binding:"omitempty,gt=0"`
	Status     *string      `json:"status" binding:"omitempty,min=1"`
	CategoryID OptionalUint `json:"categoryId"`
}

type ProductResponse struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`
	StockKeepingUnit string    `json:"stockKeepingUnit"`
	Price            float64   `json:"price"`
	Status           string    `json:"status"`
	CategoryID       *uint     `json:"categoryId,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
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

type OptionalUint struct {
	Set   bool
	Null  bool
	Value uint
}

func (r CreateProductRequest) ResolvedStockKeepingUnit() (string, error) {
	return resolveStockKeepingUnit(r.StockKeepingUnit, r.SKUAlias, true)
}

func (r UpdateProductRequest) ResolvedStockKeepingUnit() (*string, error) {
	value, err := resolveStockKeepingUnit(r.StockKeepingUnit, r.SKUAlias, false)
	if err != nil {
		return nil, err
	}
	if value == "" {
		return nil, nil
	}

	return &value, nil
}

func (o *OptionalUint) UnmarshalJSON(data []byte) error {
	o.Set = true

	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		o.Null = true
		o.Value = 0
		return nil
	}

	o.Null = false
	return unmarshalUint(data, &o.Value)
}

func resolveStockKeepingUnit(stockKeepingUnit *string, skuAlias *string, required bool) (string, error) {
	switch {
	case stockKeepingUnit != nil && skuAlias != nil && *stockKeepingUnit != *skuAlias:
		return "", commonapi.NewValidationError([]commonapi.ErrorDetail{{
			Field: "stockKeepingUnit",
			Constraints: []string{
				"stockKeepingUnit and deprecated sku must match when both are provided",
			},
		}})
	case stockKeepingUnit != nil:
		return *stockKeepingUnit, nil
	case skuAlias != nil:
		return *skuAlias, nil
	case required:
		return "", commonapi.NewValidationError([]commonapi.ErrorDetail{{
			Field:       "stockKeepingUnit",
			Constraints: []string{"stockKeepingUnit must not be empty"},
		}})
	default:
		return "", nil
	}
}

func unmarshalUint(data []byte, value *uint) error {
	type rawUint uint

	var parsed rawUint
	if err := json.Unmarshal(data, &parsed); err != nil {
		return err
	}

	*value = uint(parsed)
	return nil
}
