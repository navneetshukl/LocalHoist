package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize a new Gin router with default middleware (logger and recovery)
	router := gin.Default()

	// Define a simple GET endpoint at the root path "/"
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
			"status":  "success",
		})
	})

	// Start the server on port 8080
	router.Run(":8080")
}