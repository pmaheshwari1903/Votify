package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/votify/api-gateway/internal/config"
	"github.com/votify/api-gateway/internal/routes"
)

func main() {
	cfg := config.Load()

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	routes.Register(router, cfg)

	addr := ":" + cfg.Port
	log.Printf("[API Gateway] Starting on %s (env=%s)", addr, cfg.Env)

	if err := router.Run(addr); err != nil {
		log.Fatalf("[API Gateway] Failed to start: %v", err)
		os.Exit(1)
	}
}
