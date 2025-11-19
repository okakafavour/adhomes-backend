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

var orderCollection *mongo.Collection
var cartCollection *mongo.Collection

func InitOrderController() {
	orderCollection = config.DB.Collection("orders")
	cartCollection = config.DB.Collection("carts")
}

func CreateOrder(c *gin.Context) {
	var order models.Order

	// Bind JSON (expecting cart_id + delivery_address)
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate CartID
	cartID, err := primitive.ObjectIDFromHex(order.CartID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	// Fetch Cart from DB
	var cart models.Cart
	if err := cartCollection.FindOne(context.Background(), map[string]interface{}{"_id": cartID}).Decode(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart not found"})
		return
	}

	// Populate Order fields
	order.ID = primitive.NewObjectID()
	order.UserID = cart.UserID
	order.DeliveryAddress = order.DeliveryAddress

	// Convert Cart.Items -> Order.ProductIDs
	productIDs := make([]string, len(cart.Items))
	for i, item := range cart.Items {
		productIDs[i] = item.ProductID.Hex()
	}
	order.ProductIDs = productIDs

	order.PaymentStatus = "pending"
	order.OrderStatus = "Processing"
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	// Insert Order into DB
	_, err = orderCollection.InsertOne(context.Background(), order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
		"order":   order,
	})
}
