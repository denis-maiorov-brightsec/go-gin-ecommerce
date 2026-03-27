package middleware

import (
	"sync"
	"time"

	commonapi "go-gin-ecommerce/internal/common/api"

	"github.com/gin-gonic/gin"
)

type WriteRateLimiterConfig struct {
	Limit   int
	Window  time.Duration
	Now     func() time.Time
	KeyFunc func(*gin.Context) string
}

type rateLimitWindow struct {
	Count     int
	ExpiresAt time.Time
}

func NewWriteRateLimiter(cfg WriteRateLimiterConfig) gin.HandlerFunc {
	limit := cfg.Limit
	if limit <= 0 {
		limit = 1
	}

	window := cfg.Window
	if window <= 0 {
		window = time.Minute
	}

	nowFn := cfg.Now
	if nowFn == nil {
		nowFn = time.Now
	}

	keyFn := cfg.KeyFunc
	if keyFn == nil {
		keyFn = defaultWriteRateLimitKey
	}

	var (
		mu      sync.Mutex
		windows = make(map[string]rateLimitWindow)
	)

	return func(c *gin.Context) {
		key := keyFn(c)
		now := nowFn()

		mu.Lock()
		state := windows[key]
		if state.ExpiresAt.IsZero() || !now.Before(state.ExpiresAt) {
			state = rateLimitWindow{
				Count:     0,
				ExpiresAt: now.Add(window),
			}
		}

		if state.Count >= limit {
			mu.Unlock()
			_ = c.Error(commonapi.NewTooManyRequestsError())
			c.Abort()
			return
		}

		state.Count++
		windows[key] = state

		for existingKey, existingState := range windows {
			if !now.Before(existingState.ExpiresAt) {
				delete(windows, existingKey)
			}
		}
		mu.Unlock()

		c.Next()
	}
}

func defaultWriteRateLimitKey(c *gin.Context) string {
	return c.ClientIP() + ":" + c.Request.Method + ":" + c.FullPath()
}
