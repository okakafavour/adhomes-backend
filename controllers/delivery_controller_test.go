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

func cleanDeliveryCollection() {
	if config.DB == nil {
		panic("Database not initialized! Did you run TestMain?")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := config.DB.Collection("deliveries").Drop(ctx)
	if err != nil {
		panic(err)
	}
}

func setUpDeliveryRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Initialize services + controller
	InitDeliveryController()

	router.POST("/deliveries/assign", deliveryController.AssignRider)
	router.PUT("/deliveries/:id/status", deliveryController.UpdateStatus)
	router.GET("/deliveries/order/:order_id", deliveryController.GetByOrder)

	return router
}

func TestAssignRider(t *testing.T) {
	cleanDeliveryCollection()
	router := setUpDeliveryRouter()

	body := []byte(`{"order_id":"order123","rider_id":"rider123"}`)

	req, _ := http.NewRequest("POST", "/deliveries/assign", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "delivery assigned", response["message"])
	assert.NotNil(t, response["delivery"])
}

func TestUpdateDeliveryStatus(t *testing.T) {
	cleanDeliveryCollection()
	router := setUpDeliveryRouter()

	// First assign a delivery
	assignBody := []byte(`{"order_id":"order123","rider_id":"rider123"}`)
	assignReq, _ := http.NewRequest("POST", "/deliveries/assign", bytes.NewBuffer(assignBody))
	assignReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, assignReq)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	delMap := resp["delivery"].(map[string]interface{})
	id := delMap["id"].(string)

	// Update status
	updateBody := []byte(`{"status":"PickedUp"}`)
	updateReq, _ := http.NewRequest("PUT", "/deliveries/"+id+"/status", bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, updateReq)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetDeliveryByOrder(t *testing.T) {
	cleanDeliveryCollection()
	router := setUpDeliveryRouter()

	assignBody := []byte(`{"order_id":"order999","rider_id":"riderABC"}`)
	assignReq, _ := http.NewRequest("POST", "/deliveries/assign", bytes.NewBuffer(assignBody))
	assignReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, assignReq)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Fetch delivery by order
	req, _ := http.NewRequest("GET", "/deliveries/order/order999", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "order999", response["delivery"].(map[string]interface{})["order_id"])
}
