package services_impl

import (
	"errors"
	"time"

	"adhomes-backend/models"
	"adhomes-backend/repositories"
	"adhomes-backend/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductServiceImpl struct {
	productRepo *repositories.ProductRepository
}

func NewProductService(productRepo *repositories.ProductRepository) *ProductServiceImpl {
	return &ProductServiceImpl{
		productRepo: productRepo,
	}
}

// --------------------
// CREATE PRODUCT
// --------------------
func (s *ProductServiceImpl) AddProduct(product *models.Product) (models.Product, error) {
	if product.ID.IsZero() {
		product.ID = primitive.NewObjectID()
	}
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	return s.productRepo.Create(*product)
}

// --------------------
// UPDATE PRODUCT (partial fields)
// --------------------
func (s *ProductServiceImpl) UpdateProduct(
	id string,
	update map[string]interface{},
) (*models.Product, error) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid product id")
	}

	if len(update) == 0 {
		return nil, errors.New("no fields provided for update")
	}

	update["updated_at"] = time.Now()
	return s.productRepo.UpdateFields(objID, update)
}

// --------------------
// DELETE PRODUCT
// --------------------
func (s *ProductServiceImpl) DeleteProduct(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid product id")
	}

	// 1️⃣ Fetch product first
	product, err := s.productRepo.FindByID(objID)
	if err != nil {
		return err
	}

	// 2️⃣ Delete image from Cloudinary (if exists)
	if product.ImageID != "" {
		if err := utils.DeleteImageFromCloudinary(product.ImageID); err != nil {
			return err
		}
	}

	// 3️⃣ Delete product from MongoDB
	return s.productRepo.Delete(objID)
}

// --------------------
// GET ALL PRODUCTS
// --------------------
func (s *ProductServiceImpl) GetAllProducts() ([]models.Product, error) {
	return s.productRepo.FindAll()
}

// --------------------
// GET PRODUCT BY ID
// --------------------
func (s *ProductServiceImpl) GetProductByID(id string) (*models.Product, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid product id")
	}
	return s.productRepo.FindByID(objID)
}
