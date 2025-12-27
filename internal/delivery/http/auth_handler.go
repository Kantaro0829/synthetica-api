package http

import (
	"net/http"
	"synthetica/internal/config"
	"synthetica/internal/domain"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Usecase domain.UserUsecase
}

func NewAuthHandler(r *gin.Engine, us domain.UserUsecase) {
	handler := &AuthHandler{
		Usecase: us,
	}
	r.GET("/auth/google/login", handler.Login)
	r.GET("/auth/google/callback", handler.Callback)
}

func (h *AuthHandler) Login(c *gin.Context) {
	url := config.GoogleOauthConfig.AuthCodeURL("randomstate")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) Callback(c *gin.Context) {
	state := c.Query("state")
	if state != "randomstate" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "states don't match"})
		return
	}

	code := c.Query("code")
	token, err := config.GoogleOauthConfig.Exchange(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Code exchange failed: " + err.Error()})
		return
	}

	user, err := h.Usecase.LoginWithGoogleOAuth(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// For MVP, just return JSON. simpler than setting cookies and redirecting to 3000 which might have CORS issues if not handled.
	// But the plan said redirect.
	// Let's set a simple cookie and redirect.
	// WARN: HttpOnly cookie won't be readable by JS easily unless we have an endpoint to fetch "me".
	// Let's just redirect with a query param for now? No, that's insecure.
	// Let's redirect to home and expect the user to be "logged in" via cookie.

	// Create a dummy session cookie for now
	c.SetCookie("user_id", user.GoogleID, 3600, "/", "localhost", false, true)
	c.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/hello")
}
