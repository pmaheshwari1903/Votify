package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/votify/pkg/health"
	"github.com/votify/vote-service/internal/config"
	"github.com/votify/vote-service/internal/handler"
)

// Register sets up all Vote Service routes.
func Register(router *gin.Engine, cfg *config.Config) {
	h := handler.NewVoteHandler()

	// Operational health & readiness endpoints
	checker := health.NewChecker()
	router.GET("/health", health.HealthCheckHandler("vote-service"))
	router.GET("/ready", health.ReadinessCheckHandler("vote-service", checker))

	// Vote endpoints boundary
	votes := router.Group("/votes")
	{
		votes.POST("", h.CastVote)
		votes.GET("/tally/:pollId", h.GetTally)
	}
}
