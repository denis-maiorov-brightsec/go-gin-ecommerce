package testutil

import (
	"log/slog"

	"go-gin-ecommerce/internal/http/routes"
	"go-gin-ecommerce/internal/platform/config"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	return routes.New(config.Config{
		AppEnv:   "test",
		Port:     "0",
		LogLevel: "error",
	}, slog.Default())
}
