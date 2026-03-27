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
	ID         uint    `json:"id"`
	ProductID  *uint   `json:"productId,omitempty"`
	Name       string  `json:"name"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unitPrice"`
	LineAmount float64 `json:"lineAmount"`
}

type OrderResponse struct {
	ID          uint                `json:"id"`
	Status      string              `json:"status"`
	CustomerID  uint                `json:"customerId"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
	Items       []OrderItemResponse `json:"items"`
	TotalAmount float64             `json:"totalAmount"`
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
