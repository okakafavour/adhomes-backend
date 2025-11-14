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

// -------------------------------
// Helpers
// -------------------------------

// Clean products collection safely
func cleanProductsCollection() {
	if config.DB == nil {
		panic("Database not initialized! Did you run TestMain?")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := config.DB.Collection("products").Drop(ctx); err != nil {
		panic(err)
	}
}

// Setup router
func setUpProductRouter() *gin.Engine {
	r := gin.New()
	r.POST("/products", CreateProduct)
	r.PUT("/products/:id", UpdateProduct)
	r.DELETE("/products/:id", DeleteProduct)
	return r
}

// -------------------------------
// Tests
// -------------------------------

func TestToCreateProduct(t *testing.T) {
	cleanProductsCollection()
	router := setUpProductRouter()

	body := []byte(`{
		"name": "Test Product",
		"description": "This is a test product",
		"price": 19.99,
		"category": "Test Category",
		"image_url": "http://example.com/image.jpg"
	}`)

	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Equal(t, "Product created successfully", resp["message"])
	product := resp["product"].(map[string]interface{})
	assert.Equal(t, "Test Product", product["name"])
}

func TestToDeleteProduct(t *testing.T) {
	cleanProductsCollection()

	router := setUpProductRouter()

	body := []byte(`{
		"name": "Test Product",
		"description": "This is a test product",
		"price": 19.99,
		"category": "Test Category",
		"image_url": "http://example.com/image.jpg"
	}`)

	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	product := resp["product"].(map[string]interface{})
	productID := product["id"].(string)
	assert.NotEmpty(t, productID)

	deleteReq, _ := http.NewRequest("DELETE", "/products/"+productID, nil)
	deleteW := httptest.NewRecorder()
	router.ServeHTTP(deleteW, deleteReq)

	assert.Equal(t, http.StatusOK, deleteW.Code)

	var deleteResp map[string]interface{}
	json.Unmarshal(deleteW.Body.Bytes(), &deleteResp)

	assert.Equal(t, "Product deleted successfully", deleteResp["message"])
}

func TestToDeleteNonExistentProduct(t *testing.T) {
	cleanProductsCollection()

	router := setUpProductRouter()

	nonExistentID := "60b8d295f1d2c12a34567890"
	deleteReq, _ := http.NewRequest("DELETE", "/products/"+nonExistentID, nil)
	deleteW := httptest.NewRecorder()
	router.ServeHTTP(deleteW, deleteReq)

	assert.Equal(t, http.StatusNotFound, deleteW.Code)

	var deleteResp map[string]interface{}
	json.Unmarshal(deleteW.Body.Bytes(), &deleteResp)

	assert.Equal(t, "Product not found", deleteResp["error"])
}

func TestToDeleteProductInvalidID(t *testing.T) {
	cleanProductsCollection()

	router := setUpProductRouter()

	invalidID := "invalid-id"
	deleteReq, _ := http.NewRequest("DELETE", "/products/"+invalidID, nil)
	deleteW := httptest.NewRecorder()
	router.ServeHTTP(deleteW, deleteReq)

	assert.Equal(t, http.StatusBadRequest, deleteW.Code)

	var deleteResp map[string]interface{}
	json.Unmarshal(deleteW.Body.Bytes(), &deleteResp)

	assert.Equal(t, "Invalid product ID", deleteResp["error"])
}

func TestToUpdateProduct(t *testing.T) {
	cleanProductsCollection()

	router := setUpProductRouter()

	body := []byte(`{
		"name": "Test Product",
		"description": "This is a test product",
		"price": 19.99,
		"category": "Test Category",
		"image_url": "http://example.com/image.jpg"
	}`)

	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	product := resp["product"].(map[string]interface{})
	productID := product["id"].(string)
	assert.NotEmpty(t, productID)

	updateBody := []byte(`{
		"name": "Updated Product",
		"description": "This is an updated test product",
		"price": 29.99,
		"category": "Updated Category",
		"image_url": "http://example.com/updated_image.jpg"
	}`)

	updateReq, _ := http.NewRequest("PUT", "/products/"+productID, bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateW := httptest.NewRecorder()
	router.ServeHTTP(updateW, updateReq)

	assert.Equal(t, http.StatusOK, updateW.Code)

	var updateResp map[string]interface{}
	json.Unmarshal(updateW.Body.Bytes(), &updateResp)

	assert.Equal(t, "Product updated successfully", updateResp["message"])
	updatedProduct := updateResp["product"].(map[string]interface{})
	assert.Equal(t, "Updated Product", updatedProduct["name"])
	assert.Equal(t, "This is an updated test product", updatedProduct["description"])
	assert.Equal(t, 29.99, updatedProduct["price"])
}

func TestToUpdateProductInvalidID(t *testing.T) {
	cleanProductsCollection()

	router := setUpProductRouter()

	invalidID := "invalid-id"
	updateBody := []byte(`{
		"name": "Updated Product",
		"description": "This is an updated test product",
		"price": 29.99,
		"category": "Updated Category",
		"image_url": "http://example.com/updated_image.jpg"
	}`)

	updateReq, _ := http.NewRequest("PUT", "/products/"+invalidID, bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateW := httptest.NewRecorder()
	router.ServeHTTP(updateW, updateReq)

	assert.Equal(t, http.StatusBadRequest, updateW.Code)

	var updateResp map[string]interface{}
	json.Unmarshal(updateW.Body.Bytes(), &updateResp)

	assert.Equal(t, "Invalid product ID", updateResp["error"])
}
