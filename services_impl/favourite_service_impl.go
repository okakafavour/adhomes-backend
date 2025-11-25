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

type FavouriteServiceImpl struct {
	FavoriteCollection *mongo.Collection
}

func NewFavoriteService() services.FavouriteService {
	return &FavouriteServiceImpl{
		FavoriteCollection: config.DB.Collection("favorites"),
	}
}

func (s *FavouriteServiceImpl) AddFavourite(userID, productID string) (*models.Favourite, error) {
	// Check for duplicate
	var existing models.Favourite
	err := s.FavoriteCollection.FindOne(context.Background(), bson.M{"user_id": userID, "product_id": productID}).Decode(&existing)
	if err == nil {
		return nil, errors.New("product already in favorites")
	}

	fav := models.Favourite{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		ProductID: productID,
		CreatedAt: time.Now(),
	}

	_, err = s.FavoriteCollection.InsertOne(context.Background(), fav)
	if err != nil {
		return nil, err
	}

	return &fav, nil
}

func (s *FavouriteServiceImpl) GetFavourites(userID string) ([]*models.Favourite, error) {
	cursor, err := s.FavoriteCollection.Find(context.Background(), bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var favorites []*models.Favourite
	for cursor.Next(context.Background()) {
		var fav models.Favourite
		cursor.Decode(&fav)
		favorites = append(favorites, &fav)
	}

	return favorites, nil
}

func (s *FavouriteServiceImpl) RemoveFavourite(favoriteID string) error {
	oid, err := primitive.ObjectIDFromHex(favoriteID)
	if err != nil {
		return errors.New("invalid favorite ID")
	}

	_, err = s.FavoriteCollection.DeleteOne(context.Background(), bson.M{"_id": oid})
	if err != nil {
		return err
	}

	return nil
}
