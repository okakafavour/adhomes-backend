package services_impl

import (
	"errors"
	"time"

	"adhomes-backend/models"
	"adhomes-backend/repositories"
	"adhomes-backend/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type orderServiceImpl struct {
	orderRepo   *repositories.OrderRepository
	productRepo *repositories.ProductRepository
}

func NewOrderService(orderRepo *repositories.OrderRepository, productRepo *repositories.ProductRepository) *orderServiceImpl {
	return &orderServiceImpl{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

// -----------------------------
// CREATE ORDER
// -----------------------------
func (s *orderServiceImpl) CreateOrder(order models.Order) (models.Order, error) {

	if len(order.Items) == 0 {
		return models.Order{}, errors.New("order must contain at least one item")
	}

	var total float64

	for _, item := range order.Items {
		productID, err := primitive.ObjectIDFromHex(item.ProductID)
		if err != nil {
			return models.Order{}, errors.New("invalid product ID")
		}

		product, err := s.productRepo.FindByID(productID)
		if err != nil {
			return models.Order{}, err
		}

		total += product.Price * float64(item.Quantity)
	}

	order.ID = primitive.NewObjectID()
	order.TotalAmount = total
	order.PaymentStatus = "unpaid"
	order.Status = "pending"
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	return s.orderRepo.CreateOrder(order)
}

// -----------------------------
// APPROVE ORDER (ADMIN)
// -----------------------------
func (s *orderServiceImpl) ApproveOrder(id string) error {
	return s.orderRepo.UpdateOrderStatus(id, "approved")
}

// -----------------------------
// CANCEL ORDER
// -----------------------------
func (s *orderServiceImpl) CancelOrder(id string) error {
	return s.orderRepo.UpdateOrderStatus(id, "cancelled")
}

// -----------------------------
// GET ORDER BY ID
// -----------------------------
func (s *orderServiceImpl) GetOrderByID(id string) (models.Order, error) {
	return s.orderRepo.FindOrderByID(id)
}

// -----------------------------
// GET ORDERS BY USER
// -----------------------------
func (s *orderServiceImpl) GetOrdersByUserID(userID string) ([]models.Order, error) {
	return s.orderRepo.FindOrdersByUserID(userID)
}

// -----------------------------
// GET ALL ORDERS (ADMIN)
// -----------------------------
func (s *orderServiceImpl) GetAllOrders() ([]models.Order, error) {
	return s.orderRepo.FindAll()
}

// -----------------------------
// UPDATE ORDER (FULL UPDATE)
// -----------------------------
func (s *orderServiceImpl) UpdateOrder(id string, order models.Order) (models.Order, error) {
	order.UpdatedAt = time.Now()
	return s.orderRepo.UpdateOrder(id, order)
}

// -----------------------------
// UPDATE ORDER STATUS ONLY
// -----------------------------
func (s *orderServiceImpl) UpdateOrderStatus(id string, newStatus string) error {

	if !utils.IsValidStatus(newStatus) {
		return errors.New("invalid order status")
	}

	return s.orderRepo.UpdateOrderStatus(id, newStatus)
}

// -----------------------------
// DELETE ORDER
// -----------------------------
func (s *orderServiceImpl) DeleteOrder(id string) error {
	return s.orderRepo.DeleteOrder(id)
}
