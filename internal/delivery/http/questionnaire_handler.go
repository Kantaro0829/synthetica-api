package http

import (
	"net/http"
	"synthetica/internal/domain"

	"github.com/gin-gonic/gin"
)

type QuestionnaireHandler struct {
	Usecase domain.QuestionnaireUsecase
}

func NewQuestionnaireHandler(r *gin.Engine, us domain.QuestionnaireUsecase) {
	handler := &QuestionnaireHandler{
		Usecase: us,
	}
	r.POST("/questionnaire", handler.Store)
	r.GET("/questionnaire/status", handler.GetStatus)
}

type StoreRequest struct {
	Answer int `json:"answer" binding:"required"`
}

func (h *QuestionnaireHandler) Store(c *gin.Context) {
	// 1. Get GoogleID from cookie
	googleID, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 2. Bind JSON
	var req StoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Call Usecase
	ctx := c.Request.Context()
	err = h.Usecase.Store(ctx, googleID, req.Answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Answer stored"})
}

func (h *QuestionnaireHandler) GetStatus(c *gin.Context) {
	googleID, err := c.Cookie("user_id")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := c.Request.Context()
	q, err := h.Usecase.GetStatus(ctx, googleID)

	if err != nil {
		// Log error internally if needed
		// For client, if it's a "user not found" issue, maybe 404 or just answered: false
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if q != nil {
		c.JSON(http.StatusOK, gin.H{
			"answered":    true,
			"user_answer": q.Answer,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"answered":    false,
			"user_answer": nil,
		})
	}
}
