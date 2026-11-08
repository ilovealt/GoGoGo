package handler

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"gogogo/internal/pkg/apperrors"
	"gogogo/internal/pkg/response"
	"gogogo/internal/service"
)

type SampleHandler struct {
	samples *service.SampleService
}

func NewSampleHandler(samples *service.SampleService) *SampleHandler {
	return &SampleHandler{samples: samples}
}

func parseTimeQuery(c *gin.Context, name string) (time.Time, error) {
	value := strings.TrimSpace(c.Query(name))
	if value == "" {
		return time.Time{}, apperrors.New(apperrors.ErrInvalid, name+" is required")
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02"} {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed.UTC(), nil
		}
	}
	return time.Time{}, apperrors.New(apperrors.ErrInvalid, name+" must be RFC3339, e.g. 2026-11-03T10:00:00Z")
}

func (h *SampleHandler) Query(c *gin.Context) {
	pointID, ok := parseID(c, "id")
	if !ok {
		return
	}
	start, err := parseTimeQuery(c, "start")
	if err != nil {
		writeError(c, err)
		return
	}
	end, err := parseTimeQuery(c, "end")
	if err != nil {
		writeError(c, err)
		return
	}
	limit := 0
	if value := strings.TrimSpace(c.Query("limit")); value != "" {
		limit, err = strconv.Atoi(value)
		if err != nil {
			writeError(c, apperrors.New(apperrors.ErrInvalid, "limit must be an integer"))
			return
		}
	}
	samples, err := h.samples.Query(c.Request.Context(), pointID, start, end, c.Query("interval"), c.Query("agg"), limit)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{
		"point_id": pointID,
		"samples":  samples,
		"count":    len(samples),
	})
}

func (h *SampleHandler) Latest(c *gin.Context) {
	pointID, ok := parseID(c, "id")
	if !ok {
		return
	}
	latest, err := h.samples.Latest(c.Request.Context(), pointID)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, latest)
}

func (h *SampleHandler) LatestByPlant(c *gin.Context) {
	plantID, ok := parseID(c, "id")
	if !ok {
		return
	}
	items, err := h.samples.LatestByPlant(c.Request.Context(), plantID)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"plant_id": plantID, "items": items, "count": len(items)})
}
