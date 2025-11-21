package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wallet struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    string             `bson:"user_id" json:"user_id"`
	Balance   float64            `bson:"balance" json:"balance"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
