package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CartID          string             `bson:"cart_id" json:"cart_id"`
	UserID          string             `bson:"user_id" json:"user_id"`
	ProductIDs      []string           `bson:"product_ids" json:"product_ids"`
	Amount          float64            `bson:"amount" json:"amount"`
	DeliveryAddress string             `bson:"delivery_address" json:"delivery_address"`
	PaymentStatus   string             `bson:"payment_status" json:"payment_status"`
	OrderStatus     string             `bson:"order_status" json:"status"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}
