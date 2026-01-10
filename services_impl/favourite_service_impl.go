package services_impl

import (
	"errors"

	"adhomes-backend/models"
	"adhomes-backend/repositories"
	"adhomes-backend/services"
)

type FavouriteServiceImpl struct {
	favRepo *repositories.FavouriteRepository
}

func NewFavoriteService() services.FavouriteService {
	return &FavouriteServiceImpl{
		favRepo: repositories.NewFavouriteRepository(),
	}
}

func (s *FavouriteServiceImpl) AddFavourite(userID, productID string) (*models.Favourite, error) {
	exists, err := s.favRepo.Exits(userID, productID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("product already in favorites")
	}

	return s.favRepo.CreateFavourite(userID, productID)
}

func (s *FavouriteServiceImpl) GetFavourites(userID string) ([]*models.Favourite, error) {
	return s.favRepo.FindByUserID(userID)
}

func (s *FavouriteServiceImpl) RemoveFavourite(favouriteID string) error {
	return s.favRepo.DeleteByID(favouriteID)
}
