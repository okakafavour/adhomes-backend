package services

import "adhomes-backend/models"

type PaymentService interface {
	MakePayment(req models.PaymentRequest) (models.Payment, string, error)

	// VerifyPayment(reference string) error
}
