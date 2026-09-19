package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/votify/api-gateway/internal/config"
	"github.com/votify/api-gateway/internal/handler"
	"github.com/votify/api-gateway/internal/middleware"
	"github.com/votify/pkg/health"
)

// Register sets up all API Gateway routes.
// The gateway is the single frontend entry point that forwards requests
// to specialized internal microservices. It contains NO business logic or DB access.
func Register(router *gin.Engine, cfg *config.Config) {
	// Global middleware
	router.Use(middleware.CORS())
	router.Use(middleware.RequestID())

	// Operational health & readiness endpoints
	checker := health.NewChecker()
	router.GET("/health", health.HealthCheckHandler("api-gateway"))
	router.GET("/ready", health.ReadinessCheckHandler("api-gateway", checker))

	// API v1 route groups
	v1 := router.Group("/api/v1")
	{
		// Auth routes → Auth Service (:8081)
		v1.Any("/auth/*path", handler.ProxyPlaceholder("auth-service", cfg.AuthServiceURL))

		// Poll routes → Poll Service (:8082)
		v1.Any("/polls/*path", handler.ProxyPlaceholder("poll-service", cfg.PollServiceURL))

		// Vote routes → Vote Service (:8083)
		v1.Any("/votes/*path", handler.ProxyPlaceholder("vote-service", cfg.VoteServiceURL))

		// Analytics routes → Analytics Service (:8085)
		v1.Any("/analytics/*path", handler.ProxyPlaceholder("analytics-service", cfg.AnalyticsServiceURL))

		// Payment routes → Payment Service (:8086)
		v1.Any("/payments/*path", handler.ProxyPlaceholder("payment-service", cfg.PaymentServiceURL))
	}

	// Realtime WebSocket proxy boundary → Realtime Service (:8084)
	router.Any("/ws/*path", handler.ProxyPlaceholder("realtime-service", cfg.RealtimeServiceURL))

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "The requested API endpoint does not exist.",
			},
		})
	})
}
