package controllers

import (
	"adhomes-backend/config"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// -------------------------------
// Clean favorites collection
// -------------------------------
func cleanFavoritesCollection() {
	if config.DB == nil {
		panic("Database not initialized! Did you run TestMain?")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := config.DB.Collection("favorites").Drop(ctx)
	if err != nil {
		panic(err)
	}
}

// -------------------------------
// Setup router for favorites
// -------------------------------
func setUpFavoriteRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Initialize Favorites controller
	InitFavouriteController()

	// Routes
	router.POST("/favorites", favouriteController.AddFavorite)
	router.GET("/favorites/:user_id", favouriteController.GetFavorites)
	router.DELETE("/favorites/:id", favouriteController.RemoveFavorite)

	return router
}

// -------------------------------
// Test Add Favorite
// -------------------------------
func TestAddFavorite(t *testing.T) {
	cleanFavoritesCollection()
	router := setUpFavoriteRouter()

	body := []byte(`{
		"user_id": "user123",
		"product_id": "prod123"
	}`)

	req, _ := http.NewRequest("POST", "/favorites", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	assert.Equal(t, "added to favorites", resp["message"])
	assert.NotNil(t, resp["favorite"])
}
