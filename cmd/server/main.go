package main

import (
	"LocalHoist/internal/server"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	wsManager := server.NewWSManager()
	http.HandleFunc("/ws", wsManager.HandleWebSocket)
	r.Any("/tunnel/*filepath", wsManager.TunnelHandler)

	log.Println("Server is listening on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal("Server crashed:", err)
	}
}
