package middleware

import (
	"time"

	"booking-app/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RequestLogger() gin.HandlerFunc {

	return func(c *gin.Context) {

		start := time.Now()

		c.Next()

		logger.Log.Info(
			"http request",

			zap.String(
				"method",
				c.Request.Method,
			),

			zap.String(
				"path",
				c.Request.URL.Path,
			),

			zap.Int(
				"status",
				c.Writer.Status(),
			),

			zap.Duration(
				"duration",
				time.Since(start),
			),
		)
	}
}