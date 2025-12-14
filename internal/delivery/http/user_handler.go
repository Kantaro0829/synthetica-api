package http

import (
	"net/http"
	"strconv"
	"synthetica/internal/domain"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Usecase domain.UserUsecase
}

func NewUserHandler(r *gin.Engine, us domain.UserUsecase) {
	handler := &UserHandler{
		Usecase: us,
	}
	r.POST("/users", handler.Store)
	r.GET("/users/:id", handler.GetByID)
	r.GET("/users", handler.Fetch)
}

func (h *UserHandler) Store(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	err := h.Usecase.Store(ctx, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	ctx := c.Request.Context()
	user, err := h.Usecase.GetByID(ctx, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Fetch(c *gin.Context) {
	ctx := c.Request.Context()
	users, err := h.Usecase.Fetch(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
