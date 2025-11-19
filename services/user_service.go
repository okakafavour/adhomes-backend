package services

import "adhomes-backend/models"

type UserService interface {
	Register(user models.User) error
	Login(email, password string) (string, error)
}
