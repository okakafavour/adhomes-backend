package main

import (
	"adhomes-backend/config"
	"adhomes-backend/routes"
	"adhomes-backend/utils"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	godotenv.Load()

	// Connect to MongoDB
	config.ConnectDB()

	// Initialize Cloudinary
	if err := utils.InitCloudinary(); err != nil {
		log.Fatal(err)
	}

	// Create Gin router
	router := gin.Default()

	// Setup all routes
	routes.SetupRoutes(router)

	// Start server
	router.Run(":8080")
}
