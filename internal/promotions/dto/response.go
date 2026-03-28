package dto

import (
	"time"

	"go-gin-ecommerce/internal/promotions/model"
)

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
