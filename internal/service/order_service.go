package service

import (
	"github.com/wingc34/mini-commerce-api/internal/config"
	"github.com/wingc34/mini-commerce-api/internal/model"
	"github.com/wingc34/mini-commerce-api/internal/repository"
)

type OrderService interface {
	CreateDraftOrder(draft *model.DraftOrder, items []model.DraftOrderItem) error
	GetOrders(userID string, page int) ([]model.Order, int64, error)
	GetOrderDetail(id string) (*model.Order, *model.DraftOrder, error)
}
type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) CreateDraftOrder(draft *model.DraftOrder, items []model.DraftOrderItem) error {
	return s.repo.CreateDraft(draft, items)
}

func (s *orderService) GetOrders(userID string, page int) ([]model.Order, int64, error) {
	return s.repo.FindByUserID(userID, page, config.OrderPageSize)
}

func (s *orderService) GetOrderDetail(id string) (*model.Order, *model.DraftOrder, error) {
	return s.repo.FindOrderOrDraftByID(id)
}
