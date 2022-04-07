package booking

import (
	"context"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
	"gym/internal/constants"
	"gym/internal/errors"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	Validate          *validator.Validate
	BookingService    BookingService
	BookingRepository BookingRepository
}

func NewHandler(r *gin.Engine,
	route string,
	val *validator.Validate,
	BookingService BookingService,
	BookingRepo BookingRepository,
) {
	handler := &Handler{
		Validate:          val,
		BookingService:    BookingService,
		BookingRepository: BookingRepo,
	}
	v1 := r.Group("v1/" + route)
	{
		v1.GET("", handler.GetAll)
		v1.GET(":id", handler.GetByID)
		v1.GET("count", handler.GetTotalCount)
		v1.GET("date", handler.GetByDateRange)
		v1.GET("member/:id", handler.GetAllClassesByMemberId)
		v1.GET("class/:id", handler.GetAllMembersByClassId)
		v1.POST("", handler.Save)
		v1.PUT(":id", handler.Update)
		v1.DELETE(":id", handler.Delete)

	}
}

func (h *Handler) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.CTX_DEFAULT*time.Second)
	defer cancel()

	limit := c.Query("limit")
	offset := c.Query("offset")

	result, err := h.BookingRepository.GetAll(ctx, limit, offset)
	if err != nil {
		c.JSON(http.StatusNotFound, errors.Response{
			Status:  http.StatusNotFound,
			Type:    constants.ErrUnknownResource,
			Message: []string{err.Error()}})
	} else {
		c.JSON(http.StatusOK, result)
	}
}
func (h *Handler) GetByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.CTX_DEFAULT*time.Second)
	defer cancel()

	id := c.Param("id")
	if _, err := strconv.Atoi(id); err != nil {
		c.JSON(http.StatusNotFound, errors.Response{
			http.StatusNotFound,
			constants.ErrUnknownResource,
			[]string{constants.ErrWrongURLParamType}})
		return
	}
	result, err := h.BookingRepository.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, errors.Response{
			Status:  http.StatusNotFound,
			Type:    constants.ErrUnknownResource,
			Message: []string{err.Error()}})
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func (h *Handler) GetTotalCount(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.CTX_DEFAULT*time.Second)
	defer cancel()

	result, err := h.BookingRepository.GetTotalCount(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, errors.Response{
			Status:  http.StatusNotFound,
			Type:    constants.ErrUnknownResource,
			Message: []string{err.Error()}})
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func (h *Handler) GetByDateRange(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.CTX_DEFAULT*time.Second)
	defer cancel()

	start := c.Query("start")
	if start == "" {
		c.JSON(http.StatusNotFound, errors.Response{
			Status:  http.StatusNotFound,
			Type:    constants.ErrUnknownResource,
			Message: []string{constants.ErrMissingStartTime}})
		return
	}
	end := c.Query("end")
	if end == "" {
		c.JSON(http.StatusNotFound, errors.Response{
			Status:  http.StatusNotFound,
			Type:    constants.ErrUnknownResource,
			Message: []string{constants.ErrMissingEndTime}})
		return
	}
	result, err := h.BookingRepository.GetByDateRange(ctx, start, end)
	if err != nil {
		c.JSON(http.StatusNotFound, errors.Response{
			Status:  http.StatusNotFound,
			Type:    constants.ErrUnknownResource,
			Message: []string{err.Error()}})
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func (h *Handler) Save(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.CTX_DEFAULT*time.Second)
	defer cancel()

	var booking Booking
	if err := c.BindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, errors.Response{
			Status:  http.StatusBadRequest,
			Type:    constants.ErrRequestDecoding,
			Message: []string{err.Error()}})
		return
	} // validate before hitting the db
	if err := h.Validate.Struct(booking); err != nil {
		c.JSON(http.StatusBadRequest, errors.Response{
			Status:  http.StatusBadRequest,
			Type:    constants.ErrRequestBody,
			Message: []string{err.Error()}})
		return
	}
	id, err := h.BookingService.Save(ctx, booking)
	if err != nil {
		if err == errors.ErrInvalidTimestamp || err == errors.ErrOldTimestamp {
			c.JSON(http.StatusBadRequest, errors.Response{
				Status:  http.StatusBadRequest,
				Type:    constants.ErrRequestBody,
				Message: []string{err.Error()}})
			return
		} else {
			c.JSON(http.StatusBadRequest, errors.Response{
				Status:  http.StatusBadRequest,
				Type:    constants.ErrDatabaseOperation,
				Message: []string{err.Error(), "Check if the class date is valid"}})
		}
	} else {
		msg := "Created Booking successfully"
		c.JSON(http.StatusCreated, gin.H{
			"Status":  http.StatusCreated,
			"Id":      id,
			"Message": msg})
	}
}

func (h *Handler) Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.CTX_DEFAULT*time.Second)
	defer cancel()

	var booking Booking
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Missing 'id' Query Parameters"})
		return
	}
	if err := c.BindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, c.Error(err))
		return
	}
	if err := h.Validate.Struct(booking); err != nil {
		c.JSON(http.StatusNotFound, c.Error(err))
		return
	}
	err := h.BookingService.Update(ctx, id, booking)
	if err != nil {
		if err == errors.ErrInvalidTimestamp {
			c.JSON(http.StatusBadRequest, errors.Response{
				Status:  http.StatusBadRequest,
				Type:    constants.ErrRequestBody,
				Message: []string{err.Error()}})
			return
		} else {
			c.JSON(http.StatusBadRequest, errors.Response{
				Status:  http.StatusBadRequest,
				Type:    constants.ErrDatabaseOperation,
				Message: []string{err.Error()}})
		}
	} else {
		msg := "Updated Booking with successfully"
		c.JSON(http.StatusOK, gin.H{
			"Status":  http.StatusOK,
			"Id":      id,
			"Message": msg})
	}
}

func (h *Handler) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.CTX_DEFAULT*time.Second)
	defer cancel()

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Missing 'id' Query Parameters"})
		return
	}
	err := h.BookingRepository.Delete(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, errors.Response{
			Status:  http.StatusNotFound,
			Type:    constants.ErrUnknownResource,
			Message: []string{err.Error()}})
	} else {
		msg := "Deleted Booking successfully"
		c.JSON(http.StatusOK, gin.H{
			"Status":  http.StatusOK,
			"Id":      id,
			"Message": msg})
	}
}

func (h *Handler) GetAllClassesByMemberId(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.CTX_DEFAULT*time.Second)
	defer cancel()

	memberID := c.Param("id")
	result, err := h.BookingRepository.GetAllClassesByMemberId(ctx, memberID)
	if err != nil {
		c.JSON(http.StatusNotFound, errors.Response{
			Status:  http.StatusNotFound,
			Type:    constants.ErrUnknownResource,
			Message: []string{err.Error()}})
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func (h *Handler) GetAllMembersByClassId(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.CTX_DEFAULT*time.Second)
	defer cancel()

	classId := c.Param("id")
	result, err := h.BookingRepository.GetAllMembersByClassId(ctx, classId)
	if err != nil {
		c.JSON(http.StatusNotFound, errors.Response{
			Status:  http.StatusNotFound,
			Type:    constants.ErrUnknownResource,
			Message: []string{err.Error()}})
	} else {
		c.JSON(http.StatusOK, result)
	}
}
