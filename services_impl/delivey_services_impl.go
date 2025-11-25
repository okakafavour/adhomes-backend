package services_impl

import (
	"context"
	"errors"
	"time"

	"adhomes-backend/config"
	"adhomes-backend/models"
	"adhomes-backend/services"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeliveryServiceImpl struct {
	DeliveryCollection *mongo.Collection
}

func NewDeliveryService() services.DeliveryService {
	return &DeliveryServiceImpl{
		DeliveryCollection: config.DB.Collection("deliveries"),
	}
}

// Assign a rider to an order
func (s *DeliveryServiceImpl) AssignRider(orderID, riderID string) (*models.Delivery, error) {
	delivery := models.Delivery{
		ID:        primitive.NewObjectID(),
		OrderID:   orderID,
		RiderID:   riderID,
		Status:    "Assigned",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := s.DeliveryCollection.InsertOne(context.Background(), delivery)
	if err != nil {
		return nil, err
	}

	return &delivery, nil
}

// Update delivery status
func (s *DeliveryServiceImpl) UpdateDeliveryStatus(deliveryID, status string) (*models.Delivery, error) {
	oid, err := primitive.ObjectIDFromHex(deliveryID)
	if err != nil {
		return nil, errors.New("invalid delivery ID")
	}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	_, err = s.DeliveryCollection.UpdateOne(context.Background(), bson.M{"_id": oid}, update)
	if err != nil {
		return nil, err
	}

	var delivery models.Delivery
	err = s.DeliveryCollection.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&delivery)
	if err != nil {
		return nil, err
	}

	return &delivery, nil
}

// Get delivery by order ID
func (s *DeliveryServiceImpl) GetDeliveryByOrder(orderID string) (*models.Delivery, error) {
	var delivery models.Delivery
	err := s.DeliveryCollection.FindOne(context.Background(), bson.M{"order_id": orderID}).Decode(&delivery)
	if err != nil {
		return nil, err
	}
	return &delivery, nil
}
