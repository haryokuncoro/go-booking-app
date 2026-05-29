package handler

import (
	"booking-app/internal/entity"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	DB *gorm.DB
}

func NewUserHandler(
	db *gorm.DB,
) *UserHandler {
	return &UserHandler{
		DB: db,
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

	h.DB.Create(&user)

	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "user inserted",
		},
	)
}
