package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type CartItem struct {
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
	Quantity  int                `json:"quantity" bson:"quantity"`
}

type Cart struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID string             `json:"user_id" bson:"user_id"`
	Items  []CartItem         `json:"items" bson:"items"`
}
