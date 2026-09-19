package health

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// PingFunc is a health check function for a dependency (e.g. MongoDB, Redis, Kafka).
type PingFunc func(ctx context.Context) error

// Checker manages readiness dependency checks.
type Checker struct {
	checkers map[string]PingFunc
}

// NewChecker creates a new Checker instance.
func NewChecker() *Checker {
	return &Checker{
		checkers: make(map[string]PingFunc),
	}
}

// RegisterDependency adds a named ping function to readiness checks.
func (c *Checker) RegisterDependency(name string, ping PingFunc) {
	c.checkers[name] = ping
}

// HealthCheckHandler returns a 200 OK indicating the service process is alive.
func HealthCheckHandler(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "UP",
			"service":   serviceName,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	}
}

// ReadinessCheckHandler returns a 200 OK if all registered dependencies respond within timeout.
func ReadinessCheckHandler(serviceName string, checker *Checker) gin.HandlerFunc {
	return func(c *gin.Context) {
		if checker == nil || len(checker.checkers) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"status":    "READY",
				"service":   serviceName,
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		results := make(map[string]string)
		allHealthy := true

		for name, ping := range checker.checkers {
			if err := ping(ctx); err != nil {
				results[name] = "UNHEALTHY: " + err.Error()
				allHealthy = false
			} else {
				results[name] = "HEALTHY"
			}
		}

		status := http.StatusOK
		state := "READY"
		if !allHealthy {
			status = http.StatusServiceUnavailable
			state = "NOT_READY"
		}

		c.JSON(status, gin.H{
			"status":       state,
			"service":      serviceName,
			"dependencies": results,
			"timestamp":    time.Now().UTC().Format(time.RFC3339),
		})
	}
}
