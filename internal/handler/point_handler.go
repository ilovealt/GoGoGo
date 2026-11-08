package handler

import (
	"github.com/gin-gonic/gin"

	"gogogo/internal/dto/request"
	"gogogo/internal/pkg/response"
	"gogogo/internal/service"
)

type PointHandler struct {
	points *service.PointService
}

func NewPointHandler(points *service.PointService) *PointHandler {
	return &PointHandler{points: points}
}

func (h *PointHandler) List(c *gin.Context) {
	page, pageSize, err := parsePagination(c)
	if err != nil {
		writeError(c, err)
		return
	}
	equipmentID, err := parseOptionalUintQuery(c, "equipment_id")
	if err != nil {
		writeError(c, err)
		return
	}
	status, err := parseOptionalStatusQuery(c)
	if err != nil {
		writeError(c, err)
		return
	}
	items, err := h.points.List(c.Request.Context(), equipmentID, c.Query("point_type"), status)
	if err != nil {
		writeError(c, err)
		return
	}
	pageItems, total := response.Paginate(items, page, pageSize)
	response.OK(c, response.NewPage(pageItems, total, page, pageSize))
}

func (h *PointHandler) Create(c *gin.Context) {
	var req request.PointCreate
	if !bindJSON(c, &req) {
		return
	}
	point, err := h.points.Create(c.Request.Context(), req)
	if err != nil {
		writeError(c, err)
		return
	}
	response.Created(c, point)
}

func (h *PointHandler) Get(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	point, err := h.points.Get(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, point)
}

func (h *PointHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req request.PointUpdate
	if !bindJSON(c, &req) {
		return
	}
	point, err := h.points.Update(c.Request.Context(), id, req)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, point)
}

func (h *PointHandler) Delete(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.points.Delete(c.Request.Context(), id); err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"id": id, "deleted": true})
}
