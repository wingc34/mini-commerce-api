package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wingc34/mini-commerce-api/internal/model"
	"github.com/wingc34/mini-commerce-api/internal/service"
	"github.com/wingc34/mini-commerce-api/pkg/ctxutil"
	"github.com/wingc34/mini-commerce-api/pkg/id"
	"github.com/wingc34/mini-commerce-api/pkg/response"
)

type UserHandler struct {
	service service.UserService
}
type UpdateMeRequest struct {
	Name        *string `json:"name"`
	PhoneNumber *string `json:"phoneNumber"`
}

type AddressRequest struct {
	FullName  string  `json:"fullName"  binding:"required"`
	Phone     string  `json:"phone"     binding:"required"`
	Line1     string  `json:"line1"     binding:"required"`
	Line2     *string `json:"line2"`
	City      string  `json:"city"      binding:"required"`
	State     *string `json:"state"`
	Postal    string  `json:"postal"    binding:"required"`
	Country   string  `json:"country"   binding:"required"`
	IsDefault bool    `json:"isDefault"`
}

type WishItemRequest struct {
	ProductID string `json:"productId" binding:"required"`
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := ctxutil.GetUserID(c)

	user, err := h.service.GetMe(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get user")
		return
	}

	totalOrders, totalSpent, err := h.service.GetOrderStats(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get order stats")
		return
	}

	var defaultAddress *model.Address
	if len(user.Addresses) > 0 {
		defaultAddress = &user.Addresses[0]
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":             user.ID,
			"email":          user.Email,
			"name":           user.Name,
			"image":          user.Image,
			"phoneNumber":    user.PhoneNumber,
			"wishlist":       user.Wishlist,
			"defaultAddress": defaultAddress,
			"totalOrders":    totalOrders,
			"totalSpent":     totalSpent,
			"createdAt":      user.CreatedAt,
		},
	})
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID := ctxutil.GetUserID(c)

	// 解析 request body
	var req UpdateMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	// 呼叫 service
	if err := h.service.UpdateMe(userID, req.Name, req.PhoneNumber); err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to update user")
		return
	}

	response.Success(c, gin.H{"message": "updated successfully"})

}

func (h *UserHandler) GetAddresses(c *gin.Context) {
	userID := ctxutil.GetUserID(c)

	address, err := h.service.GetAddresses(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get addresses")
		return
	}

	response.Success(c, address)
}

func (h *UserHandler) CreateAddress(c *gin.Context) {
	userID := ctxutil.GetUserID(c)

	var req AddressRequest
	if !ctxutil.BindJSON(c, &req) {
		return
	}

	// 建立 model
	address := &model.Address{
		ID:        id.New(),
		UserID:    userID,
		FullName:  req.FullName,
		Phone:     req.Phone,
		Line1:     req.Line1,
		Line2:     req.Line2,
		City:      req.City,
		State:     req.State,
		Postal:    req.Postal,
		Country:   req.Country,
		IsDefault: req.IsDefault,
	}

	if err := h.service.CreateAddress(address); err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to create address")
		return
	}

	response.Success(c, address)
}

func (h *UserHandler) UpdateAddress(c *gin.Context) {
	userID := ctxutil.GetUserID(c)

	var req AddressRequest
	if !ctxutil.BindJSON(c, &req) {
		return
	}

	// 建立 model
	address := &model.Address{
		ID:        c.Param("id"),
		UserID:    userID,
		FullName:  req.FullName,
		Phone:     req.Phone,
		Line1:     req.Line1,
		Line2:     req.Line2,
		City:      req.City,
		State:     req.State,
		Postal:    req.Postal,
		Country:   req.Country,
		IsDefault: req.IsDefault,
	}

	if err := h.service.UpdateAddress(address); err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to update address")
		return
	}

	response.Success(c, address)
}

func (h *UserHandler) DeleteAddress(c *gin.Context) {
	userID := ctxutil.GetUserID(c)
	addressID := c.Param("id")

	if err := h.service.DeleteAddress(addressID, userID); err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to delete address")
		return
	}

	response.Success(c, gin.H{"message": "deleted successfully"})
}

func (h *UserHandler) AddWishItem(c *gin.Context) {
	userID := ctxutil.GetUserID(c)

	var req WishItemRequest
	if !ctxutil.BindJSON(c, &req) {
		return
	}

	if err := h.service.AddWishItem(userID, req.ProductID); err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to add wish item")
		return
	}

	response.Success(c, gin.H{"message": "added successfully"})
}

func (h *UserHandler) RemoveWishItem(c *gin.Context) {
	userID := ctxutil.GetUserID(c)
	productID := c.Param("productId")

	if err := h.service.RemoveWishItem(userID, productID); err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to remove wish item")
		return
	}

	response.Success(c, gin.H{"message": "removed successfully"})
}
