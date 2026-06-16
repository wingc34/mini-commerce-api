package model

type Wishlist struct {
	UserID    string  `gorm:"primaryKey;not null" json:"userId"`
	ProductID string  `gorm:"primaryKey;not null" json:"productId"`
	User      User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}
