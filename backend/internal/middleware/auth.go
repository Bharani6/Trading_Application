package middleware

import (
	"net/http"
	"strings"

	"stock-trading/internal/response"
	"stock-trading/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(pCtx *gin.Context) {
		lTokenString := ""

		// 1. Try to get token from Authorization header
		lAuthHeader := pCtx.GetHeader("Authorization")
		if lAuthHeader != "" && strings.HasPrefix(lAuthHeader, "Bearer ") {
			lTokenString = strings.TrimPrefix(lAuthHeader, "Bearer ")
		}

		// 2. Fallback to HttpOnly cookie
		if lTokenString == "" {
			lCookie, lErr := pCtx.Cookie("access_token")
			if lErr == nil {
				lTokenString = lCookie
			}
		}

		if lTokenString == "" {
			response.Error(pCtx, http.StatusUnauthorized, "UNAUTHORIZED", "Missing authentication token", nil)
			pCtx.Abort()
			return
		}

		lClaims, lErr := utils.ValidateToken(lTokenString)
		if lErr != nil {
			response.Error(pCtx, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or expired token", lErr.Error())
			pCtx.Abort()
			return
		}

		// Store user data in context for subsequent handlers
		pCtx.Set("userID", lClaims.UserID)
		pCtx.Set("role", lClaims.Role)

		pCtx.Next()
	}
}

func RoleMiddleware(pRequiredRole string) gin.HandlerFunc {
	return func(pCtx *gin.Context) {
		lRole, lExists := pCtx.Get("role")
		if !lExists || lRole.(string) != pRequiredRole {
			response.Error(pCtx, http.StatusForbidden, "FORBIDDEN", "You do not have permission to access this resource", nil)
			pCtx.Abort()
			return
		}
		pCtx.Next()
	}
}
