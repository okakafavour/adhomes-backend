package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ProductID string `bson:"product_id" json:"product_id"`
	Quantity  int    `bson:"quantity" json:"quantity"`
}

type Order struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          string             `bson:"user_id" json:"user_id"`
	Items           []OrderItem        `bson:"items" json:"items"`
	Total           float64            `bson:"total" json:"total"`
	DeliveryAddress string             `bson:"delivery_address" json:"delivery_address"`
	PaymentStatus   string             `bson:"payment_status" json:"payment_status"`
	OrderStatus     string             `bson:"order_status" json:"status"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}
