package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"

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

func (ws *WSManager) TunnelHandler(ctx *gin.Context) {
	payload, rawBytes, err := ReadAndPrepareRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// --- DEMONSTRATION OF READ DATA ---
	fmt.Printf("\n--- [Incoming %s Request] ---\n", payload.Method)
	fmt.Printf("URL Path & Query: %s\n", payload.URL)
	fmt.Printf("Headers Count: %d\n", len(payload.Headers))

	if len(payload.Body) > 0 {
		fmt.Printf("Body Content: %s\n", string(payload.Body))
	} else {
		fmt.Println("Body: (empty)")
	}

	clientId := payload.URL
	// Example: Serialize structured payload to JSON frame to send across your WebSocket/TCP channel
	jsonFrame, err := json.Marshal(payload)
	if err != nil {
		log.Println("error in marshaling to json ", err)
		return
	}

	err = ws.ForwardRequest(clientId, jsonFrame)
	if err != nil {
		log.Println("error in forwarding request ", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":         "Request captured",
		"method":         payload.Method,
		"raw_bytes_len":  len(rawBytes),
		"json_frame_len": len(jsonFrame),
	})
}
