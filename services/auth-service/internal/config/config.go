package config

import "os"

// Config holds the Auth Service configuration.
type Config struct {
	Env  string
	Port string

	// MongoDB
	MongoURI string
	MongoDB  string

	// JWT
	JWTSecret string
	JWTExpiry string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		Env:       getEnv("APP_ENV", "development"),
		Port:      getEnv("AUTH_SERVICE_PORT", "8081"),
		MongoURI:  getEnv("AUTH_MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:   getEnv("AUTH_MONGO_DB", "votify_auth"),
		JWTSecret: getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		JWTExpiry: getEnv("JWT_EXPIRY", "24h"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
