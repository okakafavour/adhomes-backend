package services_impl

import (
	"context"
	"errors"
	"time"

	"adhomes-backend/config"
	"adhomes-backend/models"
	"adhomes-backend/services"
	"adhomes-backend/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type userServiceImpl struct {
	collection *mongo.Collection
}

func NewUserService() services.UserService {
	return &userServiceImpl{
		collection: config.GetCollection("users"),
	}
}

// ------------------
// Register
// ------------------
func (s *userServiceImpl) Register(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if email exists
	count, err := s.collection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("email already registered")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err = s.collection.InsertOne(ctx, user)
	return err
}

// ------------------
// Login
// ------------------
func (s *userServiceImpl) Login(email, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var found models.User
	err := s.collection.FindOne(ctx, bson.M{"email": email}).Decode(&found)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", errors.New("invalid email or password")
		}
		return "", err
	}

	if bcrypt.CompareHashAndPassword([]byte(found.Password), []byte(password)) != nil {
		return "", errors.New("invalid email or password")
	}

	// Generate token
	token, err := utils.GenerateToken(found.Email)
	if err != nil {
		return "", err
	}
	return token, nil
}
