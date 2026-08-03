// cmd/portbridge-agent/main.go
package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

func main() {
	// The URL of your Server
	url := "ws://localhost:3000/ws"
	log.Printf("Connecting to %s...", url)

	// 1. Dial the server
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer conn.Close()

	log.Println("Connected to Server! Tunnel is open.")

	// 2. Start a background Goroutine to continuously read messages
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Server disconnected or read error:", err)
				return
			}
			log.Printf("Message from Server: %s\n", message)
		}
	}()

	// 3. Send a test message up the tunnel
	testMessage := []byte("Hello from the Local Agent!")
	err = conn.WriteMessage(websocket.TextMessage, testMessage)
	if err != nil {
		log.Println("Failed to send message:", err)
	}

	// 4. Block the main thread so the program doesn't exit immediately
	// This waits for you to press Ctrl+C
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	log.Println("Shutting down Agent...")
}