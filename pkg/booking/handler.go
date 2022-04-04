package booking

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

type Handler struct {
	Validate          *validator.Validate
	BookingRepository BookingRepository
}

func NewHandler(r *gin.Engine,
	route string,
	val *validator.Validate,
	BookingRepo BookingRepository,
) {
	handler := &Handler{
		Validate:          val,
		BookingRepository: BookingRepo,
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
	limit := c.Query("limit")
	offset := c.Query("offset")

	result, err := h.BookingRepository.GetAll(limit, offset)
	if err != nil {
		c.JSON(http.StatusNotFound, c.Error(err))
	} else {
		c.JSON(http.StatusOK, result)
	}
}
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Missing 'id' Query Parameters"})
		return
	}
	result, err := h.BookingRepository.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, c.Error(err))
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func (h *Handler) GetTotalCount(c *gin.Context) {
	result, err := h.BookingRepository.GetTotalCount()
	if err != nil {
		c.JSON(http.StatusNotFound, c.Error(err))
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func (h *Handler) GetByDateRange(c *gin.Context) {
	start := c.Query("start")
	if start == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Missing 'start' Query Parameters"})
	}
	end := c.Query("end")
	if end == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Missing 'end' Query Parameters"})
	}
	result, err := h.BookingRepository.GetByDateRange(start, end)
	if err != nil {
		c.JSON(http.StatusNotFound, c.Error(err))
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func (h *Handler) Save(c *gin.Context) {
	var class Booking
	err := c.BindJSON(&class)
	if err != nil {
		c.JSON(http.StatusBadRequest, c.Error(err))
		return
	}
	if err := h.Validate.Struct(class); err != nil {
		c.JSON(http.StatusNotFound, c.Error(err))
		return
	}
	err = h.BookingRepository.Save(class)
	if err != nil {
		c.JSON(http.StatusNotFound, c.Error(err))
	} else {
		c.JSON(http.StatusOK, "Created Booking Successfully")
	}
}

func (h *Handler) Update(c *gin.Context) {
	var class Booking
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Missing 'id' Query Parameters"})
		return
	}
	if err := c.BindJSON(&class); err != nil {
		c.JSON(http.StatusBadRequest, c.Error(err))
		return
	}
	if err := h.Validate.Struct(class); err != nil {
		c.JSON(http.StatusNotFound, c.Error(err))
		return
	}
	err := h.BookingRepository.Update(id, class)
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
		c.JSON(http.StatusNotFound, c.Error(err))
	} else {
		c.JSON(http.StatusOK, "Deleted Booking with id "+id+" successfully")
	}
}
