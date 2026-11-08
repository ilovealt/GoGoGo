package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"gogogo/internal/model"
	"gogogo/internal/pkg/response"
	"gogogo/internal/pkg/utils"
	"gogogo/internal/repository"
)

const (
	UserContextKey   = "current_user"
	ClaimsContextKey = "access_claims"
)

func Auth(store repository.Store, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "missing bearer token")
			return
		}
		raw := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		claims, err := utils.ParseAccessToken(secret, raw)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "invalid or expired access token")
			return
		}
		userID, err := strconv.ParseUint(claims.Subject, 10, 64)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "invalid access token subject")
			return
		}
		user, err := store.GetUserByID(c.Request.Context(), userID)
		if err != nil || user.Status != model.StatusEnabled {
			response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "user is not available")
			return
		}
		if user.TokenVersion != claims.TokenVersion {
			response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "access token has been revoked")
			return
		}

		c.Set(UserContextKey, user)
		c.Set(ClaimsContextKey, claims)
		c.Next()
	}
}

func CurrentUser(c *gin.Context) (model.User, bool) {
	value, ok := c.Get(UserContextKey)
	if !ok {
		return model.User{}, false
	}
	user, ok := value.(model.User)
	return user, ok
}
