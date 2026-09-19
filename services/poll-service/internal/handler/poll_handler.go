package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PollHandler handles HTTP requests for poll endpoints.
type PollHandler struct{}

// NewPollHandler creates a new PollHandler.
func NewPollHandler() *PollHandler {
	return &PollHandler{}
}

// Create handles POST /polls
func (h *PollHandler) Create(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "Poll creation is not yet implemented.",
		},
	})
}

// GetByID handles GET /polls/:id
func (h *PollHandler) GetByID(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "Poll retrieval is not yet implemented.",
		},
	})
}

// Update handles PUT /polls/:id
func (h *PollHandler) Update(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "Poll update is not yet implemented.",
		},
	})
}

// Delete handles DELETE /polls/:id
func (h *PollHandler) Delete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "Poll deletion is not yet implemented.",
		},
	})
}

// List handles GET /polls
func (h *PollHandler) List(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "Poll listing is not yet implemented.",
		},
	})
}

// HealthCheck returns service health.
func (h *PollHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"service": "poll-service",
			"status":  "healthy",
		},
	})
}
