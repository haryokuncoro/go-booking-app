package main

import (
	_ "booking-app/docs"
	"booking-app/config"
	"booking-app/internal/cache"
	"booking-app/internal/database"
	"booking-app/internal/handler"
	"booking-app/internal/logger"
	"booking-app/internal/repository"
	"booking-app/internal/router"
	"booking-app/internal/service"
	"booking-app/internal/worker"
	"go.uber.org/zap"
)

// @title Booking API
// @version 1.0
// @description Production-style Booking API using Gin, PostgreSQL, Redis, JWT, Goroutines and Worker Pools.
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	logger.Init()
	defer logger.Sync()

	cfg := config.LoadConfig()
	db := database.ConnectDB(cfg)
	redisClient := cache.NewRedis(cfg)
	worker.StartWorkers()

	userRepo    := repository.NewUserRepository(db)
	bookingRepo := repository.NewBookingRepository(db)

	authService    := service.NewAuthService(userRepo, cfg)
	bookingService := service.NewBookingService(bookingRepo, userRepo, redisClient)

	r := router.Setup(
		cfg,
		handler.NewHealthHandler(),
		handler.NewAuthHandler(authService),
		handler.NewUserHandler(userRepo),
		handler.NewBookingHandler(bookingService),
	)

	logger.Log.Info("application started",
		zap.String("app", cfg.AppName),
		zap.String("env", cfg.AppEnv),
		zap.String("port", cfg.AppPort),
	)

	if err := r.Run(":" + cfg.AppPort); err != nil {
		panic(err)
	}
}