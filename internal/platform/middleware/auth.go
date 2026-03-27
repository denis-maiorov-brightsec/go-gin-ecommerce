package middleware

import (
	"strings"

	commonapi "go-gin-ecommerce/internal/common/api"
	platformauth "go-gin-ecommerce/internal/platform/auth"

	"github.com/gin-gonic/gin"
)

const identityContextKey = "auth.identity"

func RequirePermission(authenticator platformauth.Authenticator, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if authenticator == nil {
			_ = c.Error(commonapi.NewInternalServerError())
			c.Abort()
			return
		}

		token, ok := parseBearerToken(c.GetHeader("Authorization"))
		if !ok {
			_ = c.Error(commonapi.NewUnauthorizedError())
			c.Abort()
			return
		}

		identity, authenticated := authenticator.AuthenticateBearerToken(token)
		if !authenticated {
			_ = c.Error(commonapi.NewUnauthorizedError())
			c.Abort()
			return
		}

		if !identity.HasPermission(permission) {
			_ = c.Error(commonapi.NewForbiddenError())
			c.Abort()
			return
		}

		c.Set(identityContextKey, *identity)
		c.Next()
	}
}

func parseBearerToken(header string) (string, bool) {
	scheme, token, found := strings.Cut(strings.TrimSpace(header), " ")
	if !found || !strings.EqualFold(scheme, "Bearer") {
		return "", false
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return "", false
	}

	return token, true
}
