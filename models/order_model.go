package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShippingAddress struct {
	Street     string `json:"street" bson:"street"`
	City       string `json:"city" bson:"city"`
	State      string `json:"state" bson:"state"`
	Country    string `json:"country" bson:"country"`
	PostalCode string `json:"postal_code" bson:"postal_code"`
}

type OrderItem struct {
	ProductID string `json:"product_id" bson:"product_id"`
	Quantity  int    `json:"quantity" bson:"quantity"`
}

type Order struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	CustomerName  string `json:"customer_name" bson:"customer_name"`
	CustomerEmail string `json:"customer_email" bson:"customer_email"`
	CustomerPhone string `json:"customer_phone" bson:"customer_phone"`

	DeliveryType    string          `json:"delivery_type" bson:"delivery_type"`
	ShippingAddress ShippingAddress `json:"shipping_address" bson:"shipping_address"`

	Items         []OrderItem `json:"items" bson:"items"`
	TotalAmount   float64     `json:"total_amount" bson:"total_amount"`
	Status        string      `json:"status" bson:"status"`
	PaymentStatus string      `bson:"payment_status" json:"payment_status"`

	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
