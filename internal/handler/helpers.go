package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"gogogo/internal/pkg/apperrors"
	"gogogo/internal/pkg/response"
)

func parseID(c *gin.Context, name string) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || id == 0 {
		response.Error(c, http.StatusBadRequest, response.CodeInvalid, "invalid "+name)
		return 0, false
	}
	return id, true
}

func bindJSON(c *gin.Context, dst any) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeInvalid, "invalid request body")
		return false
	}
	return true
}

func parsePagination(c *gin.Context) (int, int, error) {
	page, pageSize := 1, 20
	if value := strings.TrimSpace(c.Query("page")); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 {
			return 0, 0, apperrors.New(apperrors.ErrInvalid, "page must be greater than 0")
		}
		page = parsed
	}
	if value := strings.TrimSpace(c.Query("page_size")); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 || parsed > 100 {
			return 0, 0, apperrors.New(apperrors.ErrInvalid, "page_size must be between 1 and 100")
		}
		pageSize = parsed
	}
	return page, pageSize, nil
}

func parseOptionalUintQuery(c *gin.Context, name string) (*uint64, error) {
	value := strings.TrimSpace(c.Query(name))
	if value == "" {
		return nil, nil
	}
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrInvalid, name+" must be a positive integer")
	}
	return &parsed, nil
}

func parseOptionalStatusQuery(c *gin.Context) (*uint8, error) {
	value := strings.TrimSpace(c.Query("status"))
	if value == "" {
		return nil, nil
	}
	parsed, err := strconv.ParseUint(value, 10, 8)
	if err != nil || parsed > 1 {
		return nil, apperrors.New(apperrors.ErrInvalid, "status must be 0 or 1")
	}
	status := uint8(parsed)
	return &status, nil
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		response.Error(c, http.StatusNotFound, response.CodeNotFound, err.Error())
	case errors.Is(err, apperrors.ErrConflict):
		response.Error(c, http.StatusConflict, response.CodeConflict, err.Error())
	case errors.Is(err, apperrors.ErrInvalid):
		response.Error(c, http.StatusBadRequest, response.CodeInvalid, err.Error())
	case errors.Is(err, apperrors.ErrUnauthorized):
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
	case errors.Is(err, apperrors.ErrForbidden):
		response.Error(c, http.StatusForbidden, response.CodeForbidden, err.Error())
	default:
		slog.Error("request failed", "error", err, "path", c.Request.URL.Path, "request_id", c.GetString("request_id"))
		response.Error(c, http.StatusInternalServerError, response.CodeInternal, "internal server error")
	}
}
