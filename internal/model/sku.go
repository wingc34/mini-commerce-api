package model

import (
	"time"

	"gorm.io/datatypes"
)

type SKU struct {
	ID         string         `gorm:"primaryKey" json:"id"`
	ProductID  string         `gorm:"not null" json:"productId"`
	Product    Product        `gorm:"foreignKey:ProductID" json:"product,omitempty"` // omitempty mean no output in json when value is empty
	SKUCode    string         `gorm:"uniqueIndex;not null" json:"skuCode"`
	Price      int32          `gorm:"not null" json:"price"`
	Stock      int32          `gorm:"not null;default:0" json:"stock"`
	Attributes datatypes.JSON `gorm:"type:jsonb;not null" json:"attributes"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
}
