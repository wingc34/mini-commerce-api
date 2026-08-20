package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/wingc34/mini-commerce-api/internal/service"
	"github.com/wingc34/mini-commerce-api/pkg/response"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// GetGoogleAuthURL redirects the user to Google's OAuth consent screen.
func (h *AuthHandler) GetGoogleAuthURL(c *gin.Context) {
	url := h.service.GetGoogleAuthURL()

	c.Redirect(http.StatusTemporaryRedirect, url)
}

// HandleGoogleCallback exchanges the OAuth code for a JWT token.
func (h *AuthHandler) HandleGoogleCallback(c *gin.Context) {
	// 拿到 Google 回傳的 code
	code := c.Query("code")

	// 換取 JWT token
	token, err := h.service.HandleGoogleCallback(code)

	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	c.Redirect(http.StatusTemporaryRedirect,
		fmt.Sprintf("%s/login-confirm?token=%s", frontendURL, token))

	// 回傳 JWT 給前端
	response.Success(c, gin.H{"token": token})
}
