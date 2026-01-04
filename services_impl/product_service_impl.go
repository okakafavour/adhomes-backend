package services_impl

import (
	"adhomes-backend/models"
	"adhomes-backend/repositories"
	"adhomes-backend/services"
)

type ProductServiceImpl struct {
	productRepo *repositories.ProductRepository
}

func NewProductService() services.ProductService {
	return &ProductServiceImpl{
		productRepo: repositories.NewProductRepository(),
	}
}

// OPTION 1 — Add product by fields
func (s *ProductServiceImpl) AddProduct(
	name, description string,
	price float64,
	category, imageURL string,
) (models.Product, error) {

	product := models.Product{
		Name:        name,
		Description: description,
		Price:       price,
		Category:    category,
		ImageURL:    imageURL,
	}

	return s.productRepo.CreateProduct(product)
}

func (s *ProductServiceImpl) UpdateProduct(id string, product models.Product) (models.Product, error) {
	return s.productRepo.UpdateProduct(id, product)
}

func (s *ProductServiceImpl) DeleteProduct(id string) error {
	return s.productRepo.DeleteProduct(id)
}

func (s *ProductServiceImpl) GetAllProducts() ([]models.Product, error) {
	return s.productRepo.FindAll()
}
