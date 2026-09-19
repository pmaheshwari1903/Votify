package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/votify/realtime-service/internal/pubsub"
	ws "github.com/votify/realtime-service/internal/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow CORS for WebSockets across local/cloud environments
		return true
	},
}

type RealtimeHandler struct {
	manager    *ws.ConnectionManager
	redisStore *pubsub.RedisRealtimeStore
}

func NewRealtimeHandler(manager *ws.ConnectionManager, redisStore *pubsub.RedisRealtimeStore) *RealtimeHandler {
	return &RealtimeHandler{
		manager:    manager,
		redisStore: redisStore,
	}
}

// WebSocketUpgrade handles GET /ws/:pollId
func (h *RealtimeHandler) WebSocketUpgrade(c *gin.Context) {
	pollID := strings.TrimSpace(c.Param("pollId"))
	if pollID == "" || pollID == ":pollId" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pollId is required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &ws.Client{
		PollID: pollID,
		Conn:   conn,
		Send:   make(chan pubsub.RealtimeUpdateMessage, 256),
	}

	// Send initial current-results snapshot immediately upon connection
	snapshot, err := h.redisStore.BuildSnapshot(context.Background(), pollID)
	if err == nil && snapshot != nil {
		_ = conn.WriteJSON(snapshot)
	}

	h.manager.Register(client)
}

// HealthCheck returns service health status.
func (h *RealtimeHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"service": "realtime-service",
			"status":  "healthy",
		},
	})
}
