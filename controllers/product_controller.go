package controllers

import (
	"net/http"

	"adhomes-backend/models"
	"adhomes-backend/services"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	productService services.ProductService
}

func NewProductController(productService services.ProductService) *ProductController {
	return &ProductController{
		productService: productService,
	}
}

// --------------------
// CREATE PRODUCT (ADMIN)
// --------------------
func (pc *ProductController) CreateProduct(c *gin.Context) {
	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	createdProduct, err := pc.productService.AddProduct(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create product",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"product": createdProduct,
	})
}

// --------------------
// UPDATE PRODUCT
// --------------------
func (pc *ProductController) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	if len(input) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No fields provided for update",
		})
		return
	}

	updatedProduct, err := pc.productService.UpdateProduct(id, input)
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update product",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"product": updatedProduct,
	})
}

// --------------------
// DELETE PRODUCT
// --------------------
func (pc *ProductController) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	if err := pc.productService.DeleteProduct(id); err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete product",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}

// --------------------
// GET ALL PRODUCTS (PUBLIC)
// --------------------
func (pc *ProductController) GetAllProducts(c *gin.Context) {
	products, err := pc.productService.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch products",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}

// --------------------
// GET PRODUCT BY ID (PUBLIC)
// --------------------
func (pc *ProductController) GetProductByID(c *gin.Context) {
	id := c.Param("id")

	product, err := pc.productService.GetProductByID(id)
	if err != nil {
		if err.Error() == "product not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch product",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}
