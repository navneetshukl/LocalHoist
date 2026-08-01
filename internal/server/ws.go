package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer conn.Close()

	ctx := 0

	for {
		log.Println("inside the request ", ctx)
		messageType, payload, err := conn.ReadMessage()
		if err != nil {
			log.Println("Agent disconnected or read error:", err)
			break
		}
		log.Println("outside the request ", ctx)

		log.Printf("Received from Agent: %s\n", payload)
		response := []byte("Server received: " + string(payload))
		err = conn.WriteMessage(messageType, response)
		if err != nil {
			log.Println("Write error:", err)
			break
		}
		ctx++
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)

	log.Println("Server is listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server crashed:", err)
	}
}
