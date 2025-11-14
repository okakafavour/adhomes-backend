package controllers

import (
	"context"
	"net/http"
	"time"

	"adhomes-backend/config"
	"adhomes-backend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var CartCollection *mongo.Collection

func InitCartController() {
	CartCollection = config.GetCollection("carts")
}

func CreateCart(c *gin.Context) {
	collection := config.GetCollection("carts")

	var cart models.Cart
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	cart.ID = primitive.NewObjectID()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, cart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating cart"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Cart created successfully",
		"cart":    cart,
	})
}

func DeleteCart(c *gin.Context) {
	idParam := c.Param("id")
	cartID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := CartCollection.DeleteOne(ctx, primitive.M{"_id": cartID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting cart"})
		return
	}
	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cart deleted successfully",
	})
}

func UpdateCart(c *gin.Context) {
	idParam := c.Param("id")
	cartID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	var cart models.Cart
	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := primitive.M{
		"$set": cart,
	}

	result, err := CartCollection.UpdateOne(ctx, primitive.M{"_id": cartID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating cart"})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cart updated successfully",
	})
}
