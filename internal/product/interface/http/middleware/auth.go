package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/scul0405/saga-orchestration/internal/product/service"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"net/http"
)

const (
	ErrTokenExpired = "token expired"
)

type JWTAuthMW struct {
	authSvc service.AuthService
	logger  logger.Logger
}

func NewJWTAuthMW(authSvc service.AuthService, logger logger.Logger) *JWTAuthMW {
	return &JWTAuthMW{
		authSvc: authSvc,
		logger:  logger,
	}
}

func (m *JWTAuthMW) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")
		if bearerToken == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		accessToken := bearerToken[7:]
		if accessToken == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		authResponse, err := m.authSvc.Auth(c.Request.Context(), accessToken)
		if err != nil {
			m.logger.Errorf("auth middleware: %v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if authResponse.Expired {
			c.JSON(http.StatusUnauthorized, gin.H{"error": ErrTokenExpired})
			c.Abort()
			return
		}

		c.Set("customer_id", authResponse.CustomerID)
		c.Next()
	}
}
