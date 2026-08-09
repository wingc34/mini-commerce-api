package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wingc34/mini-commerce-api/internal/handler"
)

func registerAuthRoutes(r *gin.RouterGroup, h *handler.AuthHandler) {
	auth := r.Group("/auth")
	{
		auth.GET("/google", h.GetGoogleAuthURL)
		auth.GET("/google/callback", h.HandleGoogleCallback)
	}
}
