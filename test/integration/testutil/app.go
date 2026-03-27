package testutil

import (
	"log/slog"
	"testing"

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

func NewRouterWithDB(t *testing.T) *gin.Engine {
	t.Helper()

	cfg := NewTestConfig(t)
	return NewRouterWithConfigAndDB(t, cfg)
}

func NewRouterWithConfigAndDB(t *testing.T, cfg config.Config) *gin.Engine {
	t.Helper()

	database := NewTestDatabase(t, cfg)

	gin.SetMode(gin.TestMode)

	return routes.NewWithDB(cfg, slog.Default(), database)
}
