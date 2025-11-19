package services

import "adhomes-backend/models"

type CartService interface {
	CreateCart(cart models.Cart) (models.Cart, error)
	GetCartByUserID(userID string) (models.Cart, error)
	UpdateCart(id string, cart models.Cart) (models.Cart, error)
	DeleteCart(id string) error
}
