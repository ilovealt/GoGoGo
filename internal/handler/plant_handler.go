package handler

import (
	"github.com/gin-gonic/gin"

	"gogogo/internal/dto/request"
	"gogogo/internal/pkg/response"
	"gogogo/internal/service"
)

type PlantHandler struct {
	plants *service.PlantService
}

func NewPlantHandler(plants *service.PlantService) *PlantHandler {
	return &PlantHandler{plants: plants}
}

func (h *PlantHandler) List(c *gin.Context) {
	page, pageSize, err := parsePagination(c)
	if err != nil {
		writeError(c, err)
		return
	}
	items, err := h.plants.List(c.Request.Context())
	if err != nil {
		writeError(c, err)
		return
	}
	pageItems, total := response.Paginate(items, page, pageSize)
	response.OK(c, response.NewPage(pageItems, total, page, pageSize))
}

func (h *PlantHandler) Create(c *gin.Context) {
	var req request.PlantCreate
	if !bindJSON(c, &req) {
		return
	}
	plant, err := h.plants.Create(c.Request.Context(), req)
	if err != nil {
		writeError(c, err)
		return
	}
	response.Created(c, plant)
}

func (h *PlantHandler) Get(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	plant, err := h.plants.Get(c.Request.Context(), id)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, plant)
}

func (h *PlantHandler) Update(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req request.PlantUpdate
	if !bindJSON(c, &req) {
		return
	}
	plant, err := h.plants.Update(c.Request.Context(), id, req)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, plant)
}

func (h *PlantHandler) Delete(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.plants.Delete(c.Request.Context(), id); err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"id": id, "deleted": true})
}
