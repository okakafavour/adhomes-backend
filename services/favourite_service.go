package services

import "adhomes-backend/models"

type FavouriteService interface {
	AddFavourite(userID, productID string) (*models.Favourite, error)
	GetFavourites(userID string) ([]*models.Favourite, error)
	RemoveFavourite(favouriteID string) error
}
