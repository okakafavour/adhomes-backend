package repositories

import (
	"context"
	"errors"
	"time"

	"adhomes-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(collection *mongo.Collection) *ProductRepository {
	return &ProductRepository{collection}
}

// CREATE
func (r *ProductRepository) Create(product models.Product) (models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, product)
	return product, err
}

// UPDATE (PARTIAL)
func (r *ProductRepository) UpdateFields(
	id primitive.ObjectID,
	update bson.M,
) (*models.Product, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if len(update) == 0 {
		return nil, errors.New("no fields to update")
	}

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": update},
	)
	if err != nil {
		return nil, err
	}

	var product models.Product
	err = r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	return &product, err
}

// DELETE
func (r *ProductRepository) Delete(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// READ
func (r *ProductRepository) FindAll() ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	err = cursor.All(ctx, &products)
	return products, err
}

func (r *ProductRepository) FindByID(id primitive.ObjectID) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var product models.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	return &product, err
}
