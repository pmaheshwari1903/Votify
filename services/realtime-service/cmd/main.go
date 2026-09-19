package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/votify/realtime-service/internal/config"
	"github.com/votify/realtime-service/internal/routes"
)

func main() {
	cfg := config.Load()

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Future: Initialize Redis client
	// Future: Initialize Kafka consumer
	// Future: Start WebSocket hub

	router := gin.Default()
	routes.Register(router, cfg)

	addr := ":" + cfg.Port
	log.Printf("[Realtime Service] Starting on %s (env=%s)", addr, cfg.Env)

	if err := router.Run(addr); err != nil {
		log.Fatalf("[Realtime Service] Failed to start: %v", err)
		os.Exit(1)
	}
}
