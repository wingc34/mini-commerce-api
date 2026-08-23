package handler

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wingc34/mini-commerce-api/internal/service"
	"github.com/wingc34/mini-commerce-api/pkg/response"
)

type WebhookHandler struct {
	service service.PaymentService
}

func NewWebhookHandler(service service.PaymentService) *WebhookHandler {
	return &WebhookHandler{service: service}
}

func (h *WebhookHandler) HandleStripeWebhook(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)

	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to read request body")
		return
	}

	sigHeader := c.GetHeader("Stripe-Signature")

	err = h.service.HandleWebhook(payload, sigHeader)
	if err != nil {
		log.Printf("❌ Webhook error: %v", err)
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}
