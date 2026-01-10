package controllers

import (
	"adhomes-backend/models"
	"adhomes-backend/services"
	"adhomes-backend/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	productService services.ProductService
	orderServices  services.OrderService
	userServices   services.UserService
}

func NewAdminController(
	productService services.ProductService,
	orderServices services.OrderService,
	userServices services.UserService,
) *AdminController {
	return &AdminController{
		productService: productService,
		orderServices:  orderServices,
		userServices:   userServices,
	}
}

// === Admin Login ===
func (ac *AdminController) AdminLogin(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if req.Email != adminEmail || req.Password != adminPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid admin credentials"})
		return
	}

	token, err := utils.GenerateToken("admin", true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Admin login successful",
		"token":   token,
	})
}

// === Product Management ===
func (ac *AdminController) AddProduct(ctx *gin.Context) {
	var product models.Product

	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product data"})
		return
	}

	created, err := ac.productService.AddProduct(&product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"product": created,
	})
}

func (ac *AdminController) UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	var product models.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product data"})
		return
	}

	updated, err := ac.productService.UpdateProduct(id, &product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"product": updated,
	})
}

func (ac *AdminController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := ac.productService.DeleteProduct(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}

func (ac *AdminController) GetAllProducts(ctx *gin.Context) {
	products, err := ac.productService.GetAllProducts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}

func (ac *AdminController) GetProductByID(ctx *gin.Context) {
	id := ctx.Param("id")

	product, err := ac.productService.GetProductByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

// === Order Management ===
func (ac *AdminController) GetAllOrders(ctx *gin.Context) {
	orders, err := ac.orderServices.GetAllOrders()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}

func (ac *AdminController) ApproveOrder(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := ac.orderServices.ApproveOrder(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve order"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Order approved successfully",
	})
}

func (ac *AdminController) CancelOrder(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := ac.orderServices.CancelOrder(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel order"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Order canceled successfully",
	})
}

// === User Management ===
func (ac *AdminController) GetAllUsers(ctx *gin.Context) {
	users, err := ac.userServices.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func (ac *AdminController) DeactivateUser(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := ac.userServices.DeactivateUser(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate user"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User deactivated successfully",
	})
}

func (ac *AdminController) ActivateUser(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := ac.userServices.ActivateUser(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate user"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User activated successfully",
	})
}

func (ac *AdminController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := ac.userServices.DeleteUser(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}
