package model

import "time"

// DraftStatus represents the lifecycle state of a draft order
// before it is converted into a real Order after successful payment.
type DraftStatus string

const (
	DraftStatusPendingPayment DraftStatus = "PENDING_PAYMENT"
	DraftStatusPaymentFailed  DraftStatus = "PAYMENT_FAILED"
	DraftStatusCompleted      DraftStatus = "COMPLETED"
	DraftStatusExpired        DraftStatus = "EXPIRED"
)

type DraftOrder struct {
	ID                string      `gorm:"primaryKey" json:"id"`
	UserID            string      `gorm:"not null" json:"user_id"`
	User              User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Total             int32       `gorm:"not null" json:"total"`
	Status            DraftStatus `gorm:"not null;type:draft_status;default:'PENDING_PAYMENT'" json:"status"`
	ShippingAddressID string      `gorm:"not null" json:"shippingAddressId"`
	ShippingAddress   Address     `gorm:"foreignKey:ShippingAddressID" json:"shipping_address,omitempty"`
	PaymentIntentID   *string     `json:"paymentIntentId"`
	StripeSessionID   *string     `json:"stripeSessionId"`
	CreatedAt         time.Time   `gorm:"autoCreateTime" json:"created_at"`
	ExpiresAt         *time.Time  `json:"expires_at"`
}
