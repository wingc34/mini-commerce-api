package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wingc34/mini-commerce-api/internal/handler"
)

func registerOrderRoutes(r *gin.RouterGroup, h *handler.OrderHandler) {
	orders := r.Group("/orders")
	{
		orders.POST("/draft", h.CreateDraftOrder)
		orders.GET("", h.GetOrders)
		orders.GET("/:id", h.GetOrderDetail)
	}
}
