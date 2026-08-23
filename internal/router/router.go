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
	paymentHandler *handler.PaymentHandler,
	jwtSecret string,
) *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	// Global middleware
	r.Use(middleware.CORS())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public routes
	public := r.Group("/api/v1")
	registerAuthRoutes(public, authHandler)
	registerProductRoutes(public, productHandler)

	// Protected routes
	protected := r.Group("/api/v1")
	protected.Use(middleware.Auth(jwtSecret))
	protected.POST("/payments/intent", paymentHandler.CreatePaymentIntent)
	registerUserRoutes(protected, userHandler)
	registerOrderRoutes(protected, orderHandler)

	// Webhook
	r.POST("/api/v1/webhooks/stripe", webhookHandler.HandleStripeWebhook)

	return r
}
