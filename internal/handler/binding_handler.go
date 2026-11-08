package handler

import (
	"github.com/gin-gonic/gin"

	"gogogo/internal/dto/request"
	"gogogo/internal/pkg/response"
	"gogogo/internal/service"
)

type BindingHandler struct {
	bindings *service.BindingService
}

func NewBindingHandler(bindings *service.BindingService) *BindingHandler {
	return &BindingHandler{bindings: bindings}
}

func (h *BindingHandler) List(c *gin.Context) {
	page, pageSize, err := parsePagination(c)
	if err != nil {
		writeError(c, err)
		return
	}
	pointID, err := parseOptionalUintQuery(c, "point_id")
	if err != nil {
		writeError(c, err)
		return
	}
	sourceID, err := parseOptionalUintQuery(c, "source_id")
	if err != nil {
		writeError(c, err)
		return
	}
	items, err := h.bindings.List(c.Request.Context(), pointID, sourceID)
	if err != nil {
		writeError(c, err)
		return
	}
	pageItems, total := response.Paginate(items, page, pageSize)
	response.OK(c, response.NewPage(pageItems, total, page, pageSize))
}

func (h *BindingHandler) Create(c *gin.Context) {
	var req request.BindingCreate
	if !bindJSON(c, &req) {
		return
	}
	binding, err := h.bindings.Create(c.Request.Context(), req)
	if err != nil {
		writeError(c, err)
		return
	}
	response.Created(c, binding)
}

func (h *BindingHandler) Get(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	binding, err := h.bindings.Get(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, binding)
}

func (h *BindingHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req request.BindingUpdate
	if !bindJSON(c, &req) {
		return
	}
	binding, err := h.bindings.Update(c.Request.Context(), id, req)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, binding)
}

func (h *BindingHandler) Delete(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.bindings.Delete(c.Request.Context(), id); err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"id": id, "deleted": true})
}
