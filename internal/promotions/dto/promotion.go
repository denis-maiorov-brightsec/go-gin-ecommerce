package dto

import (
	"bytes"
	"encoding/json"
	"time"

	"go-gin-ecommerce/internal/promotions/model"
)

type CreatePromotionRequest struct {
	Name          string     `json:"name" binding:"required,min=1"`
	Code          string     `json:"code" binding:"required,min=1"`
	DiscountType  string     `json:"discountType" binding:"required,min=1"`
	DiscountValue float64    `json:"discountValue" binding:"required,gt=0"`
	StartsAt      *time.Time `json:"startsAt"`
	EndsAt        *time.Time `json:"endsAt"`
	Status        string     `json:"status" binding:"required,min=1"`
}

type UpdatePromotionRequest struct {
	Name          *string      `json:"name" binding:"omitempty,min=1"`
	Code          *string      `json:"code" binding:"omitempty,min=1"`
	DiscountType  *string      `json:"discountType" binding:"omitempty,min=1"`
	DiscountValue *float64     `json:"discountValue" binding:"omitempty,gt=0"`
	StartsAt      OptionalTime `json:"startsAt"`
	EndsAt        OptionalTime `json:"endsAt"`
	Status        *string      `json:"status" binding:"omitempty,min=1"`
}

type PromotionResponse struct {
	ID            uint       `json:"id"`
	Name          string     `json:"name"`
	Code          string     `json:"code"`
	DiscountType  string     `json:"discountType"`
	DiscountValue float64    `json:"discountValue"`
	StartsAt      *time.Time `json:"startsAt,omitempty"`
	EndsAt        *time.Time `json:"endsAt,omitempty"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

func NewPromotionResponse(promotion model.Promotion) PromotionResponse {
	return PromotionResponse{
		ID:            promotion.ID,
		Name:          promotion.Name,
		Code:          promotion.Code,
		DiscountType:  promotion.DiscountType,
		DiscountValue: promotion.DiscountValue,
		StartsAt:      promotion.StartsAt,
		EndsAt:        promotion.EndsAt,
		Status:        promotion.Status,
		CreatedAt:     promotion.CreatedAt,
		UpdatedAt:     promotion.UpdatedAt,
	}
}

func NewPromotionResponses(promotions []model.Promotion) []PromotionResponse {
	responses := make([]PromotionResponse, 0, len(promotions))
	for _, promotion := range promotions {
		responses = append(responses, NewPromotionResponse(promotion))
	}

	return responses
}

type OptionalTime struct {
	Set   bool
	Null  bool
	Value time.Time
}

func (o *OptionalTime) UnmarshalJSON(data []byte) error {
	o.Set = true

	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		o.Null = true
		o.Value = time.Time{}
		return nil
	}

	o.Null = false
	return json.Unmarshal(data, &o.Value)
}
