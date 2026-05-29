package main

import (
	"booking-app/config"
	"booking-app/internal/database"
	"booking-app/internal/handler"
	"booking-app/internal/middleware"
	"booking-app/internal/repository"
	"booking-app/internal/service"
	"booking-app/internal/worker"
	"booking-app/internal/cache"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()

	db := database.ConnectDB(cfg)

	_ = db

	r := gin.Default()

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

	fmt.Printf(
		"Starting %s on port %s\n",
		cfg.AppName,
		cfg.AppPort,
	)

	err := r.Run(
		":" + cfg.AppPort,
	)

	if err != nil {
		panic(err)
	}
}
