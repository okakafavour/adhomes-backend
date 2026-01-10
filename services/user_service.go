package services

import "adhomes-backend/models"

type UserService interface {
	Register(user models.User) error
	Login(email, password string) (string, error)

	// Admin actions
	GetAllUsers() ([]models.User, error)
	DeactivateUser(userID string) error
	ActivateUser(userID string) error
	DeleteUser(userID string) error
}
