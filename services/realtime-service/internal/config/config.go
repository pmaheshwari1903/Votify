package config

import "os"

// Config holds the Realtime Service configuration.
type Config struct {
	Env  string
	Port string

	// Redis (pub/sub for live vote distribution)
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// Kafka (consumes VoteCreated events from Vote Service)
	KafkaBrokers string
	KafkaGroupID string
}

// Load reads configuration from environment variables.
func Load() *Config {
	return &Config{
		Env:           getEnv("APP_ENV", "development"),
		Port:          getEnv("REALTIME_SERVICE_PORT", "8084"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       0,
		KafkaBrokers:  getEnv("KAFKA_BROKERS", "localhost:9092"),
		KafkaGroupID:  getEnv("KAFKA_GROUP_ID", "votify-realtime"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
