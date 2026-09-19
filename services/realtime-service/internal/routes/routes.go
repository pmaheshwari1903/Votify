package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/votify/pkg/health"
	"github.com/votify/realtime-service/internal/config"
	"github.com/votify/realtime-service/internal/handler"
)

// Register sets up all Realtime Service routes.
// This service has a minimal HTTP surface — its primary transport is WebSocket.
func Register(router *gin.Engine, cfg *config.Config) {
	h := handler.NewRealtimeHandler()

	// Operational health & readiness endpoints
	checker := health.NewChecker()
	router.GET("/health", health.HealthCheckHandler("realtime-service"))
	router.GET("/ready", health.ReadinessCheckHandler("realtime-service", checker))

	// WebSocket browser-facing transport boundary
	router.GET("/ws/:pollId", h.WebSocketUpgrade)
}
