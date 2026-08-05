package service

import (
	"encoding/json"

	"github.com/wingc34/mini-commerce-api/internal/config"
	"github.com/wingc34/mini-commerce-api/internal/model"
	"github.com/wingc34/mini-commerce-api/internal/repository"
)

type ProductService interface {
	GetProducts(page int) ([]model.Product, int64, error)
	GetRecommended() ([]model.Product, error)
	GetProductByID(id string) (*model.Product, error)
	CheckStock(productID string, attributes map[string]string) (bool, error)
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetProducts(page int) ([]model.Product, int64, error) {
	return s.repo.FindAll(page, config.ProductPageSize)
}

func (s *productService) GetRecommended() ([]model.Product, error) {
	products, err := s.repo.FindRecommended(config.RecommendedLimit)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *productService) GetProductByID(id string) (*model.Product, error) {
	// 呼叫 repo
	return s.repo.FindByID(id)
}

func (s *productService) CheckStock(productID string, attributes map[string]string) (bool, error) {
	// 1. 把 attributes 轉成 JSON string
	// 2. 呼叫 repo.CheckStock
	attrsJSON, err := json.Marshal(attributes)
	if err != nil {
		return false, err
	}
	// 第二步：呼叫 repo
	return s.repo.CheckStock(productID, string(attrsJSON))
}
