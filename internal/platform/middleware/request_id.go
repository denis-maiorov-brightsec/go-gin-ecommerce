package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

const RequestIDHeader = "X-Request-ID"

type contextKey string

const requestIDContextKey contextKey = "request_id"

const ginRequestIDKey = "request_id"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = newRequestID()
		}

		ctx := context.WithValue(c.Request.Context(), requestIDContextKey, requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Set(ginRequestIDKey, requestID)
		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}

func GetRequestID(c *gin.Context) string {
	if value, ok := c.Get(ginRequestIDKey); ok {
		if requestID, valid := value.(string); valid {
			return requestID
		}
	}

	requestID, _ := RequestIDFromContext(c.Request.Context())
	return requestID
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(requestIDContextKey).(string)
	return requestID, ok && requestID != ""
}

func newRequestID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "request-id-unavailable"
	}

	return hex.EncodeToString(bytes)
}
