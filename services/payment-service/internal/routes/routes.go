package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/votify/payment-service/internal/config"
	"github.com/votify/payment-service/internal/handler"
	"github.com/votify/pkg/health"
)

// Register sets up all Payment Service routes.
func Register(router *gin.Engine, cfg *config.Config) {
	h := handler.NewPaymentHandler()

	// Operational health & readiness endpoints
	checker := health.NewChecker()
	router.GET("/health", health.HealthCheckHandler("payment-service"))
	router.GET("/ready", health.ReadinessCheckHandler("payment-service", checker))

	// Payment & subscription endpoints boundary
	payments := router.Group("/payments")
	{
		payments.GET("/plans", h.GetPlans)
		payments.POST("/checkout", h.CreateCheckout)
		payments.POST("/webhook", h.Webhook)
	}
}
