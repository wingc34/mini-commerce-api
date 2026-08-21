package model

type Wishlist struct {
	UserID    string  `gorm:"primaryKey;not null" json:"-"`
	ProductID string  `gorm:"primaryKey;not null" json:"productId"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}
