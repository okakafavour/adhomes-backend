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
	r := gin.New()
	r.POST("/carts", CreateCart)
	r.DELETE("/carts/:id", DeleteCart)
	r.PUT("/carts/:id", UpdateCart)
	return r
}

func TestToCreateCart(t *testing.T) {
	cleanCartCollection()
	router := setUpCartRouter()

	body := []byte(`{
		"user_id": "user123",
		"product_ids": ["prod1", "prod2"]
	}`)

	req, _ := http.NewRequest("POST", "/carts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "Cart created successfully", response["message"])
	assert.NotNil(t, response["cart"])
}

func TestToDeleteCart(t *testing.T) {
	if config.DB == nil {
		t.Fatal("Database not initialized! Did you call config.ConnectDB()?")
	}
	InitCartController()

	cleanCartCollection()

	router := setUpCartRouter()

	createBody := []byte(`{
        "user_id": "user123",
        "product_ids": ["prod1", "prod2"]
    }`)
	createReq, _ := http.NewRequest("POST", "/carts", bytes.NewBuffer(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, createReq)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201 when creating cart, got %d", w.Code)
	}

	// Parse the response to get cart ID
	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("failed to parse create cart response: %v", err)
	}

	cartMap, ok := resp["cart"].(map[string]interface{})
	if !ok {
		t.Fatal("response does not contain cart object")
	}

	cartID, ok := cartMap["id"].(string)
	if !ok {
		t.Fatal("cart object does not contain id")
	}

	// Step 4: Delete the cart
	deleteReq, _ := http.NewRequest("DELETE", "/carts/"+cartID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, deleteReq)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 when deleting cart, got %d", w.Code)
	}
}

func TestToDeleteNonExistentCart(t *testing.T) {
	cleanCartCollection()

	router := setUpCartRouter()

	nonExistentID := "60b8d295f1d2c12a34567890"
	deleteReq, _ := http.NewRequest("DELETE", "/carts/"+nonExistentID, nil)
	deleteW := httptest.NewRecorder()
	router.ServeHTTP(deleteW, deleteReq)

	assert.Equal(t, http.StatusNotFound, deleteW.Code)

	var deleteResp map[string]interface{}
	json.Unmarshal(deleteW.Body.Bytes(), &deleteResp)

	assert.Equal(t, "Cart not found", deleteResp["error"])
}

func TestToDeleteCartInvalidID(t *testing.T) {
	cleanCartCollection()

	router := setUpCartRouter()

	invalidID := "invalid-id"
	deleteReq, _ := http.NewRequest("DELETE", "/carts/"+invalidID, nil)
	deleteW := httptest.NewRecorder()
	router.ServeHTTP(deleteW, deleteReq)

	assert.Equal(t, http.StatusBadRequest, deleteW.Code)

	var deleteResp map[string]interface{}
	json.Unmarshal(deleteW.Body.Bytes(), &deleteResp)

	assert.Equal(t, "Invalid cart ID", deleteResp["error"])
}

func TestToUpdateCart(t *testing.T) {
	if config.DB == nil {
		t.Fatal("Database not initialized! Did you call config.ConnectDB()?")
	}
	InitCartController()

	cleanCartCollection()

	router := setUpCartRouter()

	createBody := []byte(`{
		"user_id": "user123",
		"product_ids": ["prod1", "prod2"]
	}`)
	createReq, _ := http.NewRequest("POST", "/carts", bytes.NewBuffer(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, createReq)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201 when creating cart, got %d", w.Code)
	}

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("failed to parse create cart response: %v", err)
	}

	cartMap, ok := resp["cart"].(map[string]interface{})
	if !ok {
		t.Fatal("response does not contain cart object")
	}

	cartID, ok := cartMap["id"].(string)
	if !ok {
		t.Fatal("cart object does not contain id")
	}

	updateBody := []byte(`{
		"user_id": "user123",
		"product_ids": ["prod3", "prod4"]
	}`)
	updateReq, _ := http.NewRequest("PUT", "/carts/"+cartID, bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateW := httptest.NewRecorder()
	router.ServeHTTP(updateW, updateReq)

	if updateW.Code != http.StatusOK {
		t.Fatalf("expected status 200 when updating cart, got %d", updateW.Code)
	}
}
