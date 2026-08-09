// main.go
package main

import (
	"encoding/json"
	"io"
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

	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("Server disconnected or read error:", err)
		return
	}

	log.Println("Request in clicllclclnet ", string(message))

	payload, err := DecodeRequest(message)
	if err != nil {
		log.Println("Error decoding request:", err)
		return
	}

	log.Println("LocalHoist URL is ", payload.URL)

	// Start a background Goroutine to continuously read messages
	go func() {
		for {
			payload := RequestPayload{}
			// _, message, err := conn.ReadMessage()

			err := conn.ReadJSON(&payload)
			if err != nil {
				log.Println("Server disconnected or read error:", err)
				return
			}

			log.Println("Request in clinet ", payload)

			log.Println("Reeceiving message is 1111")

			// Execute local HTTP request
			resp, err := Execute(&payload)
			if err != nil {
				log.Println("Error making request to local server:", err)
				// Optional: send error payload back to WebSocket server here
				continue
			}

			// Read the EXACT body bytes (handles JSON, HTML, images, binary streams, etc.)
			bodyBytes, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				log.Println("Error reading local response body:", err)
				continue
			}

			var responseBody map[string]interface{}

			err=json.Unmarshal(bodyBytes,&responseBody)
			if err!=nil{
				log.Println("error in unmarshalling json ",err)
				return
			}

			// Construct response with exact headers and raw bytes
			respPayload := ResponsePayload{
				StatusCode: resp.StatusCode,
				Headers:    resp.Header,
				Body:       responseBody,
			}

			// Send back over WebSocket
			if err := conn.WriteJSON(respPayload); err != nil {
				log.Println("Error writing response to server:", err)
				return
			}
			log.Println("Written back to socket is ",respPayload)
		}
	}()

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
