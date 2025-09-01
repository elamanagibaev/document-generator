package httpapi

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Authenticator interface {
	Authenticate(token string) bool
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := uuid.NewString()
		c.Set("request_id", reqID)
		c.Writer.Header().Set("X-Request-ID", reqID)
		c.Next()
	}
}

func LoggerMiddleware(logger func(msg string, args ...any)) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		logger("[%s] %s %s %d %s",
			c.GetString("request_id"),
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			latency,
		)
	}
}

func AuthMiddleware(authenticator Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !authenticator.Authenticate(auth) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}
		c.Next()
	}
}
