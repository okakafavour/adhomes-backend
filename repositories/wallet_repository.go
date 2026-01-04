package repositories

import (
	"context"
	"errors"
	"time"

	"adhomes-backend/config"
	"adhomes-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type WalletRepository struct {
	collection *mongo.Collection
}

func NewWalletRepository() *WalletRepository {
	return &WalletRepository{
		collection: config.DB.Collection("wallets"),
	}
}

func (r *WalletRepository) FindByUserID(ctx context.Context, userID string) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&wallet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("wallet not found")
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *WalletRepository) IncreaseBalance(
	ctx context.Context,
	userID string,
	amount float64,
) (*models.Wallet, error) {

	update := bson.M{
		"$inc": bson.M{"balance": amount},
		"$set": bson.M{"updated_at": time.Now()},
	}

	res := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"user_id": userID},
		update,
	)

	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return nil, errors.New("wallet not found")
		}
		return nil, res.Err()
	}

	var wallet models.Wallet
	if err := res.Decode(&wallet); err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *WalletRepository) DecreaseBalance(
	ctx context.Context,
	userID string,
	amount float64,
) (*models.Wallet, error) {

	update := bson.M{
		"$inc": bson.M{"balance": -amount},
		"$set": bson.M{"updated_at": time.Now()},
	}

	res := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"user_id": userID},
		update,
	)

	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return nil, errors.New("wallet not found")
		}
		return nil, res.Err()
	}

	var wallet models.Wallet
	if err := res.Decode(&wallet); err != nil {
		return nil, err
	}

	return &wallet, nil
}
