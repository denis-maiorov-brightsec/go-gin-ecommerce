package http

import (
	nethttp "net/http"

	"go-gin-ecommerce/internal/orders/commands/service"
	"go-gin-ecommerce/internal/orders/dto"
	queryhttp "go-gin-ecommerce/internal/orders/queries/http"

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

func (h *Handler) Cancel(c *gin.Context) {
	id, err := queryhttp.ParseOrderID(c.Param("id"))
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
