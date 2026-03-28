package dto

import (
	"bytes"
	"encoding/json"
	"time"
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
