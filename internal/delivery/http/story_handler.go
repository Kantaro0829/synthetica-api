package http

import (
	"net/http"
	"strconv"
	"synthetica/internal/domain"

	"github.com/gin-gonic/gin"
)

type StoryHandler struct {
	Usecase domain.StoryUsecase
}

func NewStoryHandler(r *gin.Engine, us domain.StoryUsecase) {
	handler := &StoryHandler{
		Usecase: us,
	}
	r.POST("/stories", handler.Store)
	r.GET("/stories", handler.Fetch)
	r.POST("/stories/:id/like", handler.Like)
}

type storeStoryRequest struct {
	Title  string `json:"title" binding:"required"`
	Detail string `json:"detail" binding:"required"`
}

type likeStoryRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}

func (h *StoryHandler) Like(c *gin.Context) {
	// Get user_id from cookie
	var userID uint
	cookie, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.Atoi(cookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}
	userID = uint(id)

	idParam := c.Param("id")
	storyID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid story ID"})
		return
	}

	ctx := c.Request.Context()
	err = h.Usecase.ToggleLike(ctx, uint(storyID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *StoryHandler) Store(c *gin.Context) {
	var req storeStoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	story := &domain.Story{
		Title:  req.Title,
		Detail: req.Detail,
		UserID: 1, // Test user ID as requested
	}

	ctx := c.Request.Context()
	err := h.Usecase.Create(ctx, story)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, story)
}

func (h *StoryHandler) Fetch(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user_id from cookie to check liked status
	var userID uint
	if cookie, err := c.Cookie("user_id"); err == nil {
		if id, err := strconv.Atoi(cookie); err == nil {
			userID = uint(id)
		}
	}

	stories, err := h.Usecase.Fetch(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stories)
}
