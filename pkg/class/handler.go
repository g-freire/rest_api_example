package class

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
	Validate        *validator.Validate
	ClassService    ClassService
	ClassRepository ClassRepository
}

func NewHandler(r *gin.Engine,
	route string,
	val *validator.Validate,
	ClassService ClassService,
	ClassRepo ClassRepository,
) {
	handler := &Handler{
		Validate:        val,
		ClassService:    ClassService,
		ClassRepository: ClassRepo,
	}

	v1 := r.Group("v1/" + route)
	{
		v1.GET("", handler.GetAll)
		v1.GET(":id", handler.GetByID)
		v1.GET("count", handler.GetTotalCount)
		v1.GET("date", handler.GetByDateRange)
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
	name := c.Query("name")

	result, err := h.ClassRepository.GetAll(ctx, limit, offset, name)
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
	result, err := h.ClassRepository.GetByID(ctx, id)
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

	result, err := h.ClassRepository.GetTotalCount(ctx)
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
	result, err := h.ClassService.GetByDateRange(ctx, start, end)
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

	var class Class
	if err := c.BindJSON(&class); err != nil {
		c.JSON(http.StatusBadRequest, errors.Response{
			Status:  http.StatusBadRequest,
			Type:    constants.ErrRequestDecoding,
			Message: []string{err.Error()}})
		return
	} // validates before hitting the db
	if err := h.Validate.Struct(class); err != nil {
		c.JSON(http.StatusBadRequest, errors.Response{
			Status:  http.StatusBadRequest,
			Type:    constants.ErrRequestBody,
			Message: []string{err.Error()}})
		return
	}
	id, err := h.ClassService.Save(ctx, class)
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
				Message: []string{err.Error()}})
		}
	} else {
		msg := "Created Class successfully"
		c.JSON(http.StatusCreated, gin.H{
			"Status":  http.StatusCreated,
			"Id":      id,
			"Message": msg})
	}
}

func (h *Handler) Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), constants.CTX_DEFAULT*time.Second)
	defer cancel()

	var class Class
	id := c.Param("id")
	if err := c.BindJSON(&class); err != nil {
		c.JSON(http.StatusBadRequest, errors.Response{
			Status:  http.StatusBadRequest,
			Type:    constants.ErrRequestDecoding,
			Message: []string{err.Error()}})
		return
	} // validates before hitting the db
	if err := h.Validate.Struct(class); err != nil {
		c.JSON(http.StatusBadRequest, errors.Response{
			Status:  http.StatusBadRequest,
			Type:    constants.ErrRequestBody,
			Message: []string{err.Error()}})
		return
	}
	err := h.ClassService.Update(ctx, id, class)
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
		msg := "Updated Class with successfully"
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
	if _, err := strconv.Atoi(id); err != nil {
		c.JSON(http.StatusNotFound, errors.Response{
			Status:  http.StatusNotFound,
			Type:    constants.ErrUnknownResource,
			Message: []string{constants.ErrWrongURLParamType}})
		return
	}
	err := h.ClassRepository.Delete(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, errors.Response{
			Status:  http.StatusNotFound,
			Type:    constants.ErrUnknownResource,
			Message: []string{err.Error()}})
	} else {
		msg := "Deleted Class successfully"
		c.JSON(http.StatusOK, gin.H{
			"Status":  http.StatusOK,
			"Id":      id,
			"Message": msg})
	}
}
