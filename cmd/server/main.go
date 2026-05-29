package main

import (
	"booking-app/config"
	"booking-app/internal/database"
	"booking-app/internal/handler"
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
	userHandler := handler.NewUserHandler(db)
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
