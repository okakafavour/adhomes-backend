package services

import "adhomes-backend/models"

type ProductService interface {
	AddProduct(product *models.Product) (models.Product, error)
	UpdateProduct(id string, product *models.Product) (models.Product, error)
	DeleteProduct(id string) error
	GetAllProducts() ([]models.Product, error)
	GetProductByID(id string) (models.Product, error)
}
