package middlewares

import (
	"net/http"
	"strings"

	"github.com/Rajkumar-coderm/go-blog-backend/internal/auth"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware protects routes by verifying JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := auth.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Ensure only access tokens are allowed
		if claims.Type != "access" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type"})
			c.Abort()
			return
		}

		// Store user ID in context
		c.Set("userID", claims.ID)
		c.Next()
	}
}
