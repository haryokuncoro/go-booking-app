package handler

import (
	"booking-app/internal/dto"
	"booking-app/internal/response"
	"booking-app/internal/service"
	customValidator "booking-app/internal/validator"
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

	userID :=
		c.MustGet(
			"userID",
		).(uint)

	err =
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

	
	response.Success(
	c,
	http.StatusCreated,
	"booking created",
	nil,
)
}

func (h *BookingHandler) ListBookings(
	c *gin.Context,
) {
	ctx := c.Request.Context()
	userID :=
		c.MustGet(
			"userID",
		).(uint)

	bookings, err :=
		h.bookingService.
			GetUserBookings(
				ctx,
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
	ctx := c.Request.Context()
	booking, err :=
		h.bookingService.
			GetBooking(
				ctx,
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
