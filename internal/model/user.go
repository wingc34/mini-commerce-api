package model

import "time"

type User struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Email       string    `gorm:"uniqueIndex;not null" json:"email"`
	Name        *string   `json:"name"`
	Image       *string   `json:"image"`
	PhoneNumber *string   `json:"phoneNumber"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}
