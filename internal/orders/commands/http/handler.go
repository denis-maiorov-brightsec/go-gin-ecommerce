package http

import (
	nethttp "net/http"

	"go-gin-ecommerce/internal/orders/commands/service"
	"go-gin-ecommerce/internal/orders/dto"
	ordershttp "go-gin-ecommerce/internal/orders/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(group *gin.RouterGroup, writeMiddlewares ...gin.HandlerFunc) {
	group.POST("/:id/cancel", append(writeMiddlewares, h.Cancel)...)
}

// Cancel godoc
// @Summary Cancel order
// @Description Executes the order cancel state transition for an existing order and returns the updated order resource.
// @Tags orders
// @Accept json
// @Produce json
// @Param id path int true "Order ID" minimum(1)
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Failure 409 {object} api.ErrorResponse
// @Failure 429 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /orders/{id}/cancel [post]
func (h *Handler) Cancel(c *gin.Context) {
	id, err := ordershttp.ParseOrderID(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	order, err := h.service.Cancel(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(nethttp.StatusOK, dto.NewOrderResponse(order))
}
