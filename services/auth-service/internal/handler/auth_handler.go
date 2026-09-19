package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/votify/auth-service/internal/dto"
	"github.com/votify/auth-service/internal/security"
	"github.com/votify/auth-service/internal/service"
	"github.com/votify/pkg/errors"
	"github.com/votify/pkg/response"

	"github.com/go-playground/validator/v10"
	val "github.com/votify/pkg/validator"
)

type AuthHandler struct {
	service   service.AuthService
	jwtSecret string
}

func NewAuthHandler(svc service.AuthService, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		service:   svc,
		jwtSecret: jwtSecret,
	}
}

// Register handles POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			response.ValidationError(c, val.FormatErrors(ve))
			return
		}
		response.BadRequest(c, "Invalid request payload", nil)
		return
	}

	res, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.Created(c, res)
}

// Login handles POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			response.ValidationError(c, val.FormatErrors(ve))
			return
		}
		response.BadRequest(c, "Invalid request payload", nil)
		return
	}

	res, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.OK(c, res)
}

// Logout handles POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	response.OK(c, gin.H{
		"message": "Successfully logged out",
	})
}

// Me handles GET /auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	// Extract user ID from header forwarded by API Gateway or Bearer token
	userID := c.GetHeader("X-User-ID")

	if userID == "" {
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := security.ValidateJWT(tokenStr, h.jwtSecret)
			if err == nil && claims != nil {
				userID = claims.Sub
			}
		}
	}

	if userID == "" {
		response.Unauthorized(c, "Authentication required")
		return
	}

	res, err := h.service.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		errors.RespondWithError(c, err)
		return
	}

	response.OK(c, res)
}
