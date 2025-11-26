package controllers

import (
	"net/http"

	"adhomes-backend/models"
	"adhomes-backend/services"
	"adhomes-backend/services_impl"

	"github.com/gin-gonic/gin"
)

var productService services.ProductService

func InitProductController() {
	productService = services_impl.NewProductService()
}

// CreateProduct (Admin)
func CreateProduct(c *gin.Context) {
	var body struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Category    string  `json:"category"`
		ImageURL    string  `json:"image_url"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	p, err := productService.AddProduct(
		body.Name,
		body.Description,
		body.Price,
		body.Category,
		body.ImageURL,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"product": p,
	})
}

// Update Product
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	p, err := productService.UpdateProduct(id, product)
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "product updated",
		"product": p,
	})
}

// Delete Product
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	err := productService.DeleteProduct(id)
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}

// Get All Products
func GetAllProducts(c *gin.Context) {
	products, err := productService.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}
