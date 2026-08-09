package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/wingc34/mini-commerce-api/internal/config"
	"github.com/wingc34/mini-commerce-api/internal/database"
	"github.com/wingc34/mini-commerce-api/internal/handler"
	"github.com/wingc34/mini-commerce-api/internal/repository"
	"github.com/wingc34/mini-commerce-api/internal/router"
	"github.com/wingc34/mini-commerce-api/internal/service"
	"github.com/wingc34/mini-commerce-api/pkg/oauth"
)

func main() {
	// 1. 載入 config
	cfg := config.Load()

	// 2. 連接 database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 3. 初始化所有 repository
	userRepo := repository.NewUserRepository(db)
	addressRepo := repository.NewAddressRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// 4. 初始化所有 service
	oauthConfig := oauth.NewGoogleOAuthConfig(
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.GoogleRedirectURL,
	)
	authService := service.NewAuthService(userRepo, oauthConfig, cfg.JWTSecret)
	userService := service.NewUserService(userRepo, addressRepo)
	productService := service.NewProductService(productRepo)
	orderService := service.NewOrderService(orderRepo)
	paymentService := service.NewPaymentService(orderRepo, cfg.StripeSecretKey, cfg.StripeWebhookSecret)

	// 5. 初始化所有 handler
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService)
	orderHandler := handler.NewOrderHandler(orderService)
	webhookHandler := handler.NewWebhookHandler(paymentService)

	// 6. 設定 Gin mode
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 7. 設定 router
	r := router.SetupRouter(
		authHandler,
		productHandler,
		userHandler,
		orderHandler,
		webhookHandler,
		cfg.JWTSecret,
	)

	// 8. 啟動 server
	addr := ":" + cfg.Port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
