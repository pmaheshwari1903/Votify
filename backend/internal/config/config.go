package config

import (
	"os"
	"strconv"
)

// Config holds all environment settings for the single Votify backend.
type Config struct {
	Env       string
	Port      string
	JWTSecret string

	// MongoDB settings
	MongoURI    string
	AuthMongoDB string
	PollMongoDB string
	VoteMongoDB string

	// Redis settings
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// OAuth 2.0 / OIDC settings
	OAuthClientID     string
	OAuthClientSecret string
	OAuthIssuer       string
	OAuthRedirectURI  string
	FrontendURL       string
}

// Load loads configuration from environment variables with sensible defaults.
func Load() *Config {
	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		redisDB = 0
	}

	return &Config{
		Env:               getEnv("APP_ENV", "development"),
		Port:              getEnv("PORT", getEnv("API_GATEWAY_PORT", "8080")),
		JWTSecret:         getEnv("JWT_SECRET", "your-jwt-secret-min-32-chars"),
		MongoURI:          getEnv("MONGO_URI", ""),
		AuthMongoDB:       getEnv("AUTH_MONGO_DB", "votify_auth"),
		PollMongoDB:       getEnv("POLL_MONGO_DB", "votify_polls"),
		VoteMongoDB:       getEnv("VOTE_MONGO_DB", "votify_votes"),
		RedisAddr:         getEnv("REDIS_ADDR", ""),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		RedisDB:           redisDB,
		OAuthClientID:     getEnv("OAUTH_CLIENT_ID", ""),
		OAuthClientSecret: getEnv("OAUTH_CLIENT_SECRET", ""),
		OAuthIssuer:       getEnv("OAUTH_ISSUER", "https://oidcauth.vercel.app"),
		OAuthRedirectURI:  getEnv("OAUTH_REDIRECT_URI", "https://votify-production-2915.up.railway.app/api/v1/auth/oauth/callback"),
		FrontendURL:       getEnv("FRONTEND_URL", "https://valiant-eagerness-production-ca89.up.railway.app"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
