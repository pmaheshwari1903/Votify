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
func Register(router *gin.Engine, cfg *config.Config) {
	// Global middleware
	router.Use(middleware.CORS())
	router.Use(middleware.RequestID())

	// Operational health & readiness endpoints
	checker := health.NewChecker()
	router.GET("/health", health.HealthCheckHandler("api-gateway"))
	router.GET("/ready", health.ReadinessCheckHandler("api-gateway", checker))

	// Optional auth extraction middleware for all API routes
	optAuth := middleware.OptionalAuth(cfg.JWTSecret)
	reqAuth := middleware.AuthGuard(cfg.JWTSecret)

	// API v1 route groups
	v1 := router.Group("/api/v1")
	{
		// Auth routes → Auth Service (:8081)
		v1.Any("/auth/*path", optAuth, handler.ReverseProxy(cfg.AuthServiceURL, "/api/v1"))

		// Public Poll endpoints → Poll Service (:8082)
		v1.GET("/polls/public/:id", handler.ReverseProxy(cfg.PollServiceURL, "/api/v1"))
		v1.GET("/polls/:id/results", handler.ReverseProxy(cfg.PollServiceURL, "/api/v1"))

		// Protected Poll endpoints → Poll Service (:8082)
		v1.POST("/polls", reqAuth, handler.ReverseProxy(cfg.PollServiceURL, "/api/v1"))
		v1.GET("/polls", reqAuth, handler.ReverseProxy(cfg.PollServiceURL, "/api/v1"))
		v1.GET("/polls/:id", reqAuth, handler.ReverseProxy(cfg.PollServiceURL, "/api/v1"))
		v1.PATCH("/polls/:id", reqAuth, handler.ReverseProxy(cfg.PollServiceURL, "/api/v1"))
		v1.DELETE("/polls/:id", reqAuth, handler.ReverseProxy(cfg.PollServiceURL, "/api/v1"))
		v1.POST("/polls/:id/open", reqAuth, handler.ReverseProxy(cfg.PollServiceURL, "/api/v1"))
		v1.POST("/polls/:id/close", reqAuth, handler.ReverseProxy(cfg.PollServiceURL, "/api/v1"))

		// Vote routes → Vote Service (:8083)
		v1.POST("/votes", optAuth, handler.ReverseProxy(cfg.VoteServiceURL, "/api/v1"))

		// Analytics routes → Analytics Service (:8085)
		v1.Any("/analytics/*path", optAuth, handler.ReverseProxy(cfg.AnalyticsServiceURL, "/api/v1"))

		// Payment routes → Payment Service (:8086)
		v1.Any("/payments/*path", reqAuth, handler.ReverseProxy(cfg.PaymentServiceURL, "/api/v1"))
	}

	// Realtime WebSocket proxy boundary → Realtime Service (:8084)
	router.Any("/ws/*path", handler.ReverseProxy(cfg.RealtimeServiceURL, ""))

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
