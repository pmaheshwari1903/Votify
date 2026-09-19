package config

import "os"

// Config holds the Payment Service configuration.
type Config struct {
	Env  string
	Port string

	// MongoDB
	MongoURI string
	MongoDB  string

	// Kafka (for publishing payment/subscription events)
	KafkaBrokers string
}

// Load reads configuration from environment variables.
func Load() *Config {
	return &Config{
		Env:          getEnv("APP_ENV", "development"),
		Port:         getEnv("PAYMENT_SERVICE_PORT", "8086"),
		MongoURI:     getEnv("PAYMENT_MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:      getEnv("PAYMENT_MONGO_DB", "votify_payments"),
		KafkaBrokers: getEnv("KAFKA_BROKERS", "localhost:9092"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
