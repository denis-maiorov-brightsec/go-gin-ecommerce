package http

import (
	nethttp "net/http"
	"strings"

	commonapi "go-gin-ecommerce/internal/common/api"
	commonpagination "go-gin-ecommerce/internal/common/pagination"
	productdto "go-gin-ecommerce/internal/products/dto"
	productservice "go-gin-ecommerce/internal/products/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	productService productservice.Service
}

func NewHandler(productService productservice.Service) *Handler {
	return &Handler{productService: productService}
}

func (h *Handler) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/products", h.SearchProducts)
}

func (h *Handler) SearchProducts(c *gin.Context) {
	query, err := parseSearchQuery(c.Query("q"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	params, err := commonpagination.Parse(c.Request.URL.Query())
	if err != nil {
		_ = c.Error(err)
		return
	}

	products, total, err := h.productService.Search(c.Request.Context(), query, params)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(nethttp.StatusOK, commonpagination.NewResponse(productdto.NewProductResponses(products), params, total))
}

func parseSearchQuery(raw string) (string, error) {
	query := strings.TrimSpace(raw)
	if query == "" {
		return "", commonapi.NewValidationError([]commonapi.ErrorDetail{{
			Field:       "q",
			Constraints: []string{"q is required"},
		}})
	}

	return query, nil
}
