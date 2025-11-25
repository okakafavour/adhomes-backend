package controllers

import (
	"net/http"

	"adhomes-backend/services"
	"adhomes-backend/services_impl"

	"github.com/gin-gonic/gin"
)

type DeliveryController struct {
	deliveryService services.DeliveryService
}

// global controller instance for test setup
var deliveryController *DeliveryController

// Constructor
func NewDeliveryController(deliveryService services.DeliveryService) *DeliveryController {
	return &DeliveryController{
		deliveryService: deliveryService,
	}
}

// Initialize the global controller (for router setup in tests)
func InitDeliveryController() {
	deliveryService := services_impl.NewDeliveryService()
	deliveryController = NewDeliveryController(deliveryService)
}

// POST /deliveries/assign
func (dc *DeliveryController) AssignRider(ctx *gin.Context) {
	var req struct {
		OrderID string `json:"order_id"`
		RiderID string `json:"rider_id"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	delivery, err := dc.deliveryService.AssignRider(req.OrderID, req.RiderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":  "delivery assigned",
		"delivery": delivery,
	})
}

// PUT /deliveries/:id/status
func (dc *DeliveryController) UpdateStatus(ctx *gin.Context) {
	deliveryID := ctx.Param("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	delivery, err := dc.deliveryService.UpdateDeliveryStatus(deliveryID, req.Status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "status updated",
		"delivery": delivery,
	})
}

// GET /deliveries/order/:order_id
func (dc *DeliveryController) GetByOrder(ctx *gin.Context) {
	orderID := ctx.Param("order_id")
	delivery, err := dc.deliveryService.GetDeliveryByOrder(orderID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"delivery": delivery,
	})
}
