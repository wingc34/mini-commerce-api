package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wingc34/mini-commerce-api/internal/service"
	"github.com/wingc34/mini-commerce-api/pkg/ctxutil"
	"github.com/wingc34/mini-commerce-api/pkg/response"
)

type PaymentHandler struct {
	service service.PaymentService
}

func NewPaymentHandler(service service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

type CreatePaymentIntentRequest struct {
	Amount       int64  `json:"amount"       binding:"required"`
	DraftOrderID string `json:"draftOrderId" binding:"required"`
}

func (h *PaymentHandler) CreatePaymentIntent(c *gin.Context) {
	var req CreatePaymentIntentRequest
	if !ctxutil.BindJSON(c, &req) {
		return
	}

	clientSecret, err := h.service.CreatePaymentIntent(req.Amount, req.DraftOrderID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to create payment intent")
		return
	}

	response.Success(c, gin.H{"clientSecret": clientSecret})
}
