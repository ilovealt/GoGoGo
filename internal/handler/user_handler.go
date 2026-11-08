package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gogogo/internal/dto/request"
	dtoresponse "gogogo/internal/dto/response"
	"gogogo/internal/middleware"
	"gogogo/internal/pkg/response"
	"gogogo/internal/service"
)

type UserHandler struct {
	users *service.UserService
}

func NewUserHandler(users *service.UserService) *UserHandler {
	return &UserHandler{users: users}
}

func (h *UserHandler) Me(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "authentication required")
		return
	}
	response.OK(c, dtoresponse.NewUser(user))
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "authentication required")
		return
	}
	var req request.UpdateProfile
	if !bindJSON(c, &req) {
		return
	}
	updated, err := h.users.UpdateProfile(c.Request.Context(), user.ID, req.Nickname)
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, dtoresponse.NewUser(updated))
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "authentication required")
		return
	}
	var req request.ChangePassword
	if !bindJSON(c, &req) {
		return
	}
	if err := h.users.ChangePassword(c.Request.Context(), user.ID, req.OldPassword, req.NewPassword); err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"message": "password changed"})
}
