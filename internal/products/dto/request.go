package dto

import (
	"bytes"
	"encoding/json"

	commonapi "go-gin-ecommerce/internal/common/api"
)

type CreateProductRequest struct {
	Name             string  `json:"name" binding:"required,min=1" example:"Trail Running Shoes"`
	StockKeepingUnit *string `json:"stockKeepingUnit" binding:"omitempty,min=1" example:"TRAIL-001"`
	// Deprecated: sku is accepted as a request alias for backward compatibility.
	SKUAlias   *string `json:"sku" binding:"omitempty,min=1" example:"TRAIL-001"`
	Price      float64 `json:"price" binding:"required,gt=0" example:"129.99"`
	Status     string  `json:"status" binding:"required,min=1" example:"active"`
	CategoryID *uint   `json:"categoryId" example:"3"`
}

type UpdateProductRequest struct {
	Name             *string `json:"name" binding:"omitempty,min=1" example:"Trail Running Shoes"`
	StockKeepingUnit *string `json:"stockKeepingUnit" binding:"omitempty,min=1" example:"TRAIL-001"`
	// Deprecated: sku is accepted as a request alias for backward compatibility.
	SKUAlias   *string      `json:"sku" binding:"omitempty,min=1" example:"TRAIL-001"`
	Price      *float64     `json:"price" binding:"omitempty,gt=0" example:"139.99"`
	Status     *string      `json:"status" binding:"omitempty,min=1" example:"active"`
	CategoryID OptionalUint `json:"categoryId"`
}

type OptionalUint struct {
	Set   bool
	Null  bool
	Value uint
}

// ProductCreateExample is a documentation-only example payload for create requests.
type ProductCreateExample struct {
	Name             string  `json:"name" example:"Trail Running Shoes"`
	StockKeepingUnit *string `json:"stockKeepingUnit,omitempty" example:"TRAIL-001"`
	// Deprecated: sku is accepted as a request alias for backward compatibility.
	SKUAlias   *string `json:"sku,omitempty" example:"TRAIL-001"`
	Price      float64 `json:"price" example:"129.99"`
	Status     string  `json:"status" example:"active"`
	CategoryID *uint   `json:"categoryId,omitempty" example:"3"`
}

// ProductUpdateExample is a documentation-only example payload for patch requests.
type ProductUpdateExample struct {
	Name             *string `json:"name,omitempty" example:"Trail Running Shoes"`
	StockKeepingUnit *string `json:"stockKeepingUnit,omitempty" example:"TRAIL-001"`
	// Deprecated: sku is accepted as a request alias for backward compatibility.
	SKUAlias   *string  `json:"sku,omitempty" example:"TRAIL-001"`
	Price      *float64 `json:"price,omitempty" example:"139.99"`
	Status     *string  `json:"status,omitempty" example:"active"`
	CategoryID *uint    `json:"categoryId,omitempty" example:"3"`
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
