package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/votify/auth-service/internal/config"
	"github.com/votify/auth-service/internal/handler"
	"github.com/votify/auth-service/internal/repository"
	"github.com/votify/auth-service/internal/service"
	"github.com/votify/pkg/health"
)

// Register sets up all Auth Service routes and injects dependencies.
func Register(router *gin.Engine, cfg *config.Config) {
	userRepo := repository.NewMemoryUserRepository()
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	h := handler.NewAuthHandler(authService, cfg.JWTSecret)

	// Operational health & readiness endpoints
	checker := health.NewChecker()
	router.GET("/health", health.HealthCheckHandler("auth-service"))
	router.GET("/ready", health.ReadinessCheckHandler("auth-service", checker))

	// Auth endpoints boundary
	auth := router.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/logout", h.Logout)
		auth.GET("/me", h.Me)
	}
}
