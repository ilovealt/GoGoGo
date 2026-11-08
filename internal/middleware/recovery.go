package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"gogogo/internal/pkg/response"
)

func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				logger.Error("panic recovered",
					"error", recovered,
					"path", c.Request.URL.Path,
					"request_id", c.GetString(RequestIDKey),
					"stack", string(debug.Stack()),
				)
				if !c.Writer.Written() {
					response.Error(c, http.StatusInternalServerError, response.CodeInternal, "internal server error")
				} else {
					c.Abort()
				}
			}
		}()
		c.Next()
	}
}
