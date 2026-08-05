package service

import (
	"github.com/wingc34/mini-commerce-api/internal/model"
	"github.com/wingc34/mini-commerce-api/internal/repository"
)

type UserService interface {
	GetMe(userID string) (*model.User, error)
	UpdateMe(userID string, name *string, phoneNumber *string) error

	// Address
	GetAddresses(userID string) ([]model.Address, error)
	CreateAddress(address *model.Address) error
	UpdateAddress(address *model.Address) error
	DeleteAddress(id string, userID string) error
	SetDefaultAddress(id string, userID string) error

	// Wishlist
	AddWishItem(userID string, productID string) error
	RemoveWishItem(userID string, productID string) error
}

type userService struct {
	userRepo    repository.UserRepository
	addressRepo repository.AddressRepository
}

func NewUserService(
	userRepo repository.UserRepository,
	addressRepo repository.AddressRepository,
) UserService {
	return &userService{
		userRepo:    userRepo,
		addressRepo: addressRepo,
	}
}

func (s *userService) GetMe(userID string) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *userService) UpdateMe(userID string, name *string, phoneNumber *string) error {
	return s.userRepo.Update(&model.User{
		ID:          userID,
		Name:        name,
		PhoneNumber: phoneNumber,
	})
}

func (s *userService) GetAddresses(userID string) ([]model.Address, error) {
	return s.addressRepo.FindAllByUserID(userID)
}

func (s *userService) CreateAddress(address *model.Address) error {
	return s.addressRepo.Create(address)
}

func (s *userService) UpdateAddress(address *model.Address) error {
	return s.addressRepo.Update(address)
}

func (s *userService) DeleteAddress(id string, userID string) error {
	return s.addressRepo.SoftDelete(id, userID)
}

func (s *userService) SetDefaultAddress(id string, userID string) error {
	if err := s.addressRepo.ClearDefault(userID); err != nil {
		return err
	}
	return s.addressRepo.SetDefault(id, userID)
}

func (s *userService) AddWishItem(userID string, productID string) error {
	return s.userRepo.AddWishItem(userID, productID)
}

func (s *userService) RemoveWishItem(userID string, productID string) error {
	return s.userRepo.RemoveWishItem(userID, productID)
}
