package services_impl

import (
	"adhomes-backend/models"
	"adhomes-backend/repositories"
)

type cartServiceImpl struct {
	cartRepo    *repositories.CartRepository
	productRepo *repositories.ProductRepository
}

func NewCartService(cartRepo *repositories.CartRepository, productRepo *repositories.ProductRepository) *cartServiceImpl {
	return &cartServiceImpl{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (s *cartServiceImpl) CreateCart(cart models.Cart) (models.Cart, error) {
	return s.cartRepo.CreateCart(cart)
}

func (s *cartServiceImpl) GetCartByUserID(userID string) (models.Cart, error) {
	return s.cartRepo.FindCartByUserID(userID)
}

func (s *cartServiceImpl) UpdateCart(id string, cart models.Cart) (models.Cart, error) {
	return s.cartRepo.UpdateCart(id, cart)
}

func (s *cartServiceImpl) DeleteCart(id string) error {
	return s.cartRepo.DeleteCart(id)
}
