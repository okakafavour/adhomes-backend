package services_impl

import (
	"context"
	"errors"

	"adhomes-backend/models"
	"adhomes-backend/repositories"
	"adhomes-backend/services"
)

type WalletServiceImpl struct {
	walletRepo *repositories.WalletRepository
}

// Constructor
func NewWalletService() services.WalletService {
	return &WalletServiceImpl{
		walletRepo: repositories.NewWalletRepository(),
	}
}

// -----------------------------
// Get wallet by user ID
// -----------------------------
func (w *WalletServiceImpl) GetWalletByUserID(
	ctx context.Context,
	userID string,
) (*models.Wallet, error) {
	return w.walletRepo.FindByUserID(ctx, userID)
}

// -----------------------------
// Increase wallet balance
// -----------------------------
func (w *WalletServiceImpl) IncreaseBalance(
	ctx context.Context,
	userID string,
	amount float64,
) (*models.Wallet, error) {

	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	return w.walletRepo.IncreaseBalance(ctx, userID, amount)
}

// -----------------------------
// Decrease wallet balance
// -----------------------------
func (w *WalletServiceImpl) DecreaseBalance(
	ctx context.Context,
	userID string,
	amount float64,
) (*models.Wallet, error) {

	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	wallet, err := w.walletRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if wallet.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	return w.walletRepo.DecreaseBalance(ctx, userID, amount)
}
