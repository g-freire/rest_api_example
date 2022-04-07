package member

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
	Validate         *validator.Validate
	MemberRepository MemberRepository
}

func NewHandler(r *gin.Engine,
	route string,
	val *validator.Validate,
	MemberRepo MemberRepository,
) {
	handler := &Handler{
		Validate:         val,
		MemberRepository: MemberRepo,
	}
	v1 := r.Group("v1/" + route)
	{
		v1.GET("", handler.GetAll)
		v1.GET(":id", handler.GetByID)
		v1.GET("count", handler.GetTotalCount)
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

	result, err := h.MemberRepository.GetAll(ctx, limit, offset, name)
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
	result, err := h.MemberRepository.GetByID(ctx, id)
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

	result, err := h.MemberRepository.GetTotalCount(ctx)
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

	var member Member
	if err := c.BindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, errors.Response{
			Status:  http.StatusBadRequest,
			Type:    constants.ErrRequestDecoding,
			Message: []string{err.Error()}})
		return
	} // validates before hitting the db
	if err := h.Validate.Struct(member); err != nil {
		c.JSON(http.StatusBadRequest, errors.Response{
			Status:  http.StatusBadRequest,
			Type:    constants.ErrRequestBody,
			Message: []string{err.Error()}})
		return
	}
	id, err := h.MemberRepository.Save(ctx, member)
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

	var member Member
	id := c.Param("id")
	if err := c.BindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, errors.Response{
			Status:  http.StatusBadRequest,
			Type:    constants.ErrRequestDecoding,
			Message: []string{err.Error()}})
		return
	} // validates before hitting the db
	if err := h.Validate.Struct(member); err != nil {
		c.JSON(http.StatusBadRequest, errors.Response{
			Status:  http.StatusBadRequest,
			Type:    constants.ErrRequestBody,
			Message: []string{err.Error()}})
		return
	}
	err := h.MemberRepository.Update(ctx, id, member)
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
		c.JSON(http.StatusOK, "Updated Member with id "+id+" successfully")
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
	err := h.MemberRepository.Delete(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, errors.Response{
			Status:  http.StatusNotFound,
			Type:    constants.ErrUnknownResource,
			Message: []string{err.Error()}})
	} else {
		c.JSON(http.StatusOK, "Deleted Member with id "+id+" successfully")
	}
}
