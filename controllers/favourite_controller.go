package controllers

import (
	"net/http"

	"adhomes-backend/services"
	"adhomes-backend/services_impl"

	"github.com/gin-gonic/gin"
)

type FavouriteController struct {
	favoriteService services.FavouriteService
}

var favouriteController *FavouriteController

func FavouriteControllerSingleton() *FavouriteController {
	if adminController == nil {
		InitFavouriteController()
	}
	return favouriteController
}

func NewFavoriteController(service services.FavouriteService) *FavouriteController {
	return &FavouriteController{
		favoriteService: service,
	}
}

func InitFavouriteController() {
	favService := services_impl.NewFavoriteService()
	favouriteController = NewFavoriteController(favService)
}

func (fc *FavouriteController) AddFavorite(ctx *gin.Context) {
	var req struct {
		UserID    string `json:"user_id"`
		ProductID string `json:"product_id"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fav, err := fc.favoriteService.AddFavourite(req.UserID, req.ProductID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":  "added to favorites",
		"favorite": fav,
	})
}

func (fc *FavouriteController) GetFavorites(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	favorites, err := fc.favoriteService.GetFavourites(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"favorites": favorites,
	})
}

func (fc *FavouriteController) RemoveFavorite(ctx *gin.Context) {
	favoriteID := ctx.Param("id")
	err := fc.favoriteService.RemoveFavourite(favoriteID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "removed from favorites",
	})
}
