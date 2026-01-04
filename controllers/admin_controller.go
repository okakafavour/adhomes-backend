package controllers

import (
	"net/http"

	"adhomes-backend/models"
	"adhomes-backend/services"
	"adhomes-backend/services_impl"

	"github.com/gin-gonic/gin"
)

// -----------------------
// Admin Controller
// -----------------------
type AdminController struct {
	productService services.ProductService
	orderService   services.OrderService
}

func AdminControllerSingleton() *AdminController {
	if adminController == nil {
		InitAdminController()
	}
	return adminController
}

var adminController *AdminController

func NewAdminController(productService services.ProductService, orderService services.OrderService) *AdminController {
	return &AdminController{
		productService: productService,
		orderService:   orderService,
	}
}

func InitAdminController() {
	productService := services_impl.NewProductService()
	orderService := services_impl.NewOrderService()
	adminController = NewAdminController(productService, orderService)
}

// -----------------------
// Product Endpoints
// -----------------------

func (ac *AdminController) AddProduct(ctx *gin.Context) {
	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Category    string  `json:"category"`
		ImageURL    string  `json:"image_url"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := ac.productService.AddProduct(req.Name, req.Description, req.Price, req.Category, req.ImageURL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Product added", "product": product})
}

func (ac *AdminController) UpdateProduct(ctx *gin.Context) {
	productID := ctx.Param("id")

	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Category    string  `json:"category"`
		ImageURL    string  `json:"image_url"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// FIXED: must pass a models.Product
	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
	}

	updatedProduct, err := ac.productService.UpdateProduct(productID, product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product updated", "product": updatedProduct})
}

func (ac *AdminController) DeleteProduct(ctx *gin.Context) {
	productID := ctx.Param("id")

	err := ac.productService.DeleteProduct(productID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

func (ac *AdminController) GetAllProducts(ctx *gin.Context) {
	products, err := ac.productService.GetAllProducts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"products": products})
}

// -----------------------
// Order Endpoints
// -----------------------

func (ac *AdminController) GetAllOrders(ctx *gin.Context) {
	// FIXED: This function must EXIST in OrderService
	orders, err := ac.orderService.GetAllOrders()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"orders": orders})
}

func (ac *AdminController) UpdateOrderStatus(ctx *gin.Context) {
	orderID := ctx.Param("id")

	var req struct {
		Status string `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// FIXED: orderService.UpdateOrderStatus returns ONLY error
	err := ac.orderService.UpdateOrderStatus(orderID, req.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order status updated"})
}
