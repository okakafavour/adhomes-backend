package services_impl

import (
	"adhomes-backend/models"
	"adhomes-backend/repositories"
	"adhomes-backend/services"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type paymentServiceImpl struct {
	paymentRepo *repositories.PaymentRepository
	orderRepo   *repositories.OrderRepository
	walletRepo  *repositories.WalletRepository
}

// NewPaymentService creates a new PaymentService
func NewPaymentService() services.PaymentService {
	return &paymentServiceImpl{
		paymentRepo: repositories.NewPaymentRepository(),
		orderRepo:   repositories.NewOrderRepository(),
		walletRepo:  repositories.NewWalletRepository(),
	}
}

// MakePayment handles wallet or paystack payment
func (s *paymentServiceImpl) MakePayment(req models.PaymentRequest) (models.Payment, string, error) {
	ctx := context.Background()

	// 1️⃣ Validate order
	order, err := s.orderRepo.FindOrderByID(req.OrderID)
	if err != nil {
		return models.Payment{}, "", err
	}

	// Ensure the order belongs to the user
	if order.CustomerEmail != req.UserID {
		return models.Payment{}, "", errors.New("order does not belong to user")
	}

	// Ensure the amount matches
	if order.TotalAmount != req.Amount {
		return models.Payment{}, "", errors.New("amount does not match order total")
	}

	// Create a new payment record
	payment := models.Payment{
		ID:        primitive.NewObjectID(),
		UserID:    req.UserID,
		OrderID:   req.OrderID,
		Amount:    req.Amount,
		Email:     req.Email,
		Method:    req.PaymentMethod,
		Status:    "pending",
		Reference: primitive.NewObjectID().Hex(),
	}

	// 2️⃣ Wallet payment
	if req.PaymentMethod == "wallet" {
		wallet, err := s.walletRepo.FindByUserID(ctx, req.UserID)
		if err != nil {
			return models.Payment{}, "", err
		}

		if wallet.Balance < req.Amount {
			return models.Payment{}, "", errors.New("insufficient wallet balance")
		}

		_, err = s.walletRepo.DecreaseBalance(ctx, req.UserID, req.Amount)
		if err != nil {
			return models.Payment{}, "", err
		}

		payment.Status = "success"
		if err := s.orderRepo.UpdateOrderStatus(req.OrderID, "paid"); err != nil {
			return models.Payment{}, "", err
		}

		// Save payment
		p, err := s.paymentRepo.Create(payment)
		return p, "", err
	}

	// 3️⃣ Paystack payment (mocked)
	if req.PaymentMethod == "paystack" {
		paymentURL := "https://paystack.com/pay/" + payment.Reference

		p, err := s.paymentRepo.Create(payment)
		return p, paymentURL, err
	}

	return models.Payment{}, "", errors.New("invalid payment method")
}
