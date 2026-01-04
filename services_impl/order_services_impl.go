package services_impl

import (
	"adhomes-backend/models"
	"adhomes-backend/repositories"
	"adhomes-backend/services"
	"adhomes-backend/utils"
	"errors"
)

type orderServiceImpl struct {
	orderRepo *repositories.OrderRepository
}

func NewOrderService() services.OrderService {
	return &orderServiceImpl{
		orderRepo: repositories.NewOrderRepository(),
	}
}

func (s *orderServiceImpl) CreateOrder(order models.Order) (models.Order, error) {
	return s.orderRepo.CreateOrder(order)
}

func (s *orderServiceImpl) GetOrderByID(id string) (models.Order, error) {
	return s.orderRepo.FindOrderByID(id)
}

func (s *orderServiceImpl) GetOrdersByUserID(userID string) ([]models.Order, error) {
	return s.orderRepo.FindOrdersByUserID(userID)
}

func (s *orderServiceImpl) GetAllOrders() ([]models.Order, error) {
	return s.orderRepo.FindAll()
}

func (s *orderServiceImpl) UpdateOrder(id string, order models.Order) (models.Order, error) {
	return s.orderRepo.UpdateOrder(id, order)
}

func (s *orderServiceImpl) UpdateOrderStatus(id string, newStatus string) error {
	if !utils.IsValidStatus(newStatus) {
		return errors.New("invalid order status")
	}
	return s.orderRepo.UpdateOrderStatus(id, newStatus)
}

func (s *orderServiceImpl) DeleteOrder(id string) error {
	return s.orderRepo.DeleteOrder(id)
}
