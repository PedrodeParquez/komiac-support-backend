package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"komiac-support-backend/internal/auth"
)

type AuthConfig struct {
	AccessSecret string
}

func RequireAuth(cfg AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no token"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(h, "Bearer ")
		claims, err := auth.Parse(token, cfg.AccessSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("uid", claims.UID)
		c.Set("role", string(claims.Role))
		c.Next()
	}
}
