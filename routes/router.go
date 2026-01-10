package routes

import (
	"adhomes-backend/controllers"
	"adhomes-backend/middleware"
	"adhomes-backend/services_impl"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// ==========================
	// SERVICES
	// ==========================
	userService := services_impl.NewUserService()
	productService := services_impl.NewProductService()
	orderService := services_impl.NewOrderService()
	favouriteService := services_impl.NewFavoriteService()
	paymentService := services_impl.NewPaymentService() // Payment service

	// ==========================
	// CONTROLLERS
	// ==========================
	userController := controllers.NewUserController(userService)
	productController := controllers.NewProductController(productService)
	orderController := controllers.NewOrderController(orderService)
	favouriteController := controllers.NewFavoriteController(favouriteService)
	paymentController := controllers.NewPaymentController(paymentService)

	adminController := controllers.NewAdminController(
		productService,
		orderService,
		userService,
	)

	// ==========================
	// PUBLIC AUTH ROUTES
	// ==========================
	r.POST("/signup", userController.SignUp)
	r.POST("/login", userController.Login)
	r.POST("/admin/login", adminController.AdminLogin)

	// ==========================
	// PUBLIC PRODUCT ROUTES
	// ==========================
	r.GET("/products", productController.GetAllProducts)
	r.GET("/products/:id", productController.GetProductByID)

	// ==========================
	// USER ROUTES (JWT PROTECTED)
	// ==========================
	userRoutes := r.Group("/user")
	userRoutes.Use(middleware.AuthMiddleware())
	{
		// Orders
		userRoutes.POST("/orders", orderController.CreateOrder)
		userRoutes.GET("/orders/:id", orderController.GetOrderByID)
		userRoutes.GET("/orders", orderController.GetOrdersByUserID)
		userRoutes.DELETE("/orders/:id", orderController.DeleteOrder)
		userRoutes.PUT("/orders/:id", orderController.UpdateOrder)
		userRoutes.PUT("/orders/:id/status", orderController.UpdateOrderStatus)

		// Favourites
		userRoutes.POST("/favourite", favouriteController.AddFavorite)
		userRoutes.GET("/favourite", favouriteController.GetFavorites)
		userRoutes.DELETE("/favourite/:id", favouriteController.RemoveFavorite)

		// Payments
		userRoutes.POST("/payments", paymentController.MakePayment)
	}

	// ==========================
	// ADMIN ROUTES (JWT + Admin Middleware)
	// ==========================
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminAuth())
	{
		// Product Management
		admin.POST("/products", adminController.AddProduct)
		admin.PUT("/products/:id", adminController.UpdateProduct)
		admin.DELETE("/products/:id", adminController.DeleteProduct)
		admin.GET("/products", adminController.GetAllProducts)
		admin.GET("/products/:id", adminController.GetProductByID)

		// Order Management
		admin.GET("/orders", adminController.GetAllOrders)
		admin.PUT("/orders/:id/approve", adminController.ApproveOrder)
		admin.PUT("/orders/:id/cancel", adminController.CancelOrder)

		// User Management
		admin.GET("/users", adminController.GetAllUsers)
		admin.PUT("/users/:id/deactivate", adminController.DeactivateUser)
		admin.PUT("/users/:id/activate", adminController.ActivateUser)
		admin.DELETE("/users/:id", adminController.DeleteUser)
	}
}
