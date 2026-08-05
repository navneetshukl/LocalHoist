// main.go
package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "LocalHoist",
		Short: "LocalHoist creates secure tunnels to your local ports",
	}

	var connectCmd = &cobra.Command{
		// [url] in the Use string tells users an optional argument is expected
		Use:   "connect [url]",
		Short: "Connect to the WebSocket server",
		// This ensures the user passes at most 1 argument (the URL)
		Args: cobra.MaximumNArgs(1),
		Run:  runConnect,
	}

	rootCmd.AddCommand(connectCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

// Notice the args []string parameter
func runConnect(cmd *cobra.Command, args []string) {
	// Set the default URL
	serverURL := "ws://localhost:3000/ws"

	// If the user passed an argument (like `mkdir name`), override the default URL
	if len(args) > 0 {
		serverURL = args[0]
	}

	log.Printf("Connecting to %s...", serverURL)

	// Dial the server
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer conn.Close()

	log.Println("Connected to Server! Tunnel is open.")

	// Start a background Goroutine to continuously read messages
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Server disconnected or read error:", err)
				return
			}
			// convert this to http request payload

			payload, err := DecodeRequest(message)
			if err != nil {
				log.Println("error in decoding request ", err)
				return
			}

			// make the request to local server

			resp, err := Execute(payload)
			if err != nil {
				log.Println("error in making request to local", err)
				return
			}

			err = conn.WriteJSON(resp)
			if err != nil {
				log.Println("error in writing to server", err)
				return
			}

			log.Printf("Message from Server: %s\n", message)
		}
	}()

	// Send a test message up the tunnel
	testMessage := []byte("Hello from the Local Agent!")
	err = conn.WriteMessage(websocket.TextMessage, testMessage)
	if err != nil {
		log.Println("Failed to send message:", err)
	}

	// Block the main thread so the program doesn't exit immediately
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	log.Println("\nShutting down Agent...")

	err = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("Error during closing handshake:", err)
	}
}
