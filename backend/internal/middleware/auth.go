package middleware

import (
	"net/http"
	"strings"

	"stock-trading/internal/response"
	"stock-trading/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := ""

		// 1. Try to get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// 2. Fallback to HttpOnly cookie
		if tokenString == "" {
			cookie, err := c.Cookie("access_token")
			if err == nil {
				tokenString = cookie
			}
		}

		if tokenString == "" {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing authentication token", nil)
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired token", err.Error())
			c.Abort()
			return
		}

		// Store user data in context for subsequent handlers
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role.(string) != requiredRole {
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "You do not have permission to access this resource", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
