package services

import "adhomes-backend/models"

// OrderService defines all operations related to Orders
type OrderService interface {
	CreateOrder(order models.Order) (models.Order, error)
	GetOrderByID(id string) (models.Order, error)
	GetOrdersByUserID(userID string) ([]models.Order, error)
	DeleteOrder(id string) error
}
