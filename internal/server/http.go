package server

import (
	"LocalHoist/utils"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"sync"

	"github.com/gin-gonic/gin"
)

// ReadAndPrepareRequest will read the whole request and convert in format to write to client
func ReadAndPrepareRequest(ctx *gin.Context) (*RequestPayload, []byte, error) {
	var bodyBytes []byte
	if ctx.Request.Body != nil {
		var err error
		bodyBytes, err = io.ReadAll(ctx.Request.Body)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to read body: %w", err)
		}
	}

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	payload := &RequestPayload{
		Method:  ctx.Request.Method,
		URL:     ctx.Request.URL.String(),
		Headers: ctx.Request.Header,
		Body:    bodyBytes,
	}

	rawHTTPBytes, err := httputil.DumpRequest(ctx.Request, true)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dump raw request: %w", err)
	}

	return payload, rawHTTPBytes, nil
}

// http://localhost:3000/tunnel/4cf2fd
// TunnelHandler will serve the client request
func (ws *WSManager) TunnelHandler(ctx *gin.Context) {
	payload, _, err := ReadAndPrepareRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("payload in Tunnelhandler is ", payload)

	// --- DEMONSTRATION OF READ DATA ---
	fmt.Printf("\n--- [Incoming %s Request] ---\n", payload.Method)
	fmt.Printf("URL Path & Query: %s\n", payload.URL)
	fmt.Printf("Headers Count: %d\n", len(payload.Headers))

	if len(payload.Body) > 0 {
		fmt.Printf("Body Content: %s\n", string(payload.Body))
	} else {
		fmt.Println("Body: (empty)")
	}

	url := payload.URL
	log.Println("ClientId in TunnelHandler is ", url)

	err, clientId := utils.GetClientIdFromRoute(url)
	if err != nil {
		log.Println("error in getting client id", err)
		return
	}
	log.Println("ClientId is ", clientId)

	err = ws.ForwardRequest(clientId, payload)
	if err != nil {
		log.Println("error in forwarding request ", err)
		return
	}

	// I need to capture the response from the backend running server and send it to frontend

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		// get the websocket connection object
		defer wg.Done()

		conn, ok := ws.wsConn[clientId]
		if !ok {
			// return valid response to client
			log.Println("No Id is present ")
			return
		}

		respPayload := ResponsePayload{}

		//_, msg, err := conn.ReadMessage()
		err = conn.ReadJSON(&respPayload)
		if err != nil {
			// return valid response to client
			log.Println("error in reading the response in waiting ", err)
			return
		}

		log.Println("Received response after hit is ", respPayload)

		ctx.JSON(http.StatusOK, respPayload)
	}()

	wg.Wait()
}
