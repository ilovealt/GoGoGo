package handler

import (
	"github.com/gin-gonic/gin"

	"gogogo/internal/dto/request"
	dtoresponse "gogogo/internal/dto/response"
	"gogogo/internal/pkg/response"
	"gogogo/internal/service"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.Register
	if !bindJSON(c, &req) {
		return
	}
	user, err := h.auth.Register(c.Request.Context(), req)
	if err != nil {
		writeError(c, err)
		return
	}
	response.Created(c, dtoresponse.NewUser(user))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.Login
	if !bindJSON(c, &req) {
		return
	}
	pair, user, err := h.auth.Login(c.Request.Context(), req.Username, req.Password, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{
		"access_token":  pair.AccessToken,
		"refresh_token": pair.RefreshToken,
		"token_type":    pair.TokenType,
		"expires_in":    pair.ExpiresIn,
		"user":          dtoresponse.NewUser(user),
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req request.Refresh
	if !bindJSON(c, &req) {
		return
	}
	pair, err := h.auth.Refresh(c.Request.Context(), req.RefreshToken, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{
		"access_token":  pair.AccessToken,
		"refresh_token": pair.RefreshToken,
		"token_type":    pair.TokenType,
		"expires_in":    pair.ExpiresIn,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req request.Refresh
	if err := c.ShouldBindJSON(&req); err != nil {
		// 允许请求体为空：此时只返回成功，不撤销任何令牌。
		response.OK(c, gin.H{"message": "logged out"})
		return
	}
	if err := h.auth.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		writeError(c, err)
		return
	}
	response.OK(c, gin.H{"message": "logged out"})
}
