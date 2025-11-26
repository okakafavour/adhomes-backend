package services_impl

import (
	"adhomes-backend/config"
	"adhomes-backend/models"
	"adhomes-backend/services"
	"adhomes-backend/utils"

	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type orderServiceImpl struct {
	collection *mongo.Collection
}

func NewOrderService() services.OrderService {
	return &orderServiceImpl{
		collection: config.DB.Collection("orders"),
	}
}

func (s *orderServiceImpl) CreateOrder(order models.Order) (models.Order, error) {
	order.ID = primitive.NewObjectID()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	_, err := s.collection.InsertOne(context.Background(), order)
	return order, err
}

func (s *orderServiceImpl) GetOrderByID(id string) (models.Order, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Order{}, errors.New("invalid order ID")
	}

	var order models.Order
	err = s.collection.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Order{}, errors.New("order not found")
		}
		return models.Order{}, err
	}
	return order, nil
}

func (s *orderServiceImpl) GetOrdersByUserID(userID string) ([]models.Order, error) {
	cursor, err := s.collection.Find(context.Background(), bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var orders []models.Order
	for cursor.Next(context.Background()) {
		var order models.Order
		if err := cursor.Decode(&order); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (s *orderServiceImpl) DeleteOrder(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid order ID")
	}

	result, err := s.collection.DeleteOne(context.Background(), bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("order not found")
	}
	return nil
}

func (s *orderServiceImpl) UpdateOrder(id string, order models.Order) (models.Order, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Order{}, errors.New("invalid order ID")
	}

	order.UpdatedAt = time.Now()
	update := bson.M{
		"$set": bson.M{
			"payment_status":   order.PaymentStatus,
			"order_status":     order.OrderStatus,
			"delivery_address": order.DeliveryAddress,
			"items":            order.Items,
			"total":            order.Total,
			"updated_at":       order.UpdatedAt,
		},
	}

	result, err := s.collection.UpdateOne(context.Background(), bson.M{"_id": oid}, update)
	if err != nil {
		return models.Order{}, err
	}
	if result.MatchedCount == 0 {
		return models.Order{}, errors.New("order not found")
	}

	order.ID = oid
	return order, nil
}

func (s *orderServiceImpl) UpdateOrderStatus(orderID string, newStatus string) error {
	// Validate ObjectID
	objID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return errors.New("invalid order ID")
	}

	// Validate status
	if !utils.IsValidStatus(newStatus) {
		return errors.New("invalid order status")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"status":     newStatus,
			"updated_at": time.Now(),
		},
	}

	result, err := s.collection.UpdateByID(ctx, objID, update)
	if err != nil {
		return errors.New("failed to update order status")
	}

	if result.MatchedCount == 0 {
		return errors.New("order not found")
	}

	return nil
}

func (s *orderServiceImpl) GetAllOrders() ([]models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}
