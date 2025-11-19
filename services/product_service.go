package services

import "adhomes-backend/models"

type ProductService interface {
	CreateProduct(product models.Product) (models.Product, error)
	UpdateProduct(id string, product models.Product) (models.Product, error)
	DeleteProduct(id string) error
}
