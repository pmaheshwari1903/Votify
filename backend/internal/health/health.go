package health

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ComponentHealth struct {
	Status  string `json:"status"` // "up" or "down"
	Details string `json:"details,omitempty"`
}

type HealthResponse struct {
	Service    string                     `json:"service"`
	Status     string                     `json:"status"` // "up", "degraded", "down"
	Timestamp  string                     `json:"timestamp"`
	Components map[string]ComponentHealth `json:"components,omitempty"`
}

func HealthCheckHandler(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, HealthResponse{
			Service:   serviceName,
			Status:    "up",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
	}
}

func ReadinessCheckHandler(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		components := map[string]ComponentHealth{
			"mongo": {Status: "up"},
			"redis": {Status: "up"},
		}
		c.JSON(http.StatusOK, HealthResponse{
			Service:    serviceName,
			Status:     "up",
			Timestamp:  time.Now().UTC().Format(time.RFC3339),
			Components: components,
		})
	}
}
