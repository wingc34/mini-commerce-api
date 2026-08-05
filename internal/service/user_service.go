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
	// 查現有地址數量
	addresses, err := s.addressRepo.FindAllByUserID(address.UserID)
	if err != nil {
		return err
	}

	// 第一個地址自動設為 default
	if len(addresses) == 0 {
		address.IsDefault = true
	}

	// 如果要設為 default，先清除舊的
	if address.IsDefault {
		if err := s.addressRepo.ClearDefault(address.UserID); err != nil {
			return err
		}
	}

	return s.addressRepo.Create(address)
}

func (s *userService) UpdateAddress(address *model.Address) error {
	if address.IsDefault {
		if err := s.addressRepo.ClearDefault(address.UserID); err != nil {
			return err
		}
	}
	return s.addressRepo.Update(address)
}

func (s *userService) DeleteAddress(id string, userID string) error {
	address, err := s.addressRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Soft delete
	if err := s.addressRepo.SoftDelete(id, userID); err != nil {
		return err
	}

	// 如果刪的是 default，提升另一個為 default
	if address.IsDefault {
		addresses, err := s.addressRepo.FindAllByUserID(userID)
		if err != nil {
			return err
		}
		if len(addresses) > 0 {
			return s.addressRepo.SetDefault(addresses[0].ID, userID)
		}
	}

	return nil
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
