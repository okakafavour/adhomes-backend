package services_impl

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"adhomes-backend/models"
	"adhomes-backend/repositories"
	"adhomes-backend/services"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type PaymentServiceImpl struct {
	paymentRepo *repositories.PaymentRepository
	httpClient  HTTPClient
}

// ------------------------
// Constructor
// ------------------------
func NewPaymentService() services.PaymentService {
	return &PaymentServiceImpl{
		paymentRepo: repositories.NewPaymentRepository(),
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}
}

// ------------------------
// Initialize Payment (Paystack)
// ------------------------
func (s *PaymentServiceImpl) InitializePayment(
	orderID,
	userID string,
	amount float64,
	email string,
) (*models.Payment, string, error) {

	if amount <= 0 {
		return nil, "", errors.New("invalid amount")
	}

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

	ctx := context.Background()

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, "", err
	}

	reqBody := map[string]interface{}{
		"email":     email,
		"amount":    int(amount * 100),
		"reference": reference,
	}

	bodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(
		"POST",
		"https://api.paystack.co/transaction/initialize",
		bytes.NewBuffer(bodyBytes),
	)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("PAYSTACK_SECRET_KEY"))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", errors.New("failed to initialize payment")
	}

	var paystackRes struct {
		Data struct {
			AuthorizationURL string `json:"authorization_url"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&paystackRes); err != nil {
		return nil, "", err
	}

	return &payment, paystackRes.Data.AuthorizationURL, nil
}

// ------------------------
// Update Payment Status (Webhook)
// ------------------------
func (s *PaymentServiceImpl) UpdatePaymentStatus(
	reference string,
	status string,
) (*models.Payment, error) {

	ctx := context.Background()

	payment, err := s.paymentRepo.FindByReference(ctx, reference)
	if err != nil {
		return nil, err
	}

	if err := s.paymentRepo.UpdateStatus(ctx, reference, status); err != nil {
		return nil, err
	}

	if status == "success" {
		_ = s.paymentRepo.UpdateOrderStatus(
			ctx,
			payment.OrderID,
			"Processing",
		)
	}

	payment.Status = status
	return payment, nil
}

// ------------------------
// Get Payment by Reference
// ------------------------
func (s *PaymentServiceImpl) GetPaymentByReference(
	reference string,
) (*models.Payment, error) {

	return s.paymentRepo.FindByReference(
		context.Background(),
		reference,
	)
}
