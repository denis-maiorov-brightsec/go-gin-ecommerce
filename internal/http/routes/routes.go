package routes

import (
	"log/slog"
	"net/http"

	commonapi "go-gin-ecommerce/internal/common/api"
	"go-gin-ecommerce/internal/platform/config"
	"go-gin-ecommerce/internal/platform/middleware"

	"github.com/gin-gonic/gin"
)

const (
	deprecationHeader          = "Deprecation"
	successorVersionLinkHeader = "Link"
)

func New(cfg config.Config, logger *slog.Logger) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(middleware.Recovery(logger))

	v1 := router.Group("/v1")
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, commonapi.StatusResponse{Status: "ok"})
	})

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
