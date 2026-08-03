package server

import (
	"LocalHoist/utils"
	"fmt"
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

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer conn.Close()

	clientID := utils.GenerateClientID()

	scheme := "https"
	if r.TLS == nil && r.Header.Get("X-Forwarded-Proto") != "https" {
		scheme = "http"
	}
	publicURL := fmt.Sprintf("%s://%s/tunnel/%s", scheme, r.Host, clientID)

	connResponse := WebSocketClientRequest{
		Message:   "proxy server connected",
		ErrorCode: 0,
		URL:       publicURL,
	}

	log.Println("Publice URL ",publicURL)

	err = conn.WriteJSON(connResponse)
	if err != nil {
		log.Println("Something went wrong")
		return

	}
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
