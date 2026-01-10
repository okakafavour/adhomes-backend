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

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		collection: config.GetCollection("products"),
	}
}

func (r *ProductRepository) CreateProduct(product models.Product) (models.Product, error) {
	product.ID = primitive.NewObjectID()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(context.Background(), product)
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}

func (r *ProductRepository) UpdateProduct(id string, product models.Product) (models.Product, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Product{}, errors.New("Invalid product id")
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

	result, err := r.collection.UpdateOne(context.Background(), bson.M{"_id": oid}, update)
	if err != nil {
		return models.Product{}, err
	}
	if result.MatchedCount == 0 {
		return models.Product{}, errors.New("Product not found")
	}
	product.ID = oid
	return product, nil
}

func (r *ProductRepository) DeleteProduct(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("Invalid product id")
	}

	result, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("Product not found")
	}
	return nil
}

func (r *ProductRepository) FindAll() ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) FindByID(id string) (models.Product, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Product{}, errors.New("invalid product id")
	}

	var product models.Product
	err = r.collection.FindOne(context.Background(), bson.M{"_id": oid}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Product{}, errors.New("product not found")
		}
		return models.Product{}, err
	}

	return product, nil
}
