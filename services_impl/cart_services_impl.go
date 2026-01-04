package services_impl

import (
	"adhomes-backend/models"
	"adhomes-backend/repositories"
	"adhomes-backend/services"
)

type cartServiceImpl struct {
	cartRepo *repositories.CartRepository
}

func NewCartService() services.CartService {
	return &cartServiceImpl{
		cartRepo: repositories.NewCartRepository(),
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
