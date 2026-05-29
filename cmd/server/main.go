package main

import (
	"booking-app/config"
	"booking-app/internal/database"
	"booking-app/internal/handler"
	"booking-app/internal/middleware"
	"booking-app/internal/repository"
	"booking-app/internal/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()

	db := database.ConnectDB(cfg)

	_ = db

	r := gin.Default()

	healthHandler := handler.NewHealthHandler()
	userRepo :=
		repository.NewUserRepository(
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
