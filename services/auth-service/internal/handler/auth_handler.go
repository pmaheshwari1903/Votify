package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles HTTP requests for authentication endpoints.
// It is responsible only for parsing requests, delegating to the service layer,
// and formatting responses. No business logic belongs here.
type AuthHandler struct {
	// service AuthService — will be injected when service layer is implemented
}

// NewAuthHandler creates a new AuthHandler.
// Dependencies will be injected here once the service layer is built.
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// Register handles POST /register
// Flow: Parse DTO → Validate → Delegate to service → Return response
//
// Not yet implemented — returns a placeholder acknowledging the endpoint.
func (h *AuthHandler) Register(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "Registration is not yet implemented.",
		},
	})
}

// Login handles POST /login
// Flow: Parse DTO → Validate → Delegate to service → Return token
//
// Not yet implemented — returns a placeholder acknowledging the endpoint.
func (h *AuthHandler) Login(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "Login is not yet implemented.",
		},
	})
}

// Me handles GET /me
// Returns the authenticated user's profile.
//
// Not yet implemented.
func (h *AuthHandler) Me(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "User profile endpoint is not yet implemented.",
		},
	})
}

// HealthCheck returns the service health status.
func (h *AuthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"service": "auth-service",
			"status":  "healthy",
		},
	})
}
