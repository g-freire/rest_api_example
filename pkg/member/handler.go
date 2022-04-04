package booking

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
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
	limit := c.Query("limit")
	offset := c.Query("offset")
	name := c.Query("name")

	result, err := h.MemberRepository.GetAll(limit, offset, name)
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
	result, err := h.MemberRepository.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, c.Error(err))
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func (h *Handler) GetTotalCount(c *gin.Context) {
	result, err := h.MemberRepository.GetTotalCount()
	if err != nil {
		c.JSON(http.StatusNotFound, c.Error(err))
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func (h *Handler) Save(c *gin.Context) {
	var member Member
	err := c.BindJSON(&member)
	if err != nil {
		c.JSON(http.StatusBadRequest, c.Error(err))
		return
	}
	if err := h.Validate.Struct(member); err != nil {
		c.JSON(http.StatusNotFound, c.Error(err))
		return
	}
	err = h.MemberRepository.Save(member)
	if err != nil {
		c.JSON(http.StatusNotFound, c.Error(err))
	} else {
		c.JSON(http.StatusOK, "Created Member Successfully")
	}
}

func (h *Handler) Update(c *gin.Context) {
	var member Member
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Missing 'id' Query Parameters"})
		return
	}
	if err := c.BindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, c.Error(err))
		return
	}
	if err := h.Validate.Struct(member); err != nil {
		c.JSON(http.StatusNotFound, c.Error(err))
		return
	}
	err := h.MemberRepository.Update(id, member)
	if err != nil {
		c.JSON(http.StatusNotFound, c.Error(err))
	} else {
		c.JSON(http.StatusOK, "Updated Member with id "+id+" successfully")
	}
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Missing 'id' Query Parameters"})
		return
	}
	err := h.MemberRepository.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, c.Error(err))
	} else {
		c.JSON(http.StatusOK, "Deleted Member with id "+id+" successfully")
	}
}
