package controllers

import (
	"adhomes-backend/models"
	"adhomes-backend/services"
	"adhomes-backend/services_impl"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var orderService services.OrderService

func InitOrderController() {
	orderService = services_impl.NewOrderService()
}

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

func GetOrderByID(c *gin.Context) {
	id := c.Param("id")

	order, err := orderService.GetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func GetOrdersByUserID(c *gin.Context) {
	userID := c.Param("user_id")

	orders, err := orderService.GetOrdersByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func DeleteOrder(c *gin.Context) {
	id := c.Param("id")

	err := orderService.DeleteOrder(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order deleted"})
}

func UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		PaymentStatus string `json:"payment_status"`
		OrderStatus   string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Fetch existing order
	order, err := orderService.GetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	// Update fields
	if input.PaymentStatus != "" {
		order.PaymentStatus = input.PaymentStatus
	}
	if input.OrderStatus != "" {
		order.OrderStatus = input.OrderStatus
	}
	order.UpdatedAt = time.Now()

	updatedOrder, err := orderService.UpdateOrder(id, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "order updated",
		"order":   updatedOrder,
	})
}

func UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var body struct {
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Call the service implementation
	err := orderService.UpdateOrderStatus(id, body.Status)
	if err != nil {
		switch err.Error() {
		case "invalid order ID":
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
			return
		case "invalid order status":
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order status"})
			return
		case "order not found":
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Order status updated successfully",
		"new_status": body.Status,
	})
}
