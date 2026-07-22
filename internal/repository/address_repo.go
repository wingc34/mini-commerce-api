package repository

import (
	"github.com/wingc34/mini-commerce-api/internal/model"
	"gorm.io/gorm"
)

// AddressRepository defines the database operations for addresses.
// Defining this as an interface (not a concrete struct) lets us inject
// a mock implementation in tests without touching a real database.
type AddressRepository interface {
	FindAllByUserID(userID string) ([]model.Address, error)
	FindByID(id string) (*model.Address, error)
	Create(address *model.Address) error
	Update(address *model.Address) error
	SoftDelete(id string, userID string) error
	SetDefault(id string, userID string) error
	ClearDefault(userID string) error
}

// addressRepository is the GORM implementation of AddressRepository.
// Lowercase name keeps it private - external packages can only use
// the AddressRepository interface, not this concrete type directly.
type addressRepository struct {
	db *gorm.DB
}

// NewAddressRepository returns a AddressRepository backed by the given GORM connection.
func NewAddressRepository(db *gorm.DB) AddressRepository {
	return &addressRepository{db: db}
}

// FindAllByUserID returns all non-deleted addresses belonging to the given user.
func (r *addressRepository) FindAllByUserID(userID string) ([]model.Address, error) {
	var addresses []model.Address
	err := r.db.Where("user_id = ?", userID).Find(&addresses).Error
	if err != nil {
		return []model.Address{}, err
	}
	return addresses, nil
}

// FindByID returns a single address by their ID.
func (r *addressRepository) FindByID(id string) (*model.Address, error) {
	var address model.Address
	err := r.db.First(&address, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

// Create inserts a new address into the database.
func (r *addressRepository) Create(address *model.Address) error {
	return r.db.Create(address).Error
}

func (r *addressRepository) Update(address *model.Address) error {
	// Updates 只更新非零值的欄位
	// 例如只傳了 name，就只更新 name，不影響其他欄位
	result := r.db.Model(address).Updates(address)
	return result.Error
}

func (r *addressRepository) SoftDelete(id string, userID string) error {
	// GORM 看到 gorm.DeletedAt 欄位，自動把 Delete 變成
	// UPDATE addresses SET deleted_at = NOW() WHERE id = ? AND user_id = ?
	// 不會真的刪除資料
	return r.db.Where("id = ? AND user_id = ?", id, userID).
		Delete(&model.Address{}).Error
}

func (r *addressRepository) SetDefault(id string, userID string) error {
	return r.db.Model(&model.Address{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_default", true).Error
}

func (r *addressRepository) ClearDefault(userID string) error {
	// 把這個 user 所有地址的 is_default 設為 false
	return r.db.Model(&model.Address{}).
		Where("user_id = ?", userID).
		Update("is_default", false).Error
}
