package repositories

import (
	"context"
	"errors"
	"time"

	"adhomes-backend/config"
	"adhomes-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeliveryRepository struct {
	collection *mongo.Collection
}

// Constructor
func NewDeliveryRepository() *DeliveryRepository {
	return &DeliveryRepository{
		collection: config.DB.Collection("deliveries"),
	}
}

// -----------------------------
// Create delivery
// -----------------------------
func (r *DeliveryRepository) Create(
	ctx context.Context,
	delivery models.Delivery,
) error {
	_, err := r.collection.InsertOne(ctx, delivery)
	return err
}

// -----------------------------
// Update delivery status
// -----------------------------
func (r *DeliveryRepository) UpdateStatus(
	ctx context.Context,
	deliveryID string,
	status string,
) (*models.Delivery, error) {

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

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": oid},
		update,
	)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, errors.New("delivery not found")
	}

	var delivery models.Delivery
	if err := r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&delivery); err != nil {
		return nil, err
	}

	return &delivery, nil
}

// -----------------------------
// Get delivery by order ID
// -----------------------------
func (r *DeliveryRepository) FindByOrderID(
	ctx context.Context,
	orderID string,
) (*models.Delivery, error) {

	var delivery models.Delivery
	err := r.collection.
		FindOne(ctx, bson.M{"order_id": orderID}).
		Decode(&delivery)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("delivery not found")
		}
		return nil, err
	}

	return &delivery, nil
}
