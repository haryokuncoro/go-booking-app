package main

import (
	"booking-app/config"
	"booking-app/internal/cache"
	"booking-app/internal/database"
	"booking-app/internal/handler"
	"booking-app/internal/logger"
	"booking-app/internal/middleware"
	"booking-app/internal/repository"
	"booking-app/internal/service"
	"booking-app/internal/worker"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
	_ "booking-app/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

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

	_ = db

	r := gin.New()

	r.GET(
	"/swagger/*any",
	ginSwagger.WrapHandler(
		swaggerFiles.Handler,
	),
)

	r.Use(
		middleware.RequestLogger(),
	)

	r.Use(
		middleware.TimeoutMiddleware(
			5 * time.Second,
		),
	)

	r.Use(
		gin.Recovery(),
	)

	worker.StartWorkers()

	redisClient :=
		cache.NewRedis(
			cfg,
		)

	healthHandler := handler.NewHealthHandler()
	userRepo :=
		repository.NewUserRepository(
			db,
		)
	bookingRepo :=
		repository.NewBookingRepository(
			db,
		)
	authService :=
		service.NewAuthService(
			userRepo, cfg,
		)
	authHandler :=
		handler.NewAuthHandler(
			authService,
		)
	userHandler := handler.NewUserHandler(userRepo)
	bookingService :=
		service.NewBookingService(
			bookingRepo, userRepo, redisClient,
		)
	bookingHandler :=
		handler.NewBookingHandler(
			bookingService,
		)

	r.GET(
		"/health",
		healthHandler.Health,
	)

	r.POST(
		"/seed-user",
		userHandler.SeedUser,
	)

	

	api := r.Group("/api/v1")

	api.GET(
		"/slow",
		userHandler.Slow,
	)

	auth := api.Group("/auth")

	auth.POST(
		"/register",
		authHandler.Register,
	)

	auth.POST(
		"/login",
		authHandler.Login,
	)

	protected := api.Group("")

	protected.Use(
		middleware.JWTMiddleware(
			cfg.JWTSecret,
		),
	)

	protected.GET(
		"/me",
		userHandler.Me,
	)

	

	booking := protected.Group(
		"/bookings",
	)

	booking.POST(
		"",
		bookingHandler.CreateBooking,
	)

	booking.GET(
		"",
		bookingHandler.ListBookings,
	)

	booking.GET(
		"/:id",
		bookingHandler.GetBooking,
	)

	logger.Log.Info(
		"application started",

		zap.String(
			"app",
			cfg.AppName,
		),

		zap.String(
			"env",
			cfg.AppEnv,
		),

		zap.String(
			"port",
			cfg.AppPort,
		),
	)

	err := r.Run(
		":" + cfg.AppPort,
	)

	if err != nil {
		panic(err)
	}
}
