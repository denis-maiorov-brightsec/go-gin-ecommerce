package routes

import (
	"log/slog"
	"net/http"

	categoryhttp "go-gin-ecommerce/internal/categories/http"
	categoryrepository "go-gin-ecommerce/internal/categories/repository"
	categoryservice "go-gin-ecommerce/internal/categories/service"
	commonapi "go-gin-ecommerce/internal/common/api"
	orderhttp "go-gin-ecommerce/internal/orders/http"
	orderrepository "go-gin-ecommerce/internal/orders/repository"
	orderservice "go-gin-ecommerce/internal/orders/service"
	platformauth "go-gin-ecommerce/internal/platform/auth"
	"go-gin-ecommerce/internal/platform/config"
	"go-gin-ecommerce/internal/platform/middleware"
	producthttp "go-gin-ecommerce/internal/products/http"
	productrepository "go-gin-ecommerce/internal/products/repository"
	productservice "go-gin-ecommerce/internal/products/service"
	promotionhttp "go-gin-ecommerce/internal/promotions/http"
	promotionrepository "go-gin-ecommerce/internal/promotions/repository"
	promotionservice "go-gin-ecommerce/internal/promotions/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	deprecationHeader          = "Deprecation"
	successorVersionLinkHeader = "Link"
)

func New(cfg config.Config, logger *slog.Logger) *gin.Engine {
	return NewWithDB(cfg, logger, nil)
}

func NewWithDB(cfg config.Config, logger *slog.Logger, database *gorm.DB) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	middleware.SetupValidation()
	writeRateLimiter := middleware.NewWriteRateLimiter(middleware.WriteRateLimiterConfig{
		Limit:  cfg.EffectiveWriteRateLimitRequests(),
		Window: cfg.EffectiveWriteRateLimitWindow(),
	})

	router := gin.New()
	router.Use(middleware.RequestID())
	router.Use(middleware.RequestLogger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.ErrorHandler(logger))
	router.NoRoute(func(c *gin.Context) {
		apiErr := commonapi.NewNotFoundError()
		c.AbortWithStatusJSON(apiErr.Status, commonapi.NewErrorResponse(c.Request.URL.Path, apiErr))
	})

	v1 := router.Group("/v1")
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, commonapi.StatusResponse{Status: "ok"})
	})
	if database != nil {
		authenticator := platformauth.NewStubAuthenticator()

		categoryHandler := categoryhttp.NewHandler(categoryservice.New(categoryrepository.New(database)))
		categoryHandler.RegisterRoutes(v1.Group("/categories"), writeRateLimiter)

		orderHandler := orderhttp.NewHandler(orderservice.New(orderrepository.New(database)))
		orderHandler.RegisterRoutes(v1.Group("/orders"), writeRateLimiter)

		promotionHandler := promotionhttp.NewHandler(promotionservice.New(promotionrepository.New(database)))
		promotionsGroup := v1.Group("/promotions")
		promotionsGroup.Use(middleware.RequirePermission(authenticator, platformauth.PermissionManagePromotions))
		promotionHandler.RegisterRoutes(promotionsGroup, writeRateLimiter)

		productHandler := producthttp.NewHandler(productservice.New(productrepository.New(database)))
		productHandler.RegisterRoutes(v1.Group("/products"), writeRateLimiter)
	}

	// Keep the root route temporarily for transition while directing clients to /v1.
	router.GET("/", func(c *gin.Context) {
		c.Header(deprecationHeader, "true")
		c.Header(successorVersionLinkHeader, `</v1/health>; rel="successor-version"`)
		c.JSON(http.StatusOK, commonapi.MessageResponse{
			Message: "The unversioned root route is deprecated. Migrate to /v1/health.",
		})
	})

	return router
}
