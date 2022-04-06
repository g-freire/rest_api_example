package booking

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
	"gym/internal/constants"
	"gym/internal/errors"
	"net/http"
	"strconv"
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
	limit := c.Query("limit")
	offset := c.Query("offset")

	result, err := h.BookingRepository.GetAll(limit, offset)
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
	id := c.Param("id")
	if _, err := strconv.Atoi(id); err != nil {
		c.JSON(http.StatusNotFound, errors.Response{
			http.StatusNotFound,
			constants.ErrUnknownResource,
			[]string{constants.ErrWrongURLParamType}})
		return
	}
	result, err := h.BookingRepository.GetByID(id)
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
	result, err := h.BookingRepository.GetTotalCount()
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
	result, err := h.BookingRepository.GetByDateRange(start, end)
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
	err := h.BookingRepository.Save(booking)
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
		c.JSON(http.StatusCreated, "Created Class Successfully")
	}
}

func (h *Handler) Update(c *gin.Context) {
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
	err := h.BookingRepository.Update(id, booking)
	if err != nil {
		c.JSON(http.StatusNotFound, c.Error(err))
	} else {
		c.JSON(http.StatusOK, "Updated Booking with id "+id+" successfully")
	}
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Missing 'id' Query Parameters"})
		return
	}
	err := h.BookingRepository.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, errors.Response{
			Status:  http.StatusNotFound,
			Type:    constants.ErrUnknownResource,
			Message: []string{err.Error()}})
	} else {
		c.JSON(http.StatusOK, "Deleted Booking with id "+id+" successfully")
	}
}

func (h *Handler) GetAllClassesByMemberId(c *gin.Context) {
	memberID := c.Param("id")
	result, err := h.BookingRepository.GetAllClassesByMemberId(memberID)
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
	classId := c.Param("id")
	result, err := h.BookingRepository.GetAllMembersByClassId(classId)
	if err != nil {
		c.JSON(http.StatusNotFound, errors.Response{
			Status:  http.StatusNotFound,
			Type:    constants.ErrUnknownResource,
			Message: []string{err.Error()}})
	} else {
		c.JSON(http.StatusOK, result)
	}
}
