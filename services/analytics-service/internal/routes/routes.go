package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/votify/analytics-service/internal/config"
	"github.com/votify/analytics-service/internal/handler"
	"github.com/votify/pkg/health"
)

// Register sets up all Analytics Service routes.
func Register(router *gin.Engine, cfg *config.Config) {
	h := handler.NewAnalyticsHandler()

	// Operational health & readiness endpoints
	checker := health.NewChecker()
	router.GET("/health", health.HealthCheckHandler("analytics-service"))
	router.GET("/ready", health.ReadinessCheckHandler("analytics-service", checker))

	// Analytics read-model endpoints boundary
	analytics := router.Group("/analytics")
	{
		analytics.GET("/polls/:pollId", h.GetPollAnalytics)
	}
}
