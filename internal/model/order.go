package model

import "time"

// OrderStatus represents the lifecycle state of a confirmed order.
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusPaid      OrderStatus = "PAID"
	OrderStatusShipped   OrderStatus = "SHIPPED"
	OrderStatusCompleted OrderStatus = "COMPLETED"
	OrderStatusCanceled  OrderStatus = "CANCELED"
)

type Order struct {
	ID                string      `gorm:"primaryKey"                                       json:"id"`
	UserID            string      `gorm:"not null"                                         json:"userId"`
	User              User        `gorm:"foreignKey:UserID"                                json:"user,omitempty"`
	Total             int32       `gorm:"not null"                                         json:"total"`
	Status            OrderStatus `gorm:"not null;type:order_status;default:'PENDING'"     json:"status"`
	ShippingAddressID string      `gorm:"not null"                                         json:"shippingAddressId"`
	ShippingAddress   Address     `gorm:"foreignKey:ShippingAddressID"                     json:"shippingAddress,omitempty"`
	PaymentIntentID   *string     `                                                         json:"paymentIntentId"`
	StripeSessionID   *string     `                                                         json:"stripeSessionId"`
	CreatedAt         time.Time   `gorm:"autoCreateTime"                                   json:"createdAt"`
	DraftOrderID      string      `gorm:"uniqueIndex;not null"                             json:"draftOrderId"`
	DraftOrder        DraftOrder  `gorm:"foreignKey:DraftOrderID"                          json:"draftOrder,omitempty"`
	OrderItems        []OrderItem `gorm:"foreignKey:OrderID"                           json:"orderItems,omitempty"`
}
