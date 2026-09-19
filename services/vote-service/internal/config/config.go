package config

import "os"

// Config holds the Vote Service configuration.
type Config struct {
	Env  string
	Port string

	// MongoDB
	MongoURI string
	MongoDB  string

	// Kafka (for publishing VoteCreated events)
	KafkaBrokers string
}

// Load reads configuration from environment variables.
func Load() *Config {
	return &Config{
		Env:          getEnv("APP_ENV", "development"),
		Port:         getEnv("VOTE_SERVICE_PORT", "8083"),
		MongoURI:     getEnv("VOTE_MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:      getEnv("VOTE_MONGO_DB", "votify_votes"),
		KafkaBrokers: getEnv("KAFKA_BROKERS", "localhost:9092"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
