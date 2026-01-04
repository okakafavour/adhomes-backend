package services_impl

import (
	"errors"
	"time"

	"adhomes-backend/models"
	"adhomes-backend/repositories"
	"adhomes-backend/services"
	"adhomes-backend/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type userServiceImpl struct {
	userRepo *repositories.UserRepository
}

func NewUserService() services.UserService {
	return &userServiceImpl{
		userRepo: repositories.NewUserRepository(),
	}
}

// ------------------
// Register
// ------------------
func (s *userServiceImpl) Register(user models.User) error {
	exists, err := s.userRepo.EmailExists(user.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already registered")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashed)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return s.userRepo.CreateUser(user)
}

// ------------------
// Login
// ------------------
func (s *userServiceImpl) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindUserByEmail(email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", errors.New("invalid email or password")
		}
		return "", err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}
