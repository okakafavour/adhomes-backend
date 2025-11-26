package services

import "adhomes-backend/models"

type ProductService interface {
	AddProduct(name, description string, price float64, category, imageURL string) (models.Product, error)
	UpdateProduct(id string, product models.Product) (models.Product, error)
	DeleteProduct(id string) error
	GetAllProducts() ([]models.Product, error)
}
