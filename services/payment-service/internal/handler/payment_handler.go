package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PaymentHandler handles HTTP requests for payment endpoints.
type PaymentHandler struct{}

// NewPaymentHandler creates a new PaymentHandler.
func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{}
}

// GetPlans handles GET /plans
func (h *PaymentHandler) GetPlans(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "Pricing plans are not yet implemented.",
		},
	})
}

// CreateCheckout handles POST /checkout
func (h *PaymentHandler) CreateCheckout(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "Checkout is not yet implemented.",
		},
	})
}

// Webhook handles POST /webhook
// Payment provider webhook endpoint.
func (h *PaymentHandler) Webhook(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "NOT_IMPLEMENTED",
			"message": "Payment webhooks are not yet implemented.",
		},
	})
}

// HealthCheck returns service health.
func (h *PaymentHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"service": "payment-service",
			"status":  "healthy",
		},
	})
}
