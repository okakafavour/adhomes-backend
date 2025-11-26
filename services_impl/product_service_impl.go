package services_impl

import (
	"context"
	"errors"
	"time"

	"adhomes-backend/config"
	"adhomes-backend/models"
	"adhomes-backend/services"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductServiceImpl struct {
	collection *mongo.Collection
}

func NewProductService() services.ProductService {
	return &ProductServiceImpl{
		collection: config.GetCollection("products"),
	}
}

// OPTION 1 — Add product by fields
func (s *ProductServiceImpl) AddProduct(name, description string, price float64, category, imageURL string) (models.Product, error) {
	product := models.Product{
		ID:          primitive.NewObjectID(),
		Name:        name,
		Description: description,
		Price:       price,
		Category:    category,
		ImageURL:    imageURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := s.collection.InsertOne(context.Background(), product)
	if err != nil {
		return models.Product{}, err
	}

	return product, nil
}

// Update Product
func (s *ProductServiceImpl) UpdateProduct(id string, product models.Product) (models.Product, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Product{}, errors.New("invalid product ID")
	}

	product.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"category":    product.Category,
			"image_url":   product.ImageURL,
			"updated_at":  product.UpdatedAt,
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

// Delete Product
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

// Get All Products
func (s *ProductServiceImpl) GetAllProducts() ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	for cursor.Next(ctx) {
		var prod models.Product
		if err := cursor.Decode(&prod); err != nil {
			return nil, err
		}
		products = append(products, prod)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
