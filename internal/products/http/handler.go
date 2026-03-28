package http

import (
	nethttp "net/http"
	"strconv"

	commonapi "go-gin-ecommerce/internal/common/api"
	commonpagination "go-gin-ecommerce/internal/common/pagination"
	"go-gin-ecommerce/internal/products/dto"
	"go-gin-ecommerce/internal/products/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service service.Service
}

func NewHandler(service service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(group *gin.RouterGroup, writeMiddlewares ...gin.HandlerFunc) {
	group.GET("", h.List)
	group.GET("/:id", h.GetByID)
	group.POST("", append(writeMiddlewares, h.Create)...)
	group.PATCH("/:id", append(writeMiddlewares, h.Update)...)
	group.DELETE("/:id", append(writeMiddlewares, h.Delete)...)
}

// List godoc
// @Summary List products
// @Description Returns a paginated list of products for the v1 backoffice API.
// @Tags products
// @Produce json
// @Param page query int false "Page number" minimum(1) default(1)
// @Param limit query int false "Page size" minimum(1) maximum(100) default(20)
// @Success 200 {object} dto.ProductListResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /products [get]
func (h *Handler) List(c *gin.Context) {
	params, err := commonpagination.Parse(c.Request.URL.Query())
	if err != nil {
		_ = c.Error(err)
		return
	}

	products, total, err := h.service.List(c.Request.Context(), params)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(nethttp.StatusOK, commonpagination.NewResponse(dto.NewProductResponses(products), params, total))
}

// GetByID godoc
// @Summary Get product
// @Description Returns a single product by numeric identifier.
// @Tags products
// @Produce json
// @Param id path int true "Product ID" minimum(1)
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /products/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := parseProductID(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	product, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(nethttp.StatusOK, dto.NewProductResponse(product))
}

// Create godoc
// @Summary Create product
// @Description Creates a product. The deprecated `sku` field is accepted as an alias for `stockKeepingUnit`; if both are sent they must match.
// @Tags products
// @Accept json
// @Produce json
// @Param request body dto.ProductCreateExample true "Create product payload"
// @Success 201 {object} dto.ProductResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 409 {object} api.ErrorResponse
// @Failure 429 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /products [post]
func (h *Handler) Create(c *gin.Context) {
	var request dto.CreateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}

	product, err := h.service.Create(c.Request.Context(), request)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(nethttp.StatusCreated, dto.NewProductResponse(product))
}

// Update godoc
// @Summary Update product
// @Description Partially updates a product. The deprecated `sku` field is accepted as an alias for `stockKeepingUnit`; if both are sent they must match.
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID" minimum(1)
// @Param request body dto.ProductUpdateExample true "Update product payload"
// @Success 200 {object} dto.ProductResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Failure 409 {object} api.ErrorResponse
// @Failure 429 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /products/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	id, err := parseProductID(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	var request dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}

	product, err := h.service.Update(c.Request.Context(), id, request)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(nethttp.StatusOK, dto.NewProductResponse(product))
}

// Delete godoc
// @Summary Delete product
// @Description Deletes a product by numeric identifier.
// @Tags products
// @Param id path int true "Product ID" minimum(1)
// @Success 204
// @Failure 400 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Failure 429 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /products/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := parseProductID(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(nethttp.StatusNoContent)
}

func parseProductID(rawID string) (uint, error) {
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		return 0, commonapi.NewValidationError([]commonapi.ErrorDetail{{
			Field:       "id",
			Constraints: []string{"id must be a positive integer"},
		}})
	}

	if strconv.IntSize == 32 && id > uint64(^uint32(0)) {
		return 0, commonapi.NewValidationError([]commonapi.ErrorDetail{{
			Field:       "id",
			Constraints: []string{"id must be a positive integer"},
		}})
	}

	return uint(id), nil
}
