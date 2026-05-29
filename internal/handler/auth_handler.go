package handler

import (
	"booking-app/internal/dto"
	"booking-app/internal/service"
	"net/http"

	"booking-app/internal/response"
	customValidator "booking-app/internal/validator"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(
	authService service.AuthService,
) *AuthHandler {

	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(
	c *gin.Context,
) {

	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(
		&req,
	); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	err :=
		customValidator.Validate.
			Struct(req)

	if err != nil {

		response.Error(
			c,
			http.StatusBadRequest,
			customValidator.
				FormatValidationError(
					err,
				),
		)

		return
	}

	err = h.authService.Register(
		req,
	)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	response.Success(
		c,
		http.StatusCreated,
		"user registered",
		nil,
	)
}

func (h *AuthHandler) Login(
	c *gin.Context,
) {

	var req dto.LoginRequest

	if err := c.ShouldBindJSON(
		&req,
	); err != nil {

		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	err :=
		customValidator.Validate.
			Struct(req)
	if err != nil {

		response.Error(
			c,
			http.StatusBadRequest,
			customValidator.
				FormatValidationError(
					err,
				),
		)

		return
	}

	token, err :=
		h.authService.Login(req)

	if err != nil {

		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	response.Success(
		c,
		http.StatusOK,
		"login success",
		gin.H{
			"access_token": token,
		},
	)
}
