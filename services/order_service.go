package services

import "adhomes-backend/models"

type OrderService interface {
	CreateOrder(order models.Order) (models.Order, error)
	GetOrderByID(id string) (models.Order, error)
	GetOrdersByUserID(userID string) ([]models.Order, error)
	DeleteOrder(id string) error
	UpdateOrder(id string, order models.Order) (models.Order, error)
	UpdateOrderStatus(orderID string, newStatus string) error

	GetAllOrders() ([]models.Order, error)
}
