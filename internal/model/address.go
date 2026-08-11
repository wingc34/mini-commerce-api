package model

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	UserID    string         `gorm:"not null" json:"userId"`
	FullName  string         `gorm:"not null" json:"fullName"`
	Phone     string         `gorm:"not null" json:"phone"`
	Line1     string         `gorm:"not null" json:"line1"`
	Line2     *string        `json:"line2"`
	City      string         `gorm:"not null" json:"city"`
	State     *string        `json:"state"`
	Postal    string         `gorm:"not null" json:"postal"`
	Country   string         `gorm:"not null" json:"country"`
	IsDefault bool           `gorm:"not null;default:false" json:"isDefault"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
