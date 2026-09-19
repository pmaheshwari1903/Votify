package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// VoteHandler handles HTTP requests for vote endpoints.
type VoteHandler struct{}

// NewVoteHandler creates a new VoteHandler.
func NewVoteHandler() *VoteHandler {
	return &VoteHandler{}
}

// CastVote handles POST /votes
// Flow: Parse DTO → Validate → Check poll exists & is published →
//       Save vote → Publish VoteCreated to Kafka → Return confirmation
func (h *VoteHandler) CastVote(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "Vote casting is not yet implemented.",
		},
	})
}

// GetTally handles GET /votes/tally/:pollId
// Returns aggregated vote counts for a poll.
func (h *VoteHandler) GetTally(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "Vote tally is not yet implemented.",
		},
	})
}

// HealthCheck returns service health.
func (h *VoteHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"service": "vote-service",
			"status":  "healthy",
		},
	})
}
