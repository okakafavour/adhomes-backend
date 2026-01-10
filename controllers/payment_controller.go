package controllers

import (
	"net/http"

	"adhomes-backend/models"
	"adhomes-backend/services"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	paymentService services.PaymentService
}

func NewPaymentController(paymentService services.PaymentService) *PaymentController {
	return &PaymentController{
		paymentService: paymentService,
	}
}

func (pc *PaymentController) MakePayment(c *gin.Context) {
	var req models.PaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, paymentURL, err := pc.paymentService.MakePayment(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"payment": payment,
	}

	if paymentURL != "" {
		response["payment_url"] = paymentURL
	}

	c.JSON(http.StatusCreated, response)
}
