package controllers

import (
	"adhomes-backend/models"
	"adhomes-backend/services"

	"net/http"

	"github.com/gin-gonic/gin"
)

type CartController struct {
	cartService services.CartService
}

func NewCartController(CartService services.CartService) *CartController {
	return &CartController{
		cartService: CartService,
	}

}

func (ca *CartController) CreateCart(c *gin.Context) {
	var cart models.Cart
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := ca.cartService.CreateCart(cart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create cart"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "cart created",
		"cart":    created,
	})
}

func (ca *CartController) GetCart(c *gin.Context) {
	userID := c.Param("user_id")
	cart, err := ca.cartService.GetCartByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cart not found"})
		return
	}
	c.JSON(http.StatusOK, cart)
}

func (ca *CartController) UpdateCart(c *gin.Context) {
	id := c.Param("id")
	var cart models.Cart
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := ca.cartService.UpdateCart(id, cart)
	if err != nil {
		if err.Error() == "cart not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "cart updated",
		"cart":    updated,
	})
}

func (ca *CartController) DeleteCart(c *gin.Context) {
	id := c.Param("id")
	err := ca.cartService.DeleteCart(id)
	if err != nil {
		if err.Error() == "cart not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err.Error() == "invalid cart ID" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cart deleted"})
}
