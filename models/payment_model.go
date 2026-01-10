package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	OrderID   string             `json:"order_id" bson:"order_id"`
	UserID    string             `json:"user_id" bson:"user_id"`
	Amount    float64            `json:"amount" bson:"amount"`
	Method    string             `json:"method" bson:"method"`
	Status    string             `json:"status" bson:"status"`
	Email     string             `json:"email" bson:"email"`
	Reference string             `json:"reference" bson:"reference"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}
