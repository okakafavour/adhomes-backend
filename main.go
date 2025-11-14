package main

import (
	"adhomes-backend/config"
	"adhomes-backend/controllers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config.ConnectDB()
	controllers.InitUserController()

	r := gin.Default()

	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)

	r.Run(":8080")
}
