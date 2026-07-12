package router

import (
	"booking-app/internal/handler"
	"booking-app/internal/middleware"
	"booking-app/config"
	"time"

	_ "booking-app/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/gin-gonic/gin"
)

func Setup(
	cfg *config.Config,
	healthHandler *handler.HealthHandler,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	bookingHandler *handler.BookingHandler,
) *gin.Engine {
	r := gin.New()
	r.Use(
		middleware.RequestLogger(),
		middleware.TimeoutMiddleware(5*time.Second),
		gin.Recovery(),
	)

	r.GET("/swagger-ui/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health", healthHandler.Health)
	r.POST("/seed-user", userHandler.SeedUser)

	api := r.Group("/api/v1")
	api.GET("/slow", userHandler.Slow)

	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	protected := api.Group("")
	protected.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	protected.GET("/me", userHandler.Me)

	bookings := protected.Group("/bookings")
	bookings.POST("", bookingHandler.CreateBooking)
	bookings.GET("", bookingHandler.ListBookings)
	bookings.GET("/:id", bookingHandler.GetBooking)

	return r
}