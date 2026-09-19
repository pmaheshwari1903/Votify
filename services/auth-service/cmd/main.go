package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/votify/auth-service/internal/config"
	"github.com/votify/auth-service/internal/routes"
)

func main() {
	cfg := config.Load()

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	routes.Register(router, cfg)

	addr := ":" + cfg.Port
	log.Printf("[Auth Service] Starting on %s (env=%s)", addr, cfg.Env)

	if err := router.Run(addr); err != nil {
		log.Fatalf("[Auth Service] Failed to start: %v", err)
		os.Exit(1)
	}
}
