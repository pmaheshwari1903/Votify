package config

import "os"

// Config holds the API Gateway configuration.
// All values are loaded from environment variables.
type Config struct {
	Env  string
	Port string

	// JWT configuration
	JWTSecret string

	// Downstream service URLs
	AuthServiceURL      string
	PollServiceURL      string
	VoteServiceURL      string
	RealtimeServiceURL  string
	AnalyticsServiceURL string
	PaymentServiceURL   string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		Env:                 getEnv("APP_ENV", "development"),
		Port:                getEnv("PORT", getEnv("API_GATEWAY_PORT", "8080")),
		JWTSecret:           getEnv("JWT_SECRET", "your-jwt-secret-min-32-chars"),
		AuthServiceURL:      getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
		PollServiceURL:      getEnv("POLL_SERVICE_URL", "http://localhost:8082"),
		VoteServiceURL:      getEnv("VOTE_SERVICE_URL", "http://localhost:8083"),
		RealtimeServiceURL:  getEnv("REALTIME_SERVICE_URL", "http://localhost:8084"),
		AnalyticsServiceURL: getEnv("ANALYTICS_SERVICE_URL", "http://localhost:8085"),
		PaymentServiceURL:   getEnv("PAYMENT_SERVICE_URL", "http://localhost:8086"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
