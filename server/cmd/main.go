package main

import (
	server "LocalHoist-Server"
	"log"
	"os" // Required to read environment variables

	"github.com/gin-gonic/gin"
)

func main() {
	// Read the PORT from Cloud Run, default to 8080 if running locally
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Set Gin to Release Mode for production
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	
	// Trust Cloud Run's proxy to avoid warnings
	_ = r.SetTrustedProxies(nil)

	wsManager := server.NewWSManager()

	r.GET("/ws", func(c *gin.Context) {
		wsManager.HandleWebSocket(c.Writer, c.Request)
	})
	r.Any("/tunnel/*filepath", wsManager.TunnelHandler)

	log.Printf("Server is listening on :%s\n", port)

	// Bind to the dynamic port here instead of ":3000"
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Server crashed:", err)
	}
}