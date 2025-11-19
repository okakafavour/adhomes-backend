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
	"adhomes-backend/models"
	"adhomes-backend/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// -------------------------------
// Helpers
// -------------------------------

func cleanOrdersCollection() {
	if config.DB == nil {
		panic("Database not initialized! Did you run TestMain?")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := config.DB.Collection("orders").Drop(ctx); err != nil {
		panic(err)
	}
}

func setUpOrderRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	InitOrderController()

	router.POST("/orders", CreateOrder)
	router.GET("/orders/:id", GetOrderByID)
	router.GET("/orders/user/:user_id", GetOrdersByUserID)
	router.DELETE("/orders/:id", DeleteOrder)
	router.PUT("/orders/:id", UpdateOrder)

	return router
}

// -------------------------------
// TESTS
// -------------------------------

func TestToCreateOrder(t *testing.T) {
	cleanOrdersCollection()
	router := setUpOrderRouter()

	body := []byte(`{
		"user_id": "user123",
		"items": [
			{ "product_id": "prod1", "quantity": 2 },
			{ "product_id": "prod2", "quantity": 1 }
		],
		"total": 150.0,
		"delivery_address": "123 Main St",
		"payment_status": "Pending",
		"status": "Pending"
	}`)

	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "order created", resp["message"])
	assert.NotNil(t, resp["order"])
}

func TestToGetOrderByID(t *testing.T) {
	cleanOrdersCollection()
	router := setUpOrderRouter()

	order := models.Order{
		UserID:          "user123",
		Items:           []models.OrderItem{{ProductID: "prod1", Quantity: 2}},
		Total:           100.0,
		DeliveryAddress: "123 Main St",
		PaymentStatus:   "Pending",
		OrderStatus:     "Pending",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	created, _ := services.NewOrderService().CreateOrder(order)

	req, _ := http.NewRequest("GET", "/orders/"+created.ID.Hex(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.Order
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, created.ID.Hex(), resp.ID.Hex())
}

func TestToGetOrdersByUserID(t *testing.T) {
	cleanOrdersCollection()
	router := setUpOrderRouter()
	service := services.NewOrderService()

	service.CreateOrder(models.Order{
		UserID:          "user123",
		Items:           []models.OrderItem{{ProductID: "prod1", Quantity: 1}},
		Total:           50.0,
		DeliveryAddress: "123 Main St",
		PaymentStatus:   "Pending",
		OrderStatus:     "Pending",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	})
	service.CreateOrder(models.Order{
		UserID:          "user123",
		Items:           []models.OrderItem{{ProductID: "prod2", Quantity: 2}},
		Total:           100.0,
		DeliveryAddress: "123 Main St",
		PaymentStatus:   "Pending",
		OrderStatus:     "Pending",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	})

	req, _ := http.NewRequest("GET", "/orders/user/user123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []models.Order
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp, 2)
}

func TestToDeleteOrder(t *testing.T) {
	cleanOrdersCollection()
	router := setUpOrderRouter()

	order := models.Order{
		UserID:          "user123",
		Items:           []models.OrderItem{{ProductID: "prod1", Quantity: 2}},
		Total:           100.0,
		DeliveryAddress: "123 Main St",
		PaymentStatus:   "Pending",
		OrderStatus:     "Pending",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	created, _ := services.NewOrderService().CreateOrder(order)

	req, _ := http.NewRequest("DELETE", "/orders/"+created.ID.Hex(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestToGetNonExistentOrder(t *testing.T) {
	cleanOrdersCollection()
	router := setUpOrderRouter()

	nonExistentID := "64b8d295f1d2c12a34567890"
	req, _ := http.NewRequest("GET", "/orders/"+nonExistentID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "order not found", resp["error"])
}

func TestToDeleteNonExistentOrder(t *testing.T) {
	cleanOrdersCollection()
	router := setUpOrderRouter()

	nonExistentID := "64b8d295f1d2c12a34567890"
	req, _ := http.NewRequest("DELETE", "/orders/"+nonExistentID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "order not found", resp["error"])
}

func TestToUpdateOrderStatusAndPayment(t *testing.T) {
	cleanOrdersCollection()
	router := setUpOrderRouter()

	order := models.Order{
		UserID:          "user123",
		Items:           []models.OrderItem{{ProductID: "prod1", Quantity: 2}},
		Total:           150.0,
		DeliveryAddress: "123 Main St",
		PaymentStatus:   "Pending",
		OrderStatus:     "Pending",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	created, _ := services.NewOrderService().CreateOrder(order)

	updateBody := []byte(`{
		"payment_status": "Paid",
		"status": "Processing"
	}`)
	req, _ := http.NewRequest("PUT", "/orders/"+created.ID.Hex(), bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	orderResp := resp["order"].(map[string]interface{})
	assert.Equal(t, "Paid", orderResp["payment_status"])
	assert.Equal(t, "Processing", orderResp["status"])
}
