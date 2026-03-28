package http

import (
	"fmt"
	nethttp "net/http"
	"net/url"
	"time"

	commonapi "go-gin-ecommerce/internal/common/api"
	commonpagination "go-gin-ecommerce/internal/common/pagination"
	"go-gin-ecommerce/internal/orders/dto"
	ordershttp "go-gin-ecommerce/internal/orders/http"
	"go-gin-ecommerce/internal/orders/queries/service"

	"github.com/gin-gonic/gin"
)

const dateLayout = "2006-01-02"

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("", h.List)
	group.GET("/:id", h.GetByID)
}

// List godoc
// @Summary List orders
// @Description Returns a paginated list of orders with optional status and date filters.
// @Tags orders
// @Produce json
// @Param page query int false "Page number" minimum(1) default(1)
// @Param limit query int false "Page size" minimum(1) maximum(100) default(20)
// @Param status query string false "Order status filter" example(pending)
// @Param from query string false "Start date filter in YYYY-MM-DD format" example(2026-03-01)
// @Param to query string false "End date filter in YYYY-MM-DD format" example(2026-03-31)
// @Success 200 {object} dto.OrderListResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /orders [get]
func (h *Handler) List(c *gin.Context) {
	params, err := parseListOrdersParams(c.Request.URL.Query())
	if err != nil {
		_ = c.Error(err)
		return
	}

	orders, total, err := h.service.List(c.Request.Context(), params)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(nethttp.StatusOK, commonpagination.NewResponse(dto.NewOrderResponses(orders), params.Pagination, total))
}

// GetByID godoc
// @Summary Get order
// @Description Returns a single order with line items by numeric identifier.
// @Tags orders
// @Produce json
// @Param id path int true "Order ID" minimum(1)
// @Success 200 {object} dto.OrderResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /orders/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := ordershttp.ParseOrderID(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	order, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(nethttp.StatusOK, dto.NewOrderResponse(order))
}

func parseListOrdersParams(values url.Values) (dto.ListOrdersParams, error) {
	paginationParams, err := commonpagination.Parse(values)
	if err != nil {
		return dto.ListOrdersParams{}, err
	}

	from, err := parseDateFilter(values.Get("from"), "from")
	if err != nil {
		return dto.ListOrdersParams{}, err
	}

	to, err := parseDateFilter(values.Get("to"), "to")
	if err != nil {
		return dto.ListOrdersParams{}, err
	}
	if to != nil {
		toValue := to.Add(24 * time.Hour)
		to = &toValue
	}

	if from != nil && to != nil && from.After(to.Add(-24*time.Hour)) {
		return dto.ListOrdersParams{}, commonapi.NewValidationError([]commonapi.ErrorDetail{{
			Field:       "from",
			Constraints: []string{"from must be before or equal to to"},
		}})
	}

	return dto.ListOrdersParams{
		Pagination: paginationParams,
		Status:     values.Get("status"),
		From:       from,
		To:         to,
	}, nil
}

func parseDateFilter(raw string, field string) (*time.Time, error) {
	if raw == "" {
		return nil, nil
	}

	value, err := time.Parse(dateLayout, raw)
	if err != nil {
		return nil, commonapi.NewValidationError([]commonapi.ErrorDetail{{
			Field:       field,
			Constraints: []string{fmt.Sprintf("%s must be a valid date in %s format", field, dateLayout)},
		}})
	}

	utcValue := value.UTC()
	return &utcValue, nil
}
