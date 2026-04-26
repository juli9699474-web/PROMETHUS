package router

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func apiKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		required := strings.TrimSpace(os.Getenv("PROMETHEUS_API_KEY"))
		if required == "" {
			c.Next()
			return
		}
		provided := strings.TrimSpace(c.GetHeader("X-API-Key"))
		if provided == "" || provided != required {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid api key"})
			return
		}
		c.Next()
	}
}
