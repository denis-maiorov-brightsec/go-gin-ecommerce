package middleware

import (
	"log/slog"

	commonapi "go-gin-ecommerce/internal/common/api"

	"github.com/gin-gonic/gin"
)

func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		logger.Error("panic recovered", "panic", recovered, "path", c.Request.URL.Path)
		apiErr := commonapi.NewInternalServerError()
		c.AbortWithStatusJSON(apiErr.Status, commonapi.NewErrorResponse(c.Request.URL.Path, apiErr))
	})
}
