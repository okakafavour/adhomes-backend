package services_impl

import (
	"context"
	"time"

	"adhomes-backend/models"
	"adhomes-backend/repositories"
	"adhomes-backend/services"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeliveryServiceImpl struct {
	deliveryRepo *repositories.DeliveryRepository
}

// Constructor
func NewDeliveryService() services.DeliveryService {
	return &DeliveryServiceImpl{
		deliveryRepo: repositories.NewDeliveryRepository(),
	}
}

// -----------------------------
// Assign rider to order
// -----------------------------
func (s *DeliveryServiceImpl) AssignRider(
	orderID,
	riderID string,
) (*models.Delivery, error) {

	delivery := models.Delivery{
		ID:        primitive.NewObjectID(),
		OrderID:   orderID,
		RiderID:   riderID,
		Status:    "Assigned",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.deliveryRepo.Create(
		context.Background(),
		delivery,
	); err != nil {
		return nil, err
	}

	return &delivery, nil
}

// -----------------------------
// Update delivery status
// -----------------------------
func (s *DeliveryServiceImpl) UpdateDeliveryStatus(
	deliveryID,
	status string,
) (*models.Delivery, error) {

	return s.deliveryRepo.UpdateStatus(
		context.Background(),
		deliveryID,
		status,
	)
}

// -----------------------------
// Get delivery by order ID
// -----------------------------
func (s *DeliveryServiceImpl) GetDeliveryByOrder(
	orderID string,
) (*models.Delivery, error) {

	return s.deliveryRepo.FindByOrderID(
		context.Background(),
		orderID,
	)
}
