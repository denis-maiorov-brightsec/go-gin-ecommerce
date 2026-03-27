package routes

import (
	"log/slog"

	"go-gin-ecommerce/internal/platform/config"
	"go-gin-ecommerce/internal/platform/middleware"

	"github.com/gin-gonic/gin"
)

func New(cfg config.Config, logger *slog.Logger) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(middleware.Recovery(logger))

	// The /v1 namespace is reserved here so spec 001 can add the first versioned routes.
	v1 := router.Group("/v1")
	_ = v1

	return router
}
