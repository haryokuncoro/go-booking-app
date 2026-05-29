package main

import (
	"booking-app/config"
	"booking-app/internal/handler"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()

	r := gin.Default()

	healthHandler := handler.NewHealthHandler()

	r.GET(
		"/health",
		healthHandler.Health,
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