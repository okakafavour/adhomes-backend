package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role") // Extracted from JWT in your auth middleware

		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Admin only."})
			c.Abort()
			return
		}

		c.Next()
	}
}
