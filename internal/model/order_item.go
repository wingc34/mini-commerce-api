package model

type OrderItem struct {
	ID       string `gorm:"primaryKey"        json:"id"`
	OrderID  string `gorm:"not null"          json:"orderId"`
	Order    Order  `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	SKUID    string `gorm:"not null"          json:"skuId"`
	SKU      SKU    `gorm:"foreignKey:SKUID"  json:"sku,omitempty"`
	Quantity int32  `gorm:"not null"          json:"quantity"`
	Price    int32  `gorm:"not null"          json:"price"`
}
