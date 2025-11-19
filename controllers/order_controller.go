package controllers

import (
	"adhomes-backend/models"
	"adhomes-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

var orderService services.OrderService

// Initialize the order controller with service
func InitOrderController() {
	orderService = services.NewOrderService()
}

// -------------------------------
// CREATE ORDER
// -------------------------------
func CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	created, err := orderService.CreateOrder(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "order created",
		"order":   created,
	})
}

// -------------------------------
// GET ORDER BY ID
// -------------------------------
func GetOrderByID(c *gin.Context) {
	id := c.Param("id")

	order, err := orderService.GetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// -------------------------------
// GET ORDERS BY USER ID
// -------------------------------
func GetOrdersByUserID(c *gin.Context) {
	userID := c.Param("user_id")

	orders, err := orderService.GetOrdersByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// -------------------------------
// DELETE ORDER
// -------------------------------
func DeleteOrder(c *gin.Context) {
	id := c.Param("id")

	err := orderService.DeleteOrder(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order deleted"})
}
