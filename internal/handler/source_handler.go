package handler

import (
	"strings"

	"github.com/gin-gonic/gin"

	"gogogo/internal/dto/request"
	"gogogo/internal/model"
	"gogogo/internal/pkg/response"
	"gogogo/internal/service"
)

type SourceHandler struct {
	sources *service.DataSourceService
}

func NewSourceHandler(sources *service.DataSourceService) *SourceHandler {
	return &SourceHandler{sources: sources}
}

var sensitiveKeys = []string{"password", "passwd", "secret", "token", "key", "cert", "username"}

func maskConfig(config map[string]any) map[string]any {
	if config == nil {
		return nil
	}
	masked := make(map[string]any, len(config))
	for key, value := range config {
		lowerKey := strings.ToLower(key)
		isSensitive := false
		for _, sensitive := range sensitiveKeys {
			if strings.Contains(lowerKey, sensitive) {
				isSensitive = true
				break
			}
		}
		if isSensitive {
			masked[key] = "***"
		} else {
			masked[key] = value
		}
	}
	return masked
}

func maskSource(source model.DataSource) model.DataSource {
	source.Config = maskConfig(source.Config)
	return source
}

func (h *SourceHandler) List(c *gin.Context) {
	page, pageSize, err := parsePagination(c)
	if err != nil {
		writeError(c, err)
		return
	}
	items, err := h.sources.List(c.Request.Context())
	if err != nil {
		writeError(c, err)
		return
	}
	for index := range items {
		items[index] = maskSource(items[index])
	}
	pageItems, total := response.Paginate(items, page, pageSize)
	response.OK(c, response.NewPage(pageItems, total, page, pageSize))
}

func (h *SourceHandler) Create(c *gin.Context) {
	var req request.SourceCreate
	if !bindJSON(c, &req) {
		return
	}
	source, err := h.sources.Create(c.Request.Context(), req)
	if err != nil {
		writeError(c, err)
		return
	}
	response.Created(c, maskSource(source))
}

func (h *SourceHandler) Get(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	source, err := h.sources.Get(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, maskSource(source))
}

func (h *SourceHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req request.SourceUpdate
	if !bindJSON(c, &req) {
		return
	}
	source, err := h.sources.Update(c.Request.Context(), id, req)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, maskSource(source))
}

func (h *SourceHandler) Delete(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.sources.Delete(c.Request.Context(), id); err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"id": id, "deleted": true})
}

func (h *SourceHandler) Test(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	result, err := h.sources.Test(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, result)
}
