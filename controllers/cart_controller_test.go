package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"adhomes-backend/config"
	"adhomes-backend/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func cleanCartCollection() {
	if config.DB == nil {
		panic("Database not initialized! Did you run TestMain?")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := config.DB.Collection("carts").Drop(ctx); err != nil {
		panic(err)
	}
}

func setUpCartRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Initialize cart service and controller
	cartService := services.NewCartService()
	InitCartController(cartService)

	router.POST("/carts", CreateCart)
	router.DELETE("/carts/:id", DeleteCart)
	router.PUT("/carts/:id", UpdateCart)

	return router
}

func TestToCreateCart(t *testing.T) {
	cleanCartCollection()
	router := setUpCartRouter()

	body := []byte(`{
		"user_id": "user123",
		"items": [
			{ "product_id": "60f6e8bb5f627b5a1c9e4f31", "quantity": 2 },
			{ "product_id": "60f6e8bb5f627b5a1c9e4f32", "quantity": 1 }
		]
	}`)

	req, _ := http.NewRequest("POST", "/carts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "cart created", response["message"])
	assert.NotNil(t, response["cart"])
}

func TestToDeleteCart(t *testing.T) {
	cleanCartCollection()
	router := setUpCartRouter()

	// Create a cart first
	createBody := []byte(`{
        "user_id": "user123",
        "items": [
			{ "product_id": "60f6e8bb5f627b5a1c9e4f31", "quantity": 2 }
		]
    }`)
	createReq, _ := http.NewRequest("POST", "/carts", bytes.NewBuffer(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, createReq)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	cartMap := resp["cart"].(map[string]interface{})
	cartID := cartMap["id"].(string)

	// Delete the cart
	deleteReq, _ := http.NewRequest("DELETE", "/carts/"+cartID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, deleteReq)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestToDeleteNonExistentCart(t *testing.T) {
	cleanCartCollection()
	router := setUpCartRouter()

	nonExistentID := "60b8d295f1d2c12a34567890"
	deleteReq, _ := http.NewRequest("DELETE", "/carts/"+nonExistentID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, deleteReq)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestToDeleteCartInvalidID(t *testing.T) {
	cleanCartCollection()
	router := setUpCartRouter()

	invalidID := "invalid-id"
	deleteReq, _ := http.NewRequest("DELETE", "/carts/"+invalidID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, deleteReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestToUpdateCart(t *testing.T) {
	cleanCartCollection()
	router := setUpCartRouter()

	// Create cart first
	createBody := []byte(`{
		"user_id": "user123",
		"items": [
			{ "product_id": "60f6e8bb5f627b5a1c9e4f31", "quantity": 2 }
		]
	}`)
	createReq, _ := http.NewRequest("POST", "/carts", bytes.NewBuffer(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, createReq)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	cartMap := resp["cart"].(map[string]interface{})
	cartID := cartMap["id"].(string)

	// Update cart
	updateBody := []byte(`{
		"user_id": "user123",
		"items": [
			{ "product_id": "60f6e8bb5f627b5a1c9e4f33", "quantity": 5 }
		]
	}`)
	updateReq, _ := http.NewRequest("PUT", "/carts/"+cartID, bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, updateReq)

	assert.Equal(t, http.StatusOK, w.Code)
}
