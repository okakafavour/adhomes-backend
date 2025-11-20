package controllers

import (
	"adhomes-backend/models"
	"adhomes-backend/services"

	"net/http"

	"github.com/gin-gonic/gin"
)

var cartService services.CartService

func InitCartController() {
	cartService = services.NewCartService()
}

func CreateCart(c *gin.Context) {
	var cart models.Cart
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := cartService.CreateCart(cart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create cart"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "cart created",
		"cart":    created,
	})
}

func GetCart(c *gin.Context) {
	userID := c.Param("user_id")
	cart, err := cartService.GetCartByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cart not found"})
		return
	}
	c.JSON(http.StatusOK, cart)
}

func UpdateCart(c *gin.Context) {
	id := c.Param("id")
	var cart models.Cart
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := cartService.UpdateCart(id, cart)
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

func DeleteCart(c *gin.Context) {
	id := c.Param("id")
	err := cartService.DeleteCart(id)
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
