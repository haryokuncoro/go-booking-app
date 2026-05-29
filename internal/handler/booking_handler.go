package handler

import (
	"booking-app/internal/dto"
	"booking-app/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type BookingHandler struct {
	bookingService service.BookingService
}

func NewBookingHandler(
	bookingService service.BookingService,
) *BookingHandler {

	return &BookingHandler{
		bookingService: bookingService,
	}
}

func (h *BookingHandler) CreateBooking(
	c *gin.Context,
) {

	var req dto.CreateBookingRequest

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

	userID :=
		c.MustGet(
			"userID",
		).(uint)

	err :=
		h.bookingService.
			CreateBooking(
				userID,
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

	c.JSON(
		http.StatusCreated,
		gin.H{
			"message": "booking created",
		},
	)
}

func (h *BookingHandler) ListBookings(
	c *gin.Context,
) {

	userID :=
		c.MustGet(
			"userID",
		).(uint)

	bookings, err :=
		h.bookingService.
			GetUserBookings(
				userID,
			)

	if err != nil {

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		bookings,
	)
}

func (h *BookingHandler) GetBooking(
	c *gin.Context,
) {

	id64, _ :=
		strconv.ParseUint(
			c.Param("id"),
			10,
			64,
		)

	booking, err :=
		h.bookingService.
			GetBooking(
				uint(id64),
			)

	if err != nil {

		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": "booking not found",
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		booking,
	)
}

