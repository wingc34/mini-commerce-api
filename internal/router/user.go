package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wingc34/mini-commerce-api/internal/handler"
)

func registerUserRoutes(r *gin.RouterGroup, h *handler.UserHandler) {
	users := r.Group("/users/me")
	{
		users.GET("", h.GetMe)
		users.PATCH("", h.UpdateMe)
		users.GET("/addresses", h.GetAddresses)
		users.POST("/addresses", h.CreateAddress)
		users.PUT("/addresses/:id", h.UpdateAddress)
		users.DELETE("/addresses/:id", h.DeleteAddress)
		users.POST("/wishlist", h.AddWishItem)
		users.DELETE("/wishlist/:productId", h.RemoveWishItem)
	}
}
