package handler

import (
	"github.com/gin-gonic/gin"

	"gogogo/internal/pkg/response"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Healthz(c *gin.Context) {
	response.OK(c, gin.H{"status": "ok"})
}

func (h *HealthHandler) Readyz(c *gin.Context) {
	response.OK(c, gin.H{"status": "ok", "storage": "postgres"})
}
