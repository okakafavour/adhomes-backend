package services_impl

import (
	"context"
	"errors"
	"time"

	"adhomes-backend/config"
	"adhomes-backend/models"
	"adhomes-backend/services"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// WalletServiceImpl implements services.WalletService
type WalletServiceImpl struct {
	Collection *mongo.Collection
}

// Constructor
func NewWalletService() services.WalletService {
	return &WalletServiceImpl{
		Collection: config.DB.Collection("wallets"),
	}
}

// -----------------------------
// Get wallet by user ID
// -----------------------------
func (w *WalletServiceImpl) GetWalletByUserID(ctx context.Context, userID string) (*models.Wallet, error) {
	var wallet models.Wallet
	err := w.Collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&wallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("wallet not found")
		}
		return nil, err
	}
	return &wallet, nil
}

// -----------------------------
// Increase wallet balance
// -----------------------------
func (w *WalletServiceImpl) IncreaseBalance(ctx context.Context, userID string, amount float64) (*models.Wallet, error) {
	update := bson.M{
		"$inc": bson.M{"balance": amount},
		"$set": bson.M{"updated_at": time.Now()},
	}

	res := w.Collection.FindOneAndUpdate(ctx, bson.M{"user_id": userID}, update)
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return nil, errors.New("wallet not found")
		}
		return nil, res.Err()
	}

	var updatedWallet models.Wallet
	res.Decode(&updatedWallet)
	return &updatedWallet, nil
}

// -----------------------------
// Decrease wallet balance
// -----------------------------
func (w *WalletServiceImpl) DecreaseBalance(ctx context.Context, userID string, amount float64) (*models.Wallet, error) {
	// Check if balance is sufficient
	wallet, err := w.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if wallet.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	update := bson.M{
		"$inc": bson.M{"balance": -amount},
		"$set": bson.M{"updated_at": time.Now()},
	}

	res := w.Collection.FindOneAndUpdate(ctx, bson.M{"user_id": userID}, update)
	if res.Err() != nil {
		return nil, res.Err()
	}

	var updatedWallet models.Wallet
	res.Decode(&updatedWallet)
	return &updatedWallet, nil
}
