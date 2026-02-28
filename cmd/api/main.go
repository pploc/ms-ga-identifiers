package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gym-api/ms-ga-identifier/internal/api/handler"
	"github.com/gym-api/ms-ga-identifier/internal/api/router"
	"github.com/gym-api/ms-ga-identifier/internal/infrastructure/external"
	"github.com/gym-api/ms-ga-identifier/internal/infrastructure/messaging"
	"github.com/gym-api/ms-ga-identifier/internal/infrastructure/persistence/gorm/repository"
	"github.com/gym-api/ms-ga-identifier/internal/middleware"
	"github.com/gym-api/ms-ga-identifier/internal/service"
	"github.com/gym-api/ms-ga-identifier/pkg/config"
	"github.com/gym-api/ms-ga-identifier/pkg/database"
	"github.com/gym-api/ms-ga-identifier/pkg/redis"
	"github.com/gym-api/ms-ga-identifier/pkg/utils"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	if err := utils.InitLogger(cfg.Server.Env); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer utils.SyncLogger()

	utils.Info("Starting ms-ga-identifier service", utils.String("environment", cfg.Server.Env))

	// Initialize database
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		utils.Fatal("Failed to connect to database", utils.ErrorField(err.Error()))
	}

	// Initialize Redis (optional, for caching and rate limiting)
	redisClient, err := redis.NewRedisClient(&cfg.Redis)
	if err != nil {
		utils.Warn("Failed to connect to Redis, continuing without it", utils.ErrorField(err.Error()))
		redisClient = nil
	}

	// Initialize repositories
	identityRepo := repository.NewIdentityRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)
	loginAttemptRepo := repository.NewLoginAttemptRepository(db)
	passwordResetRepo := repository.NewPasswordResetRepository(db)

	// Initialize external clients
	authClient := external.NewAuthClient(&cfg.Auth)

	// Initialize Kafka producer (optional)
	var kafkaProducer *messaging.KafkaProducer
	kafkaProducer = messaging.NewKafkaProducer(&cfg.Kafka)
	defer func() {
		if kafkaProducer != nil {
			kafkaProducer.Close()
		}
	}()

	// Initialize JWT utility
	jwtUtil := utils.NewJWTUtil(cfg.JWT.Secret, cfg.JWT.ExpirationTime)

	// Initialize services
	identityService := service.NewIdentityService(
		identityRepo,
		refreshTokenRepo,
		loginAttemptRepo,
		passwordResetRepo,
		authClient,
		kafkaProducer,
		jwtUtil,
		cfg,
	)

	tokenService := service.NewTokenService(
		identityRepo,
		refreshTokenRepo,
		authClient,
		jwtUtil,
		cfg,
	)

	passwordService := service.NewPasswordService(
		identityRepo,
		passwordResetRepo,
	)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtUtil)

	// Initialize handlers
	identityHandler := handler.NewIdentityHandler(identityService)
	tokenHandler := handler.NewTokenHandler(tokenService)
	passwordHandler := handler.NewPasswordHandler(passwordService)

	// Initialize router
	r := router.NewRouter(identityHandler, tokenHandler, passwordHandler, authMiddleware, cfg)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	// Start server in goroutine
	go func() {
		utils.Info("Starting HTTP server", utils.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Fatal("Failed to start HTTP server", utils.ErrorField(err.Error()))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	utils.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		utils.Fatal("Server forced to shutdown", utils.ErrorField(err.Error()))
	}

	// Close Redis connection
	if redisClient != nil {
		redisClient.Close()
	}

	utils.Info("Server exited")
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
