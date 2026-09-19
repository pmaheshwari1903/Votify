package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RealtimeHandler handles HTTP and WebSocket requests for the Realtime Service.
type RealtimeHandler struct {
	// hub *websocket.Hub — will be injected when WebSocket hub is implemented
}

// NewRealtimeHandler creates a new RealtimeHandler.
func NewRealtimeHandler() *RealtimeHandler {
	return &RealtimeHandler{}
}

// WebSocketUpgrade handles GET /ws/:pollId
// This endpoint upgrades HTTP connections to WebSocket for real-time poll updates.
//
// Future flow:
// 1. Client connects to /ws/:pollId
// 2. Connection is registered in the WebSocket hub for that poll
// 3. When VoteCreated events arrive via Kafka → Redis pub/sub,
//    the hub broadcasts updated tallies to all connected clients
// 4. Client receives live updates without page refresh
//
// Not yet implemented.
func (h *RealtimeHandler) WebSocketUpgrade(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "WebSocket connections are not yet implemented.",
		},
	})
}

// HealthCheck returns service health.
func (h *RealtimeHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"service": "realtime-service",
			"status":  "healthy",
		},
	})
}
