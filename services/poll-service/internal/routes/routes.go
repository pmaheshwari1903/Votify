package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/votify/pkg/health"
	"github.com/votify/poll-service/internal/config"
	"github.com/votify/poll-service/internal/handler"
)

// Register sets up all Poll Service routes.
func Register(router *gin.Engine, cfg *config.Config) {
	h := handler.NewPollHandler()

	// Operational health & readiness endpoints
	checker := health.NewChecker()
	router.GET("/health", health.HealthCheckHandler("poll-service"))
	router.GET("/ready", health.ReadinessCheckHandler("poll-service", checker))

	// Poll endpoints boundary
	polls := router.Group("/polls")
	{
		polls.POST("", h.Create)
		polls.GET("", h.List)
		polls.GET("/:id", h.GetByID)
		polls.PUT("/:id", h.Update)
		polls.DELETE("/:id", h.Delete)
	}
}
