package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gogogo/internal/pkg/response"
)

// RequireRole 只允许给定角色访问；admin 不会自动绕过，需要在调用处显式加入。
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, role := range roles {
		allowed[role] = true
	}
	return func(c *gin.Context) {
		user, ok := CurrentUser(c)
		if !ok {
			response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "authentication required")
			return
		}
		if !allowed[user.Role] {
			response.Error(c, http.StatusForbidden, response.CodeForbidden, "insufficient permissions")
			return
		}
		c.Next()
	}
}
