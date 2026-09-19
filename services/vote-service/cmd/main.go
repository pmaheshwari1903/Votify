package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/votify/vote-service/internal/config"
	"github.com/votify/vote-service/internal/routes"
)

func main() {
	cfg := config.Load()

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	routes.Register(router, cfg)

	addr := ":" + cfg.Port
	log.Printf("[Vote Service] Starting on %s (env=%s)", addr, cfg.Env)

	if err := router.Run(addr); err != nil {
		log.Fatalf("[Vote Service] Failed to start: %v", err)
		os.Exit(1)
	}
}
