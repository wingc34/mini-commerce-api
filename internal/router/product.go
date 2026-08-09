package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wingc34/mini-commerce-api/internal/handler"
)

func registerProductRoutes(r *gin.RouterGroup, h *handler.ProductHandler) {
	products := r.Group("/products")
	{
		products.GET("", h.GetProducts)
		products.GET("/recommended", h.GetRecommended)
		products.GET("/:id", h.GetProductByID)
		products.POST("/:id/stock", h.CheckStock)
	}
}
