package model

type DraftOrderItem struct {
	ID           string     `gorm:"primaryKey" json:"id"`
	DraftOrderID string     `gorm:"not null" json:"draftOrderId"`
	DraftOrder   DraftOrder `gorm:"foreignKey:DraftOrderID" json:"draftOrder,omitempty"`
	SKUID        string     `gorm:"not null;column:sku_id" json:"skuId"`
	SKU          SKU        `gorm:"foreignKey:SKUID" json:"sku,omitempty"`
	Quantity     int32      `gorm:"not null" json:"quantity"`
	Price        int32      `gorm:"not null" json:"price"`
}
