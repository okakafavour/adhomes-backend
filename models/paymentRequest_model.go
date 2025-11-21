package models

type PaymentRequest struct {
	UserID  string  `json:"user_id"`
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
	Email   string  `json:"email"`
}
