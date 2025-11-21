package services_impl

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"adhomes-backend/config"
	"adhomes-backend/models"
	"adhomes-backend/services"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type PaymentServiceImpl struct {
	PaymentCollection *mongo.Collection
	OrderCollection   *mongo.Collection
	HttpClient        HTTPClient
}

// ------------------------
// Constructor
// ------------------------
func NewPaymentService() services.PaymentService {
	return &PaymentServiceImpl{
		PaymentCollection: config.DB.Collection("payments"),
		OrderCollection:   config.DB.Collection("orders"),
		HttpClient:        &http.Client{Timeout: 10 * time.Second},
	}
}

// ------------------------
// Initialize Payment (Paystack)
// ------------------------
func (s *PaymentServiceImpl) InitializePayment(orderID, userID string, amount float64, email string) (*models.Payment, string, error) {
	reference := uuid.New().String()

	payment := models.Payment{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		OrderID:   orderID,
		Amount:    amount,
		Reference: reference,
		Gateway:   "paystack",
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save payment in DB
	_, err := s.PaymentCollection.InsertOne(context.Background(), payment)
	if err != nil {
		return nil, "", err
	}

	// Paystack request
	reqBody := map[string]interface{}{
		"email":     email,
		"amount":    int(amount * 100),
		"reference": reference,
	}

	bodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "https://api.paystack.co/transaction/initialize", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("PAYSTACK_SECRET_KEY"))

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	var paystackRes map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&paystackRes)

	if resp.StatusCode != 200 {
		return nil, "", errors.New("failed to initialize payment with Paystack")
	}

	data := paystackRes["data"].(map[string]interface{})
	paymentURL := data["authorization_url"].(string)

	return &payment, paymentURL, nil
}

// ------------------------
// Update Payment Status (callback from Paystack)
// ------------------------
func (s *PaymentServiceImpl) UpdatePaymentStatus(reference, status string) (*models.Payment, error) {
	var payment models.Payment

	// Find payment by reference
	err := s.PaymentCollection.FindOne(context.Background(), bson.M{"reference": reference}).Decode(&payment)
	if err != nil {
		return nil, errors.New("payment not found")
	}

	// Update payment status
	_, err = s.PaymentCollection.UpdateOne(
		context.Background(),
		bson.M{"reference": reference},
		bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}},
	)
	if err != nil {
		return nil, err
	}

	if status == "success" {
		_, _ = s.OrderCollection.UpdateOne(
			context.Background(),
			bson.M{"_id": payment.OrderID},
			bson.M{"$set": bson.M{"status": "Processing"}},
		)
	}

	payment.Status = status
	return &payment, nil
}

// ------------------------
// Get Payment by Reference
// ------------------------
func (s *PaymentServiceImpl) GetPaymentByReference(reference string) (*models.Payment, error) {
	var payment models.Payment
	err := s.PaymentCollection.FindOne(context.Background(), bson.M{"reference": reference}).Decode(&payment)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}
