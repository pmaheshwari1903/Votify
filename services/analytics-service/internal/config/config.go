package config

import "os"

// Config holds the Analytics Service configuration.
type Config struct {
	Env  string
	Port string

	// MongoDB (analytics read-model database)
	MongoURI string
	MongoDB  string

	// Kafka (consumes domain events to build analytics)
	KafkaBrokers string
	KafkaGroupID string
}

// Load reads configuration from environment variables.
func Load() *Config {
	return &Config{
		Env:          getEnv("APP_ENV", "development"),
		Port:         getEnv("ANALYTICS_SERVICE_PORT", "8085"),
		MongoURI:     getEnv("ANALYTICS_MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:      getEnv("ANALYTICS_MONGO_DB", "votify_analytics"),
		KafkaBrokers: getEnv("KAFKA_BROKERS", "localhost:9092"),
		KafkaGroupID: getEnv("KAFKA_GROUP_ID", "votify-analytics"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
