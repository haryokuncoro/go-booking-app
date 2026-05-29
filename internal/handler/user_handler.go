package handler

import (
	"booking-app/internal/entity"
	"booking-app/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

func NewUserHandler(
	userRepo repository.UserRepository,
) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

func (h *UserHandler) SeedUser(
	c *gin.Context,
) {

	user := entity.User{
		Name:     "Haryo",
		Email:    "haryo@test.com",
		Password: "123456",
	}

	h.userRepo.Create(&user)

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "user inserted",
		},
	)
}

func (h *UserHandler) Me(
	c *gin.Context,
) {

	userID :=
		c.MustGet(
			"userID",
		).(uint)

	user, err :=
		h.userRepo.FindByID(
			userID,
		)

	if err != nil {

		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "user not found",
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		user,
	)
}
