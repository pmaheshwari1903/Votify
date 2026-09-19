package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/votify/pkg/health"
	"github.com/votify/vote-service/internal/config"
	"github.com/votify/vote-service/internal/handler"
	"github.com/votify/vote-service/internal/repository"
	"github.com/votify/vote-service/internal/service"
)

// Register sets up all Vote Service routes.
func Register(router *gin.Engine, cfg *config.Config) {
	voteRepo := repository.NewMemoryVoteRepository()
	voteService := service.NewVoteService(voteRepo, cfg.PollServiceURL)
	h := handler.NewVoteHandler(voteService)

	// Operational health & readiness endpoints
	checker := health.NewChecker()
	router.GET("/health", health.HealthCheckHandler("vote-service"))
	router.GET("/ready", health.ReadinessCheckHandler("vote-service", checker))

	// Vote endpoints boundary
	votes := router.Group("/votes")
	{
		votes.POST("", h.CastVote)
	}
}
