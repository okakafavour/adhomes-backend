package services

import (
	"adhomes-backend/models"
	"context"
)

type WalletService interface {
	GetWalletByUserID(ctx context.Context, userID string) (*models.Wallet, error)
	IncreaseBalance(ctx context.Context, userID string, amount float64) (*models.Wallet, error)
	DecreaseBalance(ctx context.Context, userID string, amount float64) (*models.Wallet, error)
}
