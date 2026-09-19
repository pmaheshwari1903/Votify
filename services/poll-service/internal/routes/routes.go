package routes

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/votify/pkg/health"
	"github.com/votify/poll-service/internal/config"
	"github.com/votify/poll-service/internal/handler"
	"github.com/votify/poll-service/internal/repository"
	"github.com/votify/poll-service/internal/service"
)

var (
	sharedRepo *repository.MemoryPollRepository
	once       sync.Once
)

// GetSharedPollRepository returns a singleton repository instance for cross-service testing/access
func GetSharedPollRepository() *repository.MemoryPollRepository {
	once.Do(func() {
		sharedRepo = repository.NewMemoryPollRepository()
	})
	return sharedRepo
}

// Register sets up all Poll Service routes.
func Register(router *gin.Engine, cfg *config.Config) {
	pollRepo := GetSharedPollRepository()
	pollService := service.NewPollService(pollRepo)
	h := handler.NewPollHandler(pollService)

	// Operational health & readiness endpoints
	checker := health.NewChecker()
	router.GET("/health", health.HealthCheckHandler("poll-service"))
	router.GET("/ready", health.ReadinessCheckHandler("poll-service", checker))

	// Public poll endpoint (unauthenticated)
	router.GET("/polls/public/:id", h.GetPublic)
	router.GET("/polls/:id/results", h.Results)

	// Poll endpoints boundary
	polls := router.Group("/polls")
	{
		polls.POST("", h.Create)
		polls.GET("", h.List)
		polls.GET("/:id", h.GetByID)
		polls.PATCH("/:id", h.Update)
		polls.DELETE("/:id", h.Delete)
		polls.POST("/:id/open", h.Open)
		polls.POST("/:id/close", h.Close)
	}

	// Internal service-to-service endpoint for vote increments
	router.POST("/internal/polls/:id/options/:optionId/vote", h.InternalIncrementVote)
}
