package handler

import (
	"github.com/gin-gonic/gin"

	"gogogo/internal/dto/request"
	"gogogo/internal/pkg/response"
	"gogogo/internal/service"
)

type EquipmentHandler struct {
	equipments *service.EquipmentService
}

func NewEquipmentHandler(equipments *service.EquipmentService) *EquipmentHandler {
	return &EquipmentHandler{equipments: equipments}
}

func (h *EquipmentHandler) List(c *gin.Context) {
	page, pageSize, err := parsePagination(c)
	if err != nil {
		writeError(c, err)
		return
	}
	areaID, err := parseOptionalUintQuery(c, "area_id")
	if err != nil {
		writeError(c, err)
		return
	}
	items, err := h.equipments.List(c.Request.Context(), areaID, c.Query("keyword"))
	if err != nil {
		writeError(c, err)
		return
	}
	pageItems, total := response.Paginate(items, page, pageSize)
	response.OK(c, response.NewPage(pageItems, total, page, pageSize))
}

func (h *EquipmentHandler) Create(c *gin.Context) {
	var req request.EquipmentCreate
	if !bindJSON(c, &req) {
		return
	}
	equipment, err := h.equipments.Create(c.Request.Context(), req)
	if err != nil {
		writeError(c, err)
		return
	}
	response.Created(c, equipment)
}

func (h *EquipmentHandler) Get(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	equipment, err := h.equipments.Get(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, equipment)
}

func (h *EquipmentHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req request.EquipmentUpdate
	if !bindJSON(c, &req) {
		return
	}
	equipment, err := h.equipments.Update(c.Request.Context(), id, req)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, equipment)
}

func (h *EquipmentHandler) Delete(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.equipments.Delete(c.Request.Context(), id); err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"id": id, "deleted": true})
}
