package repository

import (
	"github.com/wingc34/mini-commerce-api/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UserRepository defines the database operations for users.
// Defining this as an interface (not a concrete struct) lets us inject
// a mock implementation in tests without touching a real database.
type UserRepository interface {
	FindByID(id string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	UpsertByGoogle(user *model.User) error // Google OAuth 登入時用
	Update(user *model.User) error
	AddWishItem(userID string, productID string) error
	RemoveWishItem(userID string, productID string) error
}

// userRepository is the GORM implementation of UserRepository.
// Lowercase name keeps it private - external packages can only use
// the UserRepository interface, not this concrete type directly.
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository returns a UserRepository backed by the given GORM connection.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// FindByID returns a single user by their ID.
func (r *userRepository) FindByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Wishlist").First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail returns a single user by their email address.
func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpsertByGoogle(user *model.User) error {
	// OnConflict 意思是：
	// 如果 email 已存在（用戶之前登入過）→ 只更新 name 和 image
	// 如果 email 不存在（新用戶）→ 建立新 user
	result := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "image"}),
	}).Create(user)
	return result.Error
}
func (r *userRepository) Update(user *model.User) error {
	// Updates 只更新非零值的欄位
	// 例如只傳了 name，就只更新 name，不影響其他欄位
	result := r.db.Model(user).Updates(user)
	return result.Error
}

// AddWishItem inserts a new row into the wishlists join table.
func (r *userRepository) AddWishItem(userID string, productID string) error {
	wishlist := model.Wishlist{
		UserID:    userID,
		ProductID: productID,
	}
	return r.db.Create(&wishlist).Error
}

// RemoveWishItem deletes a row from the wishlists join table.
func (r *userRepository) RemoveWishItem(userID string, productID string) error {
	return r.db.Where("user_id = ? AND product_id = ?", userID, productID).
		Delete(&model.Wishlist{}).Error
}
