package handler

import (
	"booking-app/internal/entity"
	"booking-app/internal/repository"
	"net/http"
	"github.com/gin-gonic/gin"
	"time"
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
	ctx := c.Request.Context()
	user := entity.User{
		Name:     "Haryo",
		Email:    "haryo@test.com",
		Password: "123456",
	}

	h.userRepo.Create(ctx, &user)

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
	ctx := c.Request.Context()
	userID :=
		c.MustGet(
			"userID",
		).(uint)

	user, err :=
		h.userRepo.FindByID(
			ctx,
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

func (h *UserHandler) Slow(
	c *gin.Context,
) {
	select {
	case <-time.After(10 * time.Second):
		c.JSON(200, gin.H{
			"message": "done",
		})

	case <-c.Request.Context().Done():
		c.JSON(408, gin.H{
			"message": "timeout",
		})
	}
}