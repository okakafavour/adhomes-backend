package services

import "adhomes-backend/models"

type PaymentService interface {
	InitializePayment(orderID string, userID string, amount float64, email string) (*models.Payment, string, error)
	UpdatePaymentStatus(refrenece string, status string) (*models.Payment, error)
	GetPaymentByReference(refrence string) (*models.Payment, error)
}
