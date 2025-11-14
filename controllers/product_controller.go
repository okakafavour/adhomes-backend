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

var productCollection *mongo.Collection

func InitProductController() {
	productCollection = config.GetCollection("products")
}

func CreateProduct(c *gin.Context) {
	var product models.Product

	// ✅ Properly bind and validate request
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// ✅ Assign new ID and timestamp before saving
	product.ID = primitive.NewObjectID()
	product.CreatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ✅ Insert product into MongoDB
	_, err := productCollection.InsertOne(ctx, product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating product"})
		return
	}

	// ✅ Respond with success message and data
	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"product": product,
	})
}

func DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	productID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := productCollection.DeleteOne(ctx, gin.H{"_id": productID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting product"})
		return
	}
	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	productID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := gin.H{
		"$set": gin.H{
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"category":    product.Category,
			"image_url":   product.ImageURL,
		},
	}

	result, err := productCollection.UpdateOne(ctx, gin.H{"_id": productID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating product"})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"product": product,
	})
}
