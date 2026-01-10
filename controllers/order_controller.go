package controllers

import (
	"net/http"

	"adhomes-backend/models"
	"adhomes-backend/services"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	orderService services.OrderService
}

func NewOrderController(orderService services.OrderService) *OrderController {
	return &OrderController{
		orderService: orderService,
	}
}

// -----------------------------
// Create Order
// -----------------------------
func (oc *OrderController) CreateOrder(c *gin.Context) {
	var order models.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	createdOrder, err := oc.orderService.CreateOrder(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
func (oc *OrderController) GetOrderByID(c *gin.Context) {
	id := c.Param("id")

	order, err := oc.orderService.GetOrderByID(id)
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
func (oc *OrderController) GetOrdersByUserID(c *gin.Context) {
	userID := c.Param("user_id")

	orders, err := oc.orderService.GetOrdersByUserID(userID)
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
func (oc *OrderController) DeleteOrder(c *gin.Context) {
	id := c.Param("id")

	if err := oc.orderService.DeleteOrder(id); err != nil {
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
func (oc *OrderController) UpdateOrder(c *gin.Context) {
	id := c.Param("id")

	var input models.Order
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	updatedOrder, err := oc.orderService.UpdateOrder(id, input)
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
func (oc *OrderController) UpdateOrderStatus(c *gin.Context) {
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

	if err := oc.orderService.UpdateOrderStatus(id, body.Status); err != nil {
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
