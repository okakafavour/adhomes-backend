package routes

import (
	"adhomes-backend/config"
	"adhomes-backend/controllers"
	"adhomes-backend/services_impl"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// -----------------------------
	// Initialize all controllers
	// -----------------------------
	controllers.InitUserController()
	controllers.InitAdminController()
	controllers.InitCartController()
	controllers.InitDeliveryController()
	controllers.InitFavouriteController()
	controllers.InitOrderController()
	controllers.InitProductController()

	// Payment Service & Controller
	paymentService := &services_impl.PaymentServiceImpl{
		PaymentCollection: config.DB.Collection("payments"),
		OrderCollection:   config.DB.Collection("orders"),
		HttpClient:        &http.Client{Timeout: 10 * time.Second},
	}
	walletService := services_impl.NewWalletService()
	paymentController := controllers.NewPaymentController(paymentService, walletService)

	// -----------------------------
	// Public / User routes
	// -----------------------------
	router.POST("/signup", controllers.SignUp)
	router.POST("/login", controllers.Login)

	router.GET("/products", controllers.GetAllProducts)
	router.POST("/payments", paymentController.MakePayment)

	// Cart routes
	cart := controllers.CartControllerSingleton()
	cartRoutes := router.Group("/cart")
	{
		cartRoutes.POST("/", cart.AddToCart)
		cartRoutes.GET("/", cart.GetCartItems)
		cartRoutes.DELETE("/:id", cart.RemoveFromCart)
	}

	// Favourite routes
	fav := controllers.FavouriteControllerSingleton()
	favRoutes := router.Group("/favourites")
	{
		favRoutes.POST("/", fav.AddToFavourites)
		favRoutes.GET("/", fav.GetFavourites)
		favRoutes.DELETE("/:id", fav.RemoveFavourite)
	}

	// Delivery routes
	delivery := controllers.DeliveryControllerSingleton()
	deliveryRoutes := router.Group("/delivery")
	{
		deliveryRoutes.POST("/", delivery.ScheduleDelivery)
		deliveryRoutes.GET("/", delivery.GetAllDeliveries)
	}

	// -----------------------------
	// Admin routes
	// -----------------------------
	admin := controllers.AdminControllerSingleton()
	adminRoutes := router.Group("/admin")
	{
		// Product management
		adminRoutes.POST("/products", admin.AddProduct)
		adminRoutes.PUT("/products/:id", admin.UpdateProduct)
		adminRoutes.DELETE("/products/:id", admin.DeleteProduct)
		adminRoutes.GET("/products", admin.GetAllProducts)

		// Order management
		adminRoutes.GET("/orders", admin.GetAllOrders)
		adminRoutes.PUT("/orders/:id/status", admin.UpdateOrderStatus)
	}
}
