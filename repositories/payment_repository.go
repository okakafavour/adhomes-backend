package repositories

import (
	"adhomes-backend/config"
	"adhomes-backend/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentRepository struct {
	collection *mongo.Collection
}

func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{
		collection: config.GetCollection("payments"),
	}
}

func (r *PaymentRepository) Create(payment models.Payment) (models.Payment, error) {
	payment.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(context.Background(), payment)
	return payment, err
}
