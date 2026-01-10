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

func (s *ProductServiceImpl) AddProduct(
	product *models.Product,
) (models.Product, error) {

	return s.productRepo.CreateProduct(*product)
}
func (s *ProductServiceImpl) UpdateProduct(
	id string,
	product *models.Product,
) (models.Product, error) {

	return s.productRepo.UpdateProduct(id, *product)
}

func (s *ProductServiceImpl) DeleteProduct(id string) error {
	return s.productRepo.DeleteProduct(id)
}

func (s *ProductServiceImpl) GetAllProducts() ([]models.Product, error) {
	return s.productRepo.FindAll()
}

func (s *ProductServiceImpl) GetProductByID(id string) (models.Product, error) {
	return s.productRepo.FindByID(id)
}
