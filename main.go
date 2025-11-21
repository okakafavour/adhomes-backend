package main

import (
	"adhomes-backend/config"
	"adhomes-backend/controllers"
	"adhomes-backend/services_impl"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	godotenv.Load()

	// Connect to MongoDB
	config.ConnectDB()

	// ------------------------------
	// Initialize services
	// ------------------------------
	paymentService := &services_impl.PaymentServiceImpl{
		PaymentCollection: config.DB.Collection("payments"),
		OrderCollection:   config.DB.Collection("orders"),
		HttpClient:        &http.Client{Timeout: 10 * time.Second},
	}

	walletService := services_impl.NewWalletService()

	paymentController := controllers.NewPaymentController(paymentService, walletService)

	r := gin.Default()

	controllers.InitUserController()

	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)

	r.POST("/payments", paymentController.MakePayment)

	r.Run(":8080")
}
