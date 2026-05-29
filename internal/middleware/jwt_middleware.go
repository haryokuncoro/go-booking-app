package middleware

import (
	"booking-app/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware(
	secret string,
) gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader :=
			c.GetHeader(
				"Authorization",
			)

		if authHeader == "" {

			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "missing token",
				},
			)

			return
		}

		tokenString :=
			strings.TrimPrefix(
				authHeader,
				"Bearer ",
			)

		claims, err :=
			utils.ParseToken(
				tokenString,
				secret,
			)

		if err != nil {

			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{
					"error": "invalid token",
				},
			)

			return
		}

		c.Set(
			"userID",
			claims.UserID,
		)

		c.Next()

	}
}
