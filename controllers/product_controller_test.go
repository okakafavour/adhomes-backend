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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func cleanProductCollection() {
	if config.DB == nil {
		panic("Database not initialized! Did you run TestMain?")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := config.DB.Collection("products").Drop(ctx); err != nil {
		panic(err)
	}
}

func setUpProductRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	InitProductController()

	router.POST("/products", CreateProduct)
	router.PUT("/products/:id", UpdateProduct)
	router.DELETE("/products/:id", DeleteProduct)
	return router
}

func TestToCreateProduct(t *testing.T) {
	cleanProductCollection()
	router := setUpProductRouter()

	body := []byte(`{
		"name": "Laptop",
		"description": "Gaming Laptop",
		"price": 1500,
		"category": "Electronics",
		"image_url": "https://example.com/laptop.jpg"
	}`)

	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, "Product created successfully", resp["message"])
	assert.NotNil(t, resp["product"])
}

func TestToUpdateProduct(t *testing.T) {
	cleanProductCollection()
	router := setUpProductRouter()

	// First create a product
	product := map[string]interface{}{
		"name":        "Phone",
		"description": "Smartphone",
		"price":       500,
		"category":    "Electronics",
		"image_url":   "https://example.com/phone.jpg",
	}
	body, _ := json.Marshal(product)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	productID := resp["product"].(map[string]interface{})["id"].(string)

	// Now update
	update := map[string]interface{}{
		"name":        "Phone Pro",
		"description": "Smartphone Pro",
		"price":       700,
		"category":    "Electronics",
		"image_url":   "https://example.com/phonepro.jpg",
	}
	updateBody, _ := json.Marshal(update)
	updateReq, _ := http.NewRequest("PUT", "/products/"+productID, bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateW := httptest.NewRecorder()
	router.ServeHTTP(updateW, updateReq)

	assert.Equal(t, http.StatusOK, updateW.Code)
}

func TestToDeleteProduct(t *testing.T) {
	cleanProductCollection()
	router := setUpProductRouter()

	// First create a product
	product := map[string]interface{}{
		"name":        "Tablet",
		"description": "Android Tablet",
		"price":       300,
		"category":    "Electronics",
		"image_url":   "https://example.com/tablet.jpg",
	}
	body, _ := json.Marshal(product)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	productID := resp["product"].(map[string]interface{})["id"].(string)

	// Delete the product
	deleteReq, _ := http.NewRequest("DELETE", "/products/"+productID, nil)
	deleteW := httptest.NewRecorder()
	router.ServeHTTP(deleteW, deleteReq)

	assert.Equal(t, http.StatusOK, deleteW.Code)
}
