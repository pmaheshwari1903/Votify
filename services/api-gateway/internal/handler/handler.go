package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck returns the gateway's health status.
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"service": "api-gateway",
			"status":  "healthy",
		},
	})
}

// ProxyPlaceholder returns a handler that acknowledges the route exists
// but the actual proxy logic is not yet implemented.
//
// In the next development phase, this will be replaced with actual
// HTTP reverse proxy logic (using net/http/httputil.ReverseProxy).
func ProxyPlaceholder(serviceName, serviceURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "SERVICE_NOT_IMPLEMENTED",
				"message": "Proxy to " + serviceName + " is not yet implemented.",
				"target":  serviceURL,
			},
		})
	}
}
