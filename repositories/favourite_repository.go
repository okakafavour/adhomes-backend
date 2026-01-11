package repositories

import (
	"adhomes-backend/models"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FavouriteRepository struct {
	collection *mongo.Collection
}

func NewFavouriteRepository(collection *mongo.Collection) *FavouriteRepository {
	return &FavouriteRepository{collection}
}

func (r *FavouriteRepository) Exits(userID, productID string) (bool, error) {
	err := r.collection.FindOne(context.Background(), bson.M{"user_id": userID, "product_id": productID}).Err()
	if err == nil {
		return true, nil
	}

	if err == mongo.ErrNoDocuments {
		return false, nil
	}

	return false, err
}

func (r *FavouriteRepository) CreateFavourite(userID, productID string) (*models.Favourite, error) {
	fav := models.Favourite{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		ProductID: productID,
		CreatedAt: time.Now(),
	}

	_, err := r.collection.InsertOne(context.Background(), fav)
	if err != nil {
		return nil, err
	}

	return &fav, nil
}

func (r *FavouriteRepository) FindByUserID(userID string) ([]*models.Favourite, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var favourites []*models.Favourite
	for cursor.Next(context.Background()) {
		var fav models.Favourite
		if err := cursor.Decode(&fav); err != nil {
			return nil, err
		}
		favourites = append(favourites, &fav)
	}
	return favourites, nil
}

func (r *FavouriteRepository) DeleteByID(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid favorite ID")
	}

	_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": oid})
	return err
}
