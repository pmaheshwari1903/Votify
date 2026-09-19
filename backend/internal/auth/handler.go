package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/votify/backend/internal/errors"
	"github.com/votify/backend/internal/response"
)

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid registration payload", map[string]string{"body": err.Error()})
		return
	}

	res, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.Created(c, res)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid login payload", map[string]string{"body": err.Error()})
		return
	}

	res, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.OK(c, res)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		userID = c.GetHeader("X-User-ID")
	}

	if userID == "" {
		response.Unauthorized(c, "User context missing")
		return
	}

	res, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.OK(c, res)
}
