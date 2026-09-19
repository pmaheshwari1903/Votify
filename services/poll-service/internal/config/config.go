package config

import "os"

// Config holds the Poll Service configuration.
type Config struct {
	Env  string
	Port string

	// MongoDB
	MongoURI string
	MongoDB  string

	// Kafka (for publishing poll events)
	KafkaBrokers string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		Env:          getEnv("APP_ENV", "development"),
		Port:         getEnv("POLL_SERVICE_PORT", "8082"),
		MongoURI:     getEnv("POLL_MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:      getEnv("POLL_MONGO_DB", "votify_polls"),
		KafkaBrokers: getEnv("KAFKA_BROKERS", "localhost:9092"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
