package services

import "adhomes-backend/models"

type OrderService interface {
	CreateOrder(order models.Order) (models.Order, error)
	GetOrderByID(id string) (models.Order, error)
	GetOrdersByUserID(userID string) ([]models.Order, error)
	GetAllOrders() ([]models.Order, error)

	UpdateOrder(id string, order models.Order) (models.Order, error)
	UpdateOrderStatus(id string, status string) error

	ApproveOrder(id string) error
	CancelOrder(id string) error

	DeleteOrder(id string) error
}
