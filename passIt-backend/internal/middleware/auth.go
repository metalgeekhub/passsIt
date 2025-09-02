package middleware

import (
	"log"
	"net/http"
	"passIt/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthenticationMiddleware checks if the user has a valid JWT token
func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authentication token"})
			c.Abort()
			return
		}

		// The Token should be prefixted with Bearer
		tokenParts := strings.Split(tokenString, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication token"})
			c.Abort()
			return
		}

		tokenString = tokenParts[1]

		claims, err := utils.VerifyToken(tokenString)
		if err != nil {
			log.Printf("Error validating the JWT token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims["user_id"])
		c.Next()
	}
}
