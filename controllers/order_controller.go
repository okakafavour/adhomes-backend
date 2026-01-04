package controllers

import (
	"net/http"

	"adhomes-backend/models"
	"adhomes-backend/services"
	"adhomes-backend/services_impl"

	"github.com/gin-gonic/gin"
)

var orderService services.OrderService

// -----------------------------
// Init
// -----------------------------
func InitOrderController() {
	orderService = services_impl.NewOrderService()
}

// -----------------------------
// Create Order
// -----------------------------
func CreateOrder(c *gin.Context) {
	var order models.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	createdOrder, err := orderService.CreateOrder(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
		"order":   createdOrder,
	})
}

// -----------------------------
// Get Order By ID
// -----------------------------
func GetOrderByID(c *gin.Context) {
	id := c.Param("id")

	order, err := orderService.GetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, order)
}

// -----------------------------
// Get Orders By User ID
// -----------------------------
func GetOrdersByUserID(c *gin.Context) {
	userID := c.Param("user_id")

	orders, err := orderService.GetOrdersByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// -----------------------------
// Delete Order
// -----------------------------
func DeleteOrder(c *gin.Context) {
	id := c.Param("id")

	if err := orderService.DeleteOrder(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order deleted successfully",
	})
}

// -----------------------------
// Update Order (full update)
// -----------------------------
func UpdateOrder(c *gin.Context) {
	id := c.Param("id")

	var input models.Order
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	updatedOrder, err := orderService.UpdateOrder(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order updated successfully",
		"order":   updatedOrder,
	})
}

// -----------------------------
// Update Order Status ONLY
// -----------------------------
func UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var body struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Status is required",
		})
		return
	}

	if err := orderService.UpdateOrderStatus(id, body.Status); err != nil {
		switch err.Error() {
		case "invalid order ID":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "invalid order status":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "order not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order status updated successfully",
		"status":  body.Status,
	})
}
