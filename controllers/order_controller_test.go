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

// Clean Orders
func cleanOrderCollection() {
	if config.DB == nil {
		panic("Database not initialized! Did you run TestMain?")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := config.DB.Collection("orders").Drop(ctx); err != nil {
		panic(err)
	}
}

// Router for Order tests
func setUpOrderRouter() *gin.Engine {
	r := gin.New()

	// IMPORTANT: You must include BOTH carts & orders routes
	r.POST("/carts", CreateCart)
	r.POST("/orders", CreateOrder)

	return r
}

func TestToCreateOrder(t *testing.T) {
	cleanOrderCollection()
	cleanCartCollection()

	router := setUpOrderRouter()

	// STEP 1 ➜ Create Cart
	cartBody := []byte(`{
		"user_id": "user_123",
		"product_ids": ["p1", "p2"]
	}`)

	cartReq, _ := http.NewRequest("POST", "/carts", bytes.NewBuffer(cartBody))
	cartReq.Header.Set("Content-Type", "application/json")

	cartW := httptest.NewRecorder()
	router.ServeHTTP(cartW, cartReq)

	assert.Equal(t, http.StatusCreated, cartW.Code)

	var cartResp map[string]interface{}
	json.Unmarshal(cartW.Body.Bytes(), &cartResp)

	cart := cartResp["cart"].(map[string]interface{})
	cartID := cart["id"].(string)

	// STEP 2 ➜ Create Order
	orderBody := []byte(`{
		"cart_id": "` + cartID + `",
        "delivery_address": "123 Food Street"
	}`)

	orderReq, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(orderBody))
	orderReq.Header.Set("Content-Type", "application/json")

	orderW := httptest.NewRecorder()
	router.ServeHTTP(orderW, orderReq)

	assert.Equal(t, http.StatusCreated, orderW.Code)

	var orderResp map[string]interface{}
	json.Unmarshal(orderW.Body.Bytes(), &orderResp)

	assert.Equal(t, "Order created successfully", orderResp["message"])
	assert.NotNil(t, orderResp["order"])

	order := orderResp["order"].(map[string]interface{})
	assert.Equal(t, "user_123", order["user_id"])
	assert.Equal(t, "Processing", order["status"])
}
