package repository

import (
	"errors"

	"github.com/wingc34/mini-commerce-api/internal/model"
	"gorm.io/gorm"
)

// OrderRepository defines the database operations for orders.
// Defining this as an interface (not a concrete struct) lets us inject
// a mock implementation in tests without touching a real database.
type OrderRepository interface {
	// Draft Order
	CreateDraft(draft *model.DraftOrder, items []model.DraftOrderItem) error
	FindDraftByID(id string) (*model.DraftOrder, error)
	FindDraftItemsByDraftID(draftID string) ([]model.DraftOrderItem, error)
	UpdateDraftStatus(id string, status model.DraftStatus) error

	// Order
	CreateOrder(order *model.Order, items []model.OrderItem) error
	FindByID(id string) (*model.Order, error)
	FindByUserID(userID string, page int, pageSize int) ([]model.Order, int64, error)

	// Dual lookup（先查 orders，找不到再查 draft_orders）
	FindOrderOrDraftByID(id string) (*model.Order, *model.DraftOrder, error)
}

// orderRepository is the GORM implementation of OrderRepository.
// Lowercase name keeps it private - external packages can only use
// the OrderRepository interface, not this concrete type directly.
type orderRepository struct {
	db *gorm.DB
}

// NewOrderRepository returns a OrderRepository backed by the given GORM connection.
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateDraft(draft *model.DraftOrder, items []model.DraftOrderItem) error {
	// 提示：GORM transaction 的寫法
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 在這裡的所有操作都在同一個 transaction 裡
		if err := tx.Create(draft).Error; err != nil {
			return err // 回傳 error 會自動 rollback
		}
		for _, item := range items {
			if err := tx.Create(&item).Error; err != nil {
				return err // 回傳 error 會自動 rollback
			}
		}
		return nil // 回傳 nil 會自動 commit
	})
}

// FindByID returns a single draftorder by their ID.
func (r *orderRepository) FindDraftByID(id string) (*model.DraftOrder, error) {
	var draftorder model.DraftOrder
	err := r.db.First(&draftorder, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &draftorder, nil
}

func (r *orderRepository) UpdateDraftStatus(id string, status model.DraftStatus) error {
	return r.db.Model(&model.DraftOrder{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *orderRepository) FindDraftItemsByDraftID(draftID string) ([]model.DraftOrderItem, error) {
	var items []model.DraftOrderItem
	err := r.db.Where("draft_order_id = ?", draftID).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

// CreateOrder 只負責 INSERT，不管業務邏輯
func (r *orderRepository) CreateOrder(order *model.Order, items []model.OrderItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		for _, item := range items {
			item.OrderID = order.ID
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// FindByUserID returns all non-deleted orders belonging to the given user.
func (r *orderRepository) FindByUserID(userID string, page int, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	offset := (page - 1) * pageSize

	// 先查總數，讓前端知道總共有幾頁
	if err := r.db.Model(&model.Order{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 再查當頁的資料
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Preload("ShippingAddress").
		Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// FindByID returns a single order by their ID.
func (r *orderRepository) FindByID(id string) (*model.Order, error) {
	var order model.Order
	err := r.db.First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindOrderOrDraftByID(id string) (*model.Order, *model.DraftOrder, error) {
	// 第一步：先查 orders 表
	var order model.Order
	err := r.db.Preload("ShippingAddress").
		Preload("OrderItems").
		First(&order, "id = ?", id).Error

	// 找到了，直接回傳
	if err == nil {
		return &order, nil, nil
	}

	// 如果是其他錯誤（不是找不到），回傳錯誤
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, err
	}

	// 第二步：order 找不到，改查 draft_orders
	var draft model.DraftOrder
	err = r.db.Preload("ShippingAddress").
		First(&draft, "id = ?", id).Error
	if err != nil {
		return nil, nil, err
	}

	return nil, &draft, nil
}
