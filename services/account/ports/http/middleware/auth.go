package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/services/account/app"
	"github.com/scul0405/saga-orchestration/services/account/domain/valueobject"
	"github.com/scul0405/saga-orchestration/services/account/service/auth"
	"net/http"
)

type JWTAuthMW struct {
	authSvc app.AuthService
	logger  logger.Logger
}

func NewJWTAuthMW(authSvc app.AuthService, logger logger.Logger) *JWTAuthMW {
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

		authPayload := &valueobject.AuthPayload{
			AccessToken: accessToken,
		}

		authResponse, err := m.authSvc.Auth(c, authPayload)
		if err != nil {
			m.logger.Errorf("auth middleware: %v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if authResponse.Expired {
			c.JSON(http.StatusUnauthorized, gin.H{"error": auth.ErrTokenExpired})
			c.Abort()
			return
		}

		c.Set("customer_id", authResponse.CustomerID)
		c.Next()
	}
}
