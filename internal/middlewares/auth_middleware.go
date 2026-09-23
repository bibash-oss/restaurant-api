package middlewares

import (
	"strings"

	"kitchen-api/internal/helper"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{
				"success": false,
				"error":   "Authorization header is required",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(401, gin.H{
				"success": false,
				"error":   "Authorization header format must be Bearer {token}",
			})
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(parts[1])
		if err != nil {
			c.JSON(401, gin.H{
				"success": false,
				"error":   "Invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set("userId", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		if claims.RestaurantID != nil {
			c.Set("restaurantId", *claims.RestaurantID)
		}

		c.Next()
	}
}
