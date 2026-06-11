package model

import (
	"time"

	"github.com/lib/pq"
)

type Product struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description *string        `json:"description"`
	Images      pq.StringArray `gorm:"type:text[];not null;default:'{}'" json:"images"`
	Category    *string        `json:"category"`
	SKUs        []SKU          `gorm:"foreignKey:ProductID" json:"skus,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
}
