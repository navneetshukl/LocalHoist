/**

This file will open the websocker connection when first time client will connect than generate the public url
and send the public url to the client.

**/

package server

import (
	"LocalHoist-Server/utils"
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

type WSManager struct {
	wsConn map[string]*websocket.Conn
}

func NewWSManager() *WSManager {
	return &WSManager{
		wsConn: make(map[string]*websocket.Conn),
	}
}

// HandleWebSocket will start the websocket connection
func (ws *WSManager) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer conn.Close()

	clientID := utils.GenerateClientID()

	ws.wsConn[clientID] = conn

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

	log.Println("Publice URL ", publicURL)

	err = conn.WriteJSON(connResponse)
	if err != nil {
		log.Println("Something went wrong")
		return

	}

	select {}
}

// ForwardRequest will forward the request from the frontend to locally running backend
func (ws *WSManager) ForwardRequest(clientId string, payload interface{}) error {

	// get the websocket object

	log.Println("Inside forward request.Payload is ", payload)

	c, ok := ws.wsConn[clientId]
	if !ok {
		return fmt.Errorf("no connection for this %s client \n", clientId)
	}
	err := c.WriteJSON(payload)
	if err != nil {
		return err
	}
	return nil

}
