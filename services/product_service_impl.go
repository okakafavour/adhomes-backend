package services

import (
	"context"
	"errors"
	"time"

	"adhomes-backend/config"
	"adhomes-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductServiceImpl struct {
	collection *mongo.Collection
}

func NewProductService() ProductService {
	return &ProductServiceImpl{
		collection: config.GetCollection("products"),
	}
}

func (s *ProductServiceImpl) CreateProduct(product models.Product) (models.Product, error) {
	product.ID = primitive.NewObjectID()
	product.CreatedAt = time.Now()

	_, err := s.collection.InsertOne(context.Background(), product)
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}

func (s *ProductServiceImpl) UpdateProduct(id string, product models.Product) (models.Product, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Product{}, errors.New("invalid product ID")
	}

	update := bson.M{
		"$set": bson.M{
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"category":    product.Category,
			"image_url":   product.ImageURL,
		},
	}

	result, err := s.collection.UpdateOne(context.Background(), bson.M{"_id": oid}, update)
	if err != nil {
		return models.Product{}, err
	}
	if result.MatchedCount == 0 {
		return models.Product{}, errors.New("product not found")
	}

	product.ID = oid
	return product, nil
}

func (s *ProductServiceImpl) DeleteProduct(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid product ID")
	}

	result, err := s.collection.DeleteOne(context.Background(), bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("product not found")
	}
	return nil
}
