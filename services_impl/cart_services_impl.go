package services_impl

import (
	"adhomes-backend/config"
	"adhomes-backend/models"
	"adhomes-backend/services"

	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type cartServiceImpl struct {
	collection *mongo.Collection
}

func NewCartService() services.CartService {
	return &cartServiceImpl{
		collection: config.DB.Collection("carts"),
	}
}

func (s *cartServiceImpl) CreateCart(cart models.Cart) (models.Cart, error) {
	cart.ID = primitive.NewObjectID()
	cart.CreatedAt = time.Now()
	cart.UpdatedAt = time.Now()

	_, err := s.collection.InsertOne(context.Background(), cart)
	return cart, err
}

func (s *cartServiceImpl) GetCartByUserID(userID string) (models.Cart, error) {
	var cart models.Cart
	err := s.collection.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&cart)
	return cart, err
}

func (s *cartServiceImpl) UpdateCart(id string, cart models.Cart) (models.Cart, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Cart{}, errors.New("invalid cart ID")
	}

	cart.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"user_id":    cart.UserID,
			"items":      cart.Items,
			"updated_at": cart.UpdatedAt,
		},
	}

	result, err := s.collection.UpdateOne(context.Background(), bson.M{"_id": oid}, update)
	if err != nil {
		return models.Cart{}, err
	}
	if result.MatchedCount == 0 {
		return models.Cart{}, errors.New("cart not found")
	}

	cart.ID = oid
	return cart, nil
}

func (s *cartServiceImpl) DeleteCart(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid cart ID")
	}

	result, err := s.collection.DeleteOne(context.Background(), bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("cart not found")
	}
	return nil
}
