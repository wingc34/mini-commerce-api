package service

import (
	"encoding/json"
	"fmt"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"github.com/stripe/stripe-go/v82/webhook"
	"github.com/wingc34/mini-commerce-api/internal/model"
	"github.com/wingc34/mini-commerce-api/internal/repository"
	"github.com/wingc34/mini-commerce-api/pkg/id"
)

type PaymentService interface {
	CreatePaymentIntent(amount int64, draftOrderID string) (string, error)
	HandleWebhook(payload []byte, sigHeader string) error
}

type paymentService struct {
	orderRepo     repository.OrderRepository
	stripeKey     string
	webhookSecret string
}

func NewPaymentService(
	orderRepo repository.OrderRepository,
	stripeKey string,
	webhookSecret string,
) PaymentService {
	return &paymentService{
		orderRepo:     orderRepo,
		stripeKey:     stripeKey,
		webhookSecret: webhookSecret,
	}
}

// CreatePaymentIntent creates a Stripe PaymentIntent and returns the clientSecret.
// The draftOrderId is stored in metadata so the webhook can find it later.
func (s *paymentService) CreatePaymentIntent(amount int64, draftOrderID string) (string, error) {
	stripe.Key = s.stripeKey

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String("hkd"),
		// Store draftOrderId in metadata so webhook can look it up on payment success
		Metadata: map[string]string{
			"draftOrderId": draftOrderID,
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return "", fmt.Errorf("failed to create payment intent: %w", err)
	}

	return pi.ClientSecret, nil
}

// HandleWebhook processes incoming Stripe webhook events.
// On payment_intent.succeeded, it creates a real Order from the DraftOrder.
func (s *paymentService) HandleWebhook(payload []byte, sigHeader string) error {
	// 第一步：驗證 Stripe 簽名，防止假冒的 webhook 請求
	event, err := webhook.ConstructEvent(payload, sigHeader, s.webhookSecret)
	if err != nil {
		return fmt.Errorf("invalid webhook signature: %w", err)
	}

	// 第二步：只處理付款成功的事件，其他事件忽略
	if event.Type != "payment_intent.succeeded" {
		return nil
	}

	// 第三步：解析 PaymentIntent 拿到 metadata
	var pi stripe.PaymentIntent
	if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
		return fmt.Errorf("failed to parse payment intent: %w", err)
	}

	draftOrderID := pi.Metadata["draftOrderId"]
	if draftOrderID == "" {
		return fmt.Errorf("draftOrderId not found in metadata")
	}

	// 第四步：找到 DraftOrder
	draft, err := s.orderRepo.FindDraftByID(draftOrderID)
	if err != nil {
		return fmt.Errorf("failed to find draft order: %w", err)
	}

	// 第五步：找到 DraftOrderItems
	draftItems, err := s.orderRepo.FindDraftItemsByDraftID(draftOrderID)
	if err != nil {
		return fmt.Errorf("failed to find draft order items: %w", err)
	}

	// 第六步：把 DraftOrderItems 轉成 OrderItems
	var orderItems []model.OrderItem
	for _, item := range draftItems {
		orderItems = append(orderItems, model.OrderItem{
			ID:       id.New(),
			SKUID:    item.SKUID,
			Quantity: item.Quantity,
			Price:    item.Price,
		})
	}

	// 第七步：建立 Order
	order := &model.Order{
		ID:                id.New(),
		UserID:            draft.UserID,
		Total:             draft.Total,
		ShippingAddressID: draft.ShippingAddressID,
		DraftOrderID:      draft.ID,
	}

	if err := s.orderRepo.CreateOrder(order, orderItems); err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	// 第八步：更新 DraftOrder 狀態為 COMPLETED
	return s.orderRepo.UpdateDraftStatus(draftOrderID, model.DraftStatusCompleted)
}
