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
	"adhomes-backend/services_impl"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// HTTPClient interface to allow fake client for testing
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Fake HTTP client for Paystack
type fakeHTTPClient struct{}

func (f *fakeHTTPClient) Do(req *http.Request) (*http.Response, error) {
	dummy := `{
		"status": true,
		"message": "Authorization URL created",
		"data": {"authorization_url": "https://paystack.com/fake-url"}
	}`
	resp := &http.Response{
		StatusCode: 200,
		Body:       &bodyCloser{data: []byte(dummy)},
		Header:     make(http.Header),
	}
	return resp, nil
}

// Implement io.ReadCloser
type bodyCloser struct{ data []byte }

func (b *bodyCloser) Read(p []byte) (int, error) {
	n := copy(p, b.data)
	b.data = b.data[n:]
	return n, nil
}
func (b *bodyCloser) Close() error { return nil }

func TestMakePayment(t *testing.T) {
	config.ConnectDB()
	db := config.DB
	db.Collection("payments").Drop(context.Background())
	db.Collection("wallets").Drop(context.Background())

	// Create test wallet
	wallet := models.Wallet{
		UserID:  "user123",
		Balance: 5000,
	}
	db.Collection("wallets").InsertOne(context.Background(), wallet)

	// Setup services
	paymentService := &services_impl.PaymentServiceImpl{
		PaymentCollection: db.Collection("payments"),
		OrderCollection:   db.Collection("orders"),
		HttpClient:        &fakeHTTPClient{}, // inject fake client
	}
	walletService := services_impl.NewWalletService()

	// Setup Gin router
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	controller := NewPaymentController(paymentService, walletService)
	router.POST("/payments", controller.MakePayment)

	// Prepare request
	reqBody := models.PaymentRequest{
		UserID:  "user123",
		OrderID: "order123",
		Amount:  1500,
		Email:   "user@test.com",
	}
	jsonReq, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/payments", bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Decode response
	var resp struct {
		Payment    models.Payment `json:"payment"`
		PaymentURL string         `json:"payment_url"`
	}
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, "user123", resp.Payment.UserID)
	assert.Equal(t, 1500.0, resp.Payment.Amount)
	assert.Equal(t, "order123", resp.Payment.OrderID)
	assert.NotZero(t, resp.Payment.ID)
	assert.WithinDuration(t, time.Now(), resp.Payment.CreatedAt, 5*time.Second)
	assert.Equal(t, "https://paystack.com/fake-url", resp.PaymentURL)
}
