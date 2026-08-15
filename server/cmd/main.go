package main

import (
	server "LocalHoist-Server"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	wsManager := server.NewWSManager()

	r.GET("/ws", func(c *gin.Context) {
		wsManager.HandleWebSocket(c.Writer, c.Request)
	})
	r.Any("/tunnel/*filepath", wsManager.TunnelHandler)

	log.Println("Server is listening on :3000")

	if err := r.Run(":3000"); err != nil {
		log.Fatal("Server crashed:", err)
	}
}
