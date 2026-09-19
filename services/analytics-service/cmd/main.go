package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/votify/analytics-service/internal/config"
	"github.com/votify/analytics-service/internal/routes"
)

func main() {
	cfg := config.Load()

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Future: Initialize Kafka consumer for VoteCreated, PollCreated events
	// Future: Initialize MongoDB connection for read-model storage

	router := gin.Default()
	routes.Register(router, cfg)

	addr := ":" + cfg.Port
	log.Printf("[Analytics Service] Starting on %s (env=%s)", addr, cfg.Env)

	if err := router.Run(addr); err != nil {
		log.Fatalf("[Analytics Service] Failed to start: %v", err)
		os.Exit(1)
	}
}
