package services

import "adhomes-backend/models"

type DeliveryService interface {
	AssignRider(orderID, riderID string) (*models.Delivery, error)
	UpdateDeliveryStatus(deliveryID, status string) (*models.Delivery, error)
	GetDeliveryByOrder(orderID string) (*models.Delivery, error)
}
