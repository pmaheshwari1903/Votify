package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/votify/backend/internal/auth"
	"github.com/votify/backend/internal/config"
	"github.com/votify/backend/internal/health"
	"github.com/votify/backend/internal/middleware"
	"github.com/votify/backend/internal/polls"
	"github.com/votify/backend/internal/realtime"
	"github.com/votify/backend/internal/repository"
	"github.com/votify/backend/internal/votes"
	"github.com/votify/backend/internal/websocket"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.Load()

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Repositories (MongoDB with in-memory fallback)
	var userRepo repository.UserRepository
	var pollRepo repository.PollRepository
	var voteRepo repository.VoteRepository

	if cfg.MongoURI != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
		if err == nil {
			err = client.Ping(ctx, nil)
		}
		cancel()

		if err == nil {
			log.Println("[Database] Connected and pinged MongoDB successfully")
			authDB := client.Database(cfg.AuthMongoDB)
			pollDB := client.Database(cfg.PollMongoDB)
			voteDB := client.Database(cfg.VoteMongoDB)

			userRepo = repository.NewMongoUserRepository(authDB)
			pollRepo = repository.NewMongoPollRepository(pollDB)
			voteRepo = repository.NewMongoVoteRepository(voteDB)
		} else {
			log.Printf("[Database] Failed to connect/ping MongoDB (%v), falling back to in-memory store", err)
		}
	}

	if userRepo == nil {
		log.Println("[Database] Using in-memory repositories (Development/Offline mode)")
		userRepo = repository.NewMemoryUserRepository()
		pollRepo = repository.NewMemoryPollRepository()
		voteRepo = repository.NewMemoryVoteRepository()
	}

	// Initialize Core Services
	authService := auth.NewAuthService(userRepo, cfg.JWTSecret)
	pollService := polls.NewPollService(pollRepo)
	realtimeService := realtime.NewRealtimeService(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB, pollService)
	voteService := votes.NewVoteService(voteRepo, pollService, realtimeService)

	// Initialize Handlers & Managers
	authHandler := auth.NewAuthHandler(authService)
	pollHandler := polls.NewPollHandler(pollService)
	voteHandler := votes.NewVoteHandler(voteService)
	wsManager := websocket.NewConnectionManager(realtimeService)

	// Setup Gin Engine
	router := gin.Default()
	router.Use(middleware.CORS())
	router.Use(middleware.RequestID())

	// Health check endpoints
	router.GET("/health", health.HealthCheckHandler("votify-backend"))
	router.GET("/ready", health.ReadinessCheckHandler("votify-backend"))

	// Middleware guards
	reqAuth := middleware.AuthGuard(cfg.JWTSecret)

	// API v1 Routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.GET("/me", reqAuth, authHandler.Me)
			authGroup.POST("/logout", authHandler.Logout)
			authGroup.GET("/oauth/login", authHandler.OAuthLogin)
			authGroup.GET("/oauth/callback", authHandler.OAuthCallback)
		}

		// Public Poll routes (placed before /polls/:id to prevent wildcard collision)
		v1.GET("/polls/public/:id", pollHandler.GetPublic)
		v1.GET("/polls/:id/results", pollHandler.Results)

		// Protected Poll routes
		v1.POST("/polls", reqAuth, pollHandler.Create)
		v1.GET("/polls", reqAuth, pollHandler.List)
		v1.GET("/polls/:id", reqAuth, pollHandler.GetByID)
		v1.PATCH("/polls/:id", reqAuth, pollHandler.Update)
		v1.DELETE("/polls/:id", reqAuth, pollHandler.Delete)
		v1.POST("/polls/:id/open", reqAuth, pollHandler.Open)
		v1.POST("/polls/:id/close", reqAuth, pollHandler.Close)

		// Vote routes (only authenticated users can vote)
		v1.POST("/votes", reqAuth, voteHandler.CastVote)
	}

	// Realtime WebSocket endpoint
	router.GET("/ws/:pollId", wsManager.HandleWS)

	// Fallback 404 handler
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "The requested API endpoint does not exist.",
			},
		})
	})

	addr := ":" + cfg.Port
	log.Printf("[Votify Backend] Starting server on %s (env=%s)", addr, cfg.Env)

	if err := router.Run(addr); err != nil {
		log.Fatalf("[Votify Backend] Failed to start server: %v", err)
		os.Exit(1)
	}
}
