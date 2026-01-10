package services_impl

import (
	"errors"
	"time"

	"adhomes-backend/models"
	"adhomes-backend/repositories"
	"adhomes-backend/services"
	"adhomes-backend/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

//
// ==================== AUTH ====================
//

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

	user.ID = primitive.NewObjectID()
	user.Password = string(hashed)
	user.IsActive = true
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return s.userRepo.CreateUser(user)
}

func (s *userServiceImpl) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindUserByEmail(email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", errors.New("invalid email or password")
		}
		return "", err
	}

	if !user.IsActive {
		return "", errors.New("account is deactivated")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.Email, false)
	if err != nil {
		return "", err
	}

	return token, nil
}

//
// ==================== ADMIN ====================
//

func (s *userServiceImpl) GetAllUsers() ([]models.User, error) {
	return s.userRepo.FindAll()
}

func (s *userServiceImpl) DeactivateUser(userID string) error {
	return s.userRepo.UpdateStatus(userID, false)
}

func (s *userServiceImpl) ActivateUser(userID string) error {
	return s.userRepo.UpdateStatus(userID, true)
}

func (s *userServiceImpl) DeleteUser(userID string) error {
	return s.userRepo.DeleteUser(userID)
}
