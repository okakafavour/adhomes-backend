package repositories

import (
	"context"
	"errors"
	"time"

	"adhomes-backend/config"
	"adhomes-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentRepository struct {
	paymentCollection *mongo.Collection
	orderCollection   *mongo.Collection
}

func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{
		paymentCollection: config.DB.Collection("payments"),
		orderCollection:   config.DB.Collection("orders"),
	}
}

// ------------------------
// Create payment
// ------------------------
func (r *PaymentRepository) Create(
	ctx context.Context,
	payment models.Payment,
) error {
	_, err := r.paymentCollection.InsertOne(ctx, payment)
	return err
}

// ------------------------
// Find payment by reference
// ------------------------
func (r *PaymentRepository) FindByReference(
	ctx context.Context,
	reference string,
) (*models.Payment, error) {

	var payment models.Payment
	err := r.paymentCollection.
		FindOne(ctx, bson.M{"reference": reference}).
		Decode(&payment)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}

	return &payment, nil
}

// ------------------------
// Update payment status
// ------------------------
func (r *PaymentRepository) UpdateStatus(
	ctx context.Context,
	reference string,
	status string,
) error {

	_, err := r.paymentCollection.UpdateOne(
		ctx,
		bson.M{"reference": reference},
		bson.M{
			"$set": bson.M{
				"status":     status,
				"updated_at": time.Now(),
			},
		},
	)

	return err
}

// ------------------------
// Update order status after payment
// ------------------------
func (r *PaymentRepository) UpdateOrderStatus(
	ctx context.Context,
	orderID interface{},
	status string,
) error {

	_, err := r.orderCollection.UpdateOne(
		ctx,
		bson.M{"_id": orderID},
		bson.M{"$set": bson.M{"status": status}},
	)

	return err
}
