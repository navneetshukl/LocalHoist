package main

import (
	server "LocalHoist-Server"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Read the PORT environment variable set by Cloud Run (fallback to 8080 for local testing)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Set Gin to release mode in production
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// Trust Cloud Run / Google Cloud front-end proxy
	_ = r.SetTrustedProxies(nil)

	wsManager := server.NewWSManager()

	r.GET("/ws", func(c *gin.Context) {
		wsManager.HandleWebSocket(c.Writer, c.Request)
	})
	r.Any("/tunnel/*filepath", wsManager.TunnelHandler)

	log.Printf("Server is listening on :%s\n", port)

	// Bind dynamically to ":PORT"
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Server crashed:", err)
	}
}