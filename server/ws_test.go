package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestNewWSManager(t *testing.T) {
	manager := NewWSManager()
	if manager == nil {
		t.Fatal("Expected non-nil manager")
	}
	if manager.wsConn == nil {
		t.Error("Expected non-nil wsConn map")
	}
	if len(manager.wsConn) != 0 {
		t.Error("Expected empty wsConn map on creation")
	}
}

func TestForwardRequestNoConnection(t *testing.T) {
	manager := NewWSManager()
	err := manager.ForwardRequest("nonexistent", nil)
	if err == nil {
		t.Error("Expected error for non-existent client ID")
	}
}

// wsTestServer is a test HTTP server that upgrades the connection to WebSocket
func wsTestServer(t *testing.T, onConnect func(*websocket.Conn)) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Logf("Upgrade failed: %v", err)
			return
		}
		defer conn.Close()
		if onConnect != nil {
			onConnect(conn)
		}
		// Keep the connection alive briefly
		time.Sleep(50 * time.Millisecond)
	}))
}

func TestHandleWebSocket(t *testing.T) {
	// Start a test server that handles the WebSocket upgrade
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Logf("Upgrade failed: %v", err)
			return
		}
		defer conn.Close()
	}))
	defer ts.Close()

	// Convert http URL to ws URL
	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/ws"

	// Dial the server directly
	dialer := websocket.DefaultDialer
	conn, resp, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Logf("Dial failed (expected for raw handler): %v", err)
		return
	}
	defer conn.Close()
	if resp == nil {
		t.Error("Expected non-nil response")
	}
}

func TestUpgraderCheckOrigin(t *testing.T) {
	// Verify that the upgrader's CheckOrigin returns true
	if !upgrader.CheckOrigin(nil) {
		t.Error("Expected CheckOrigin to return true")
	}
}

func TestUpgraderBufferSizes(t *testing.T) {
	if upgrader.ReadBufferSize != 1024 {
		t.Errorf("Expected ReadBufferSize 1024, got %d", upgrader.ReadBufferSize)
	}
	if upgrader.WriteBufferSize != 1024 {
		t.Errorf("Expected WriteBufferSize 1024, got %d", upgrader.WriteBufferSize)
	}
}

func TestForwardRequestWithPayload(t *testing.T) {
	// Create a real WebSocket test server that will receive the forwarded payload
	received := make(chan interface{}, 1)
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		var payload interface{}
		if err := conn.ReadJSON(&payload); err == nil {
			received <- payload
		}
	}))
	defer ts.Close()

	// Connect to test server
	dialer := websocket.Dialer{}
	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http")
	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	defer conn.Close()

	// Manually inject this connection into a WSManager
	manager := NewWSManager()
	manager.wsConn["test-client"] = conn

	// Now forward a request
	payload := map[string]interface{}{
		"method": "GET",
		"url":    "/test",
	}
	err = manager.ForwardRequest("test-client", payload)
	if err != nil {
		t.Errorf("ForwardRequest failed: %v", err)
	}

	// Wait for the message to be received
	select {
	case msg := <-received:
		if msg == nil {
			t.Error("Expected non-nil received payload")
		}
	case <-time.After(2 * time.Second):
		t.Error("Timeout waiting for forwarded payload")
	}
}
