package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	OrderID   string             `bson:"order_id" json:"order_id"`
	Amount    float64            `bson:"amount" json:"amount"`
	Reference string             `bson:"reference" json:"reference"`
	Gateway   string             `bson:"gateway" json:"gateway"`
	Status    string             `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
