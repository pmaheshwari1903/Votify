package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/votify/poll-service/internal/config"
	"github.com/votify/poll-service/internal/routes"
)

func main() {
	cfg := config.Load()

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	routes.Register(router, cfg)

	addr := ":" + cfg.Port
	log.Printf("[Poll Service] Starting on %s (env=%s)", addr, cfg.Env)

	if err := router.Run(addr); err != nil {
		log.Fatalf("[Poll Service] Failed to start: %v", err)
		os.Exit(1)
	}
}
