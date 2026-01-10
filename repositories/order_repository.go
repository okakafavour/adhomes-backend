package repositories

import (
	"adhomes-backend/config"
	"adhomes-backend/models"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepository struct {
	collection *mongo.Collection
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		collection: config.GetCollection("orders"),
	}
}

func (r *OrderRepository) CreateOrder(order models.Order) (models.Order, error) {
	order.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(context.Background(), order)
	return order, err
}

func (r *OrderRepository) FindOrderByID(id string) (models.Order, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Order{}, errors.New("Invalid order id")
	}

	var order models.Order
	err = r.collection.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Order{}, errors.New("Order not found")
		}
		return models.Order{}, err
	}
	return order, nil
}

func (r *OrderRepository) FindOrdersByUserID(userID string) ([]models.Order, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var orders []models.Order
	if err := cursor.All(context.Background(), &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) FindAll() ([]models.Order, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var orders []models.Order
	if err := cursor.All(context.Background(), &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) UpdateOrder(id string, order models.Order) (models.Order, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Order{}, errors.New("invalid order id")
	}

	update := bson.M{
		"$set": bson.M{
			"customer_name":    order.CustomerName,
			"customer_email":   order.CustomerEmail,
			"customer_phone":   order.CustomerPhone,
			"delivery_type":    order.DeliveryType,
			"shipping_address": order.ShippingAddress,
			"items":            order.Items,
			"total_amount":     order.TotalAmount,
			"status":           order.Status,
		},
	}

	result, err := r.collection.UpdateOne(context.Background(), bson.M{"_id": oid}, update)
	if err != nil {
		return models.Order{}, err
	}
	if result.MatchedCount == 0 {
		return models.Order{}, errors.New("order not found")
	}

	order.ID = oid
	return order, nil
}

func (r *OrderRepository) UpdateOrderStatus(id string, status string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("Invalid order id")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateByID(ctx, oid, update)
	if err != nil {
		return errors.New("failed to update order status")
	}
	if result.MatchedCount == 0 {
		return errors.New("Order not found")
	}

	return nil
}

func (r *OrderRepository) DeleteOrder(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("Invalid order id")
	}

	result, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("Order not found")
	}
	return nil
}
