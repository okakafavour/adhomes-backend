package controllers

import (
	"net/http"

	"adhomes-backend/models"
	"adhomes-backend/services"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	paymentService services.PaymentService
	walletService  services.WalletService
}

func NewPaymentController(paymentService services.PaymentService, walletService services.WalletService) *PaymentController {
	return &PaymentController{
		paymentService: paymentService,
		walletService:  walletService,
	}
}

func (pc *PaymentController) MakePayment(ctx *gin.Context) {
	var req models.PaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	// Optional: Deduct wallet if available
	wallet, err := pc.walletService.GetWalletByUserID(ctx, req.UserID)
	if err == nil && wallet.Balance >= req.Amount {
		_, err = pc.walletService.DecreaseBalance(ctx, req.UserID, req.Amount)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to deduct wallet balance",
				"error":   err.Error(),
			})
			return
		}
	}

	// Initialize Paystack payment
	payment, paymentURL, err := pc.paymentService.InitializePayment(
		req.OrderID,
		req.UserID,
		req.Amount,
		req.Email,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to initialize payment",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"payment":     payment,
		"payment_url": paymentURL,
	})
}
