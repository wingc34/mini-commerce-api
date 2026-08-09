package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wingc34/mini-commerce-api/internal/handler"
	"github.com/wingc34/mini-commerce-api/internal/middleware"
)

func SetupRouter(
	authHandler *handler.AuthHandler,
	productHandler *handler.ProductHandler,
	userHandler *handler.UserHandler,
	orderHandler *handler.OrderHandler,
	webhookHandler *handler.WebhookHandler,
	jwtSecret string,
) *gin.Engine {
	r := gin.Default()

	// Public routes
	public := r.Group("/api/v1")
	registerAuthRoutes(public, authHandler)
	registerProductRoutes(public, productHandler)

	// Protected routes
	protected := r.Group("/api/v1")
	protected.Use(middleware.Auth(jwtSecret))
	registerUserRoutes(protected, userHandler)
	registerOrderRoutes(protected, orderHandler)

	// Webhook
	r.POST("/api/v1/webhooks/stripe", webhookHandler.HandleStripeWebhook)

	return r
}
