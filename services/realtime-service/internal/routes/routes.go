package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/votify/pkg/health"
	"github.com/votify/pkg/logging"
	"github.com/votify/realtime-service/internal/broker"
	"github.com/votify/realtime-service/internal/config"
	"github.com/votify/realtime-service/internal/handler"
	"github.com/votify/realtime-service/internal/pubsub"
	"github.com/votify/realtime-service/internal/websocket"
)

// Register sets up all Realtime Service routes and initializes Kafka & Redis subscribers.
func Register(router *gin.Engine, cfg *config.Config) {
	logger := logging.NewLogger("realtime-service", cfg.Env, "info")

	redisStore := pubsub.NewRedisRealtimeStore(
		cfg.RedisAddr,
		cfg.RedisPassword,
		cfg.PollServiceURL,
	)

	_ = broker.NewKafkaConsumer(
		"votify-realtime-service",
		redisStore,
		logger,
	)

	connManager := websocket.NewConnectionManager(redisStore, logger)

	h := handler.NewRealtimeHandler(connManager, redisStore)

	// Operational health & readiness endpoints
	checker := health.NewChecker()
	router.GET("/health", health.HealthCheckHandler("realtime-service"))
	router.GET("/ready", health.ReadinessCheckHandler("realtime-service", checker))

	// WebSocket browser-facing transport boundary
	router.GET("/ws/:pollId", h.WebSocketUpgrade)
}