package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter creates and configures the Gin router with all routes
func SetupRouter() *gin.Engine {
	// Initialize a new Gin router
	router := gin.Default()

	// 1. Base route (GET /)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"route":   "/",
			"method":  "GET",
			"message": "Hello, World! The server is running.",
		})
	})

	// 2. Simple GET route (GET /a)
	router.GET("/a", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"route":   "/a",
			"method":  "GET",
			"message": "You hit route A!",
		})
	})

	// 3. Nested GET route (GET /a/b)
	router.GET("/a/b", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"route":   "/a/b",
			"method":  "GET",
			"message": "You hit the nested route A/B!",
		})
	})

	// 4. POST route that reads the request body (POST /a)
	// This is great for testing if your tunnel forwards JSON correctly!
	router.POST("/a", func(c *gin.Context) {
		// Create a map to hold the incoming JSON data
		var requestBody map[string]interface{}

		// Try to parse the incoming JSON
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to parse JSON body. Did you send valid JSON?",
			})
			return
		}

		// Echo the data back in the response
		c.JSON(http.StatusOK, gin.H{
			"route":         "/a",
			"method":        "POST",
			"message":       "Successfully received your data!",
			"received_data": requestBody,
		})
	})

	// 5. Nested POST route (POST /a/b)
	router.POST("/a/b", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"route":   "/a/b",
			"method":  "POST",
			"message": "POST request to nested route successful.",
		})
	})

	return router
}

func main() {
	router := SetupRouter()
	// Start the server on port 8080
	router.Run(":8080")
}