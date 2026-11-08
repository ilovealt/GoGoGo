package handler

import (
	"github.com/gin-gonic/gin"

	"gogogo/internal/dto/request"
	"gogogo/internal/model"
	"gogogo/internal/pkg/response"
	"gogogo/internal/service"
)

type AreaHandler struct {
	areas *service.AreaService
}

func NewAreaHandler(areas *service.AreaService) *AreaHandler {
	return &AreaHandler{areas: areas}
}

type areaNode struct {
	model.Area
	Children []*areaNode `json:"children"`
}

func buildAreaTree(areas []model.Area) []*areaNode {
	nodes := make(map[uint64]*areaNode, len(areas))
	for _, area := range areas {
		nodes[area.ID] = &areaNode{Area: area, Children: []*areaNode{}}
	}
	roots := make([]*areaNode, 0)
	for _, area := range areas {
		node := nodes[area.ID]
		if area.ParentID == nil {
			roots = append(roots, node)
			continue
		}
		parent, ok := nodes[*area.ParentID]
		if !ok {
			roots = append(roots, node)
			continue
		}
		parent.Children = append(parent.Children, node)
	}
	return roots
}

func (h *AreaHandler) ListByPlant(c *gin.Context) {
	plantID, ok := parseID(c, "id")
	if !ok {
		return
	}
	areas, err := h.areas.ListByPlant(c.Request.Context(), plantID)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, buildAreaTree(areas))
}

func (h *AreaHandler) Create(c *gin.Context) {
	var req request.AreaCreate
	if !bindJSON(c, &req) {
		return
	}
	area, err := h.areas.Create(c.Request.Context(), req)
	if err != nil {
		writeError(c, err)
		return
	}
	response.Created(c, area)
}

func (h *AreaHandler) Get(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	area, err := h.areas.Get(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, area)
}

func (h *AreaHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req request.AreaUpdate
	if !bindJSON(c, &req) {
		return
	}
	area, err := h.areas.Update(c.Request.Context(), id, req)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, area)
}

func (h *AreaHandler) Delete(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.areas.Delete(c.Request.Context(), id); err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"id": id, "deleted": true})
}
