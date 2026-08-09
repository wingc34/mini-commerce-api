package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wingc34/mini-commerce-api/internal/model"
	"github.com/wingc34/mini-commerce-api/internal/service"
	"github.com/wingc34/mini-commerce-api/pkg/ctxutil"
	"github.com/wingc34/mini-commerce-api/pkg/id"
	"github.com/wingc34/mini-commerce-api/pkg/response"
)

type CreateDraftOrderRequest struct {
	Total             int32                   `json:"total"             binding:"required"`
	ShippingAddressID string                  `json:"shippingAddressId" binding:"required"`
	OrderItems        []DraftOrderItemRequest `json:"orderItem"         binding:"required"`
}

type DraftOrderItemRequest struct {
	SKUID    string `json:"skuId"    binding:"required"`
	Quantity int32  `json:"quantity" binding:"required"`
	Price    int32  `json:"price"    binding:"required"`
}

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	userID := ctxutil.GetUserID(c)
	pageStr := c.DefaultQuery("page", "1")
	page, _ := strconv.Atoi(pageStr)

	orders, total, err := h.service.GetOrders(userID, page)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get orders")
		return
	}

	response.Success(c, gin.H{"orders": orders, "total": total})
}

func (h *OrderHandler) GetOrderDetail(c *gin.Context) {
	id := c.Param("id")

	order, draftOrder, err := h.service.GetOrderDetail(id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get order detail")
		return
	}

	// 找到 order 就回傳 order，找到 draft 就回傳 draft
	if order != nil {
		response.Success(c, order)
		return
	}

	response.Success(c, draftOrder)
}

func (h *OrderHandler) CreateDraftOrder(c *gin.Context) {
	userID := ctxutil.GetUserID(c)

	var req CreateDraftOrderRequest
	if !ctxutil.BindJSON(c, &req) {
		return
	}

	// 建立 DraftOrder
	draft := &model.DraftOrder{
		ID:                id.New(),
		UserID:            userID,
		Total:             req.Total,
		ShippingAddressID: req.ShippingAddressID,
	}

	// 建立 DraftOrderItems
	var items []model.DraftOrderItem
	for _, item := range req.OrderItems {
		items = append(items, model.DraftOrderItem{
			ID:           id.New(),
			DraftOrderID: draft.ID,
			SKUID:        item.SKUID,
			Quantity:     item.Quantity,
			Price:        item.Price,
		})
	}

	if err := h.service.CreateDraftOrder(draft, items); err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to create draft order")
		return
	}

	response.Success(c, draft)
}
