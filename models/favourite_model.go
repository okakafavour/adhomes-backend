package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Favourite struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	ProductID string             `bson:"product_id" json:"product_id"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
