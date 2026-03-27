package http

import (
	nethttp "net/http"
	"strconv"

	commonapi "go-gin-ecommerce/internal/common/api"
	"go-gin-ecommerce/internal/promotions/dto"
	"go-gin-ecommerce/internal/promotions/service"

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

func (h *Handler) List(c *gin.Context) {
	promotions, err := h.service.List(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(nethttp.StatusOK, dto.NewPromotionResponses(promotions))
}

func (h *Handler) GetByID(c *gin.Context) {
	id, err := parsePromotionID(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	promotion, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(nethttp.StatusOK, dto.NewPromotionResponse(promotion))
}

func (h *Handler) Create(c *gin.Context) {
	var request dto.CreatePromotionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}

	promotion, err := h.service.Create(c.Request.Context(), request)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(nethttp.StatusCreated, dto.NewPromotionResponse(promotion))
}

func (h *Handler) Update(c *gin.Context) {
	id, err := parsePromotionID(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	var request dto.UpdatePromotionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}

	promotion, err := h.service.Update(c.Request.Context(), id, request)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(nethttp.StatusOK, dto.NewPromotionResponse(promotion))
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := parsePromotionID(c.Param("id"))
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

func parsePromotionID(rawID string) (uint, error) {
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
