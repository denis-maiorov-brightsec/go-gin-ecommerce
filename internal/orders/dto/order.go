package dto

import (
	"time"

	commonpagination "go-gin-ecommerce/internal/common/pagination"
	"go-gin-ecommerce/internal/orders/model"
)

type ListOrdersParams struct {
	Pagination commonpagination.Params
	Status     string
	From       *time.Time
	To         *time.Time
}

type OrderItemResponse struct {
	ID         uint    `json:"id" example:"501"`
	ProductID  *uint   `json:"productId,omitempty" example:"42"`
	Name       string  `json:"name" example:"Trail Running Shoes"`
	Quantity   int     `json:"quantity" example:"2"`
	UnitPrice  float64 `json:"unitPrice" example:"129.99"`
	LineAmount float64 `json:"lineAmount" example:"259.98"`
}

type OrderResponse struct {
	ID          uint                `json:"id" example:"1001"`
	Status      string              `json:"status" example:"pending"`
	CustomerID  uint                `json:"customerId" example:"88"`
	CreatedAt   time.Time           `json:"createdAt" example:"2026-03-01T10:00:00Z"`
	UpdatedAt   time.Time           `json:"updatedAt" example:"2026-03-02T11:00:00Z"`
	Items       []OrderItemResponse `json:"items"`
	TotalAmount float64             `json:"totalAmount" example:"259.98"`
}

type OrderListResponse struct {
	Items      []OrderResponse `json:"items"`
	Page       int             `json:"page" example:"1"`
	Limit      int             `json:"limit" example:"20"`
	Total      int64           `json:"total" example:"1"`
	TotalPages int             `json:"totalPages" example:"1"`
}

func NewOrderResponse(order model.Order) OrderResponse {
	return OrderResponse{
		ID:          order.ID,
		Status:      order.Status,
		CustomerID:  order.CustomerID,
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
		Items:       NewOrderItemResponses(order.Items),
		TotalAmount: order.TotalAmount,
	}
}

func NewOrderResponses(orders []model.Order) []OrderResponse {
	responses := make([]OrderResponse, 0, len(orders))
	for _, order := range orders {
		responses = append(responses, NewOrderResponse(order))
	}

	return responses
}

func NewOrderItemResponses(items []model.OrderItem) []OrderItemResponse {
	responses := make([]OrderItemResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, OrderItemResponse{
			ID:         item.ID,
			ProductID:  item.ProductID,
			Name:       item.Name,
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
			LineAmount: item.LineAmount,
		})
	}

	return responses
}
