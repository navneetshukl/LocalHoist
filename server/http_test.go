package server

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func TestReadAndPrepareRequestGET(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/test?key=value", nil)

	payload, raw, err := ReadAndPrepareRequest(c)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if payload == nil {
		t.Fatal("Expected payload")
	}
	if payload.Method != "GET" {
		t.Errorf("Expected GET, got %s", payload.Method)
	}
	if raw == nil {
		t.Error("Expected raw bytes")
	}
}

func TestReadAndPrepareRequestPOST(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"test":"data"}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	payload, _, err := ReadAndPrepareRequest(c)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if payload.Method != "POST" {
		t.Errorf("Expected POST, got %s", payload.Method)
	}
	if len(payload.Body) == 0 {
		t.Error("Expected body to be read")
	}
}

func TestReadAndPrepareRequestNoBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("DELETE", "/api/1", nil)

	payload, _, err := ReadAndPrepareRequest(c)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if len(payload.Body) != 0 {
		t.Error("Expected empty body")
	}
}

func TestReadAndPrepareRequestWithHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Custom", "test")
	c.Request = req

	payload, _, err := ReadAndPrepareRequest(c)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if len(payload.Headers["X-Custom"]) == 0 || payload.Headers["X-Custom"][0] != "test" {
		t.Error("Headers not preserved")
	}
}

func TestTunnelHandlerNoClient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/tunnel/nonexistent/api", nil)

	manager := NewWSManager()
	manager.TunnelHandler(c)
}

func TestReadAndPrepareRequestPUT(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"name":"updated"}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/api/1", bytes.NewReader(body))

	payload, raw, err := ReadAndPrepareRequest(c)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if payload.Method != "PUT" {
		t.Errorf("Expected PUT, got %s", payload.Method)
	}
	if !bytes.Contains(payload.Body, []byte(`"name":"updated"`)) {
		t.Error("Body content mismatch")
	}
	if raw == nil {
		t.Error("Expected raw bytes")
	}
}

func TestReadAndPrepareRequestURLPreservation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api?query=123&sort=name", nil)

	payload, _, err := ReadAndPrepareRequest(c)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if payload.URL != "/api?query=123&sort=name" {
		t.Errorf("Expected URL '/api?query=123&sort=name', got '%s'", payload.URL)
	}
}

func TestReadAndPrepareRequestBodyReusability(t *testing.T) {
	gin.SetMode(gin.TestMode)
	body := []byte(`{"test":"data"}`)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api", bytes.NewReader(body))

	payload, _, err := ReadAndPrepareRequest(c)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	bodyRead2, err := io.ReadAll(c.Request.Body)
	if err != nil {
		t.Fatalf("Failed to re-read body: %v", err)
	}
	if !bytes.Equal(payload.Body, bodyRead2) {
		t.Error("Body should be reusable after ReadAndPrepareRequest")
	}
}

func TestTunnelHandlerSuccessful(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		var payload RequestPayload
		if err := conn.ReadJSON(&payload); err != nil {
			return
		}

		resp := ResponsePayload{
			StatusCode: 200,
			Headers:    map[string][]string{"Content-Type": {"application/json"}},
			Body:       map[string]interface{}{"status": "ok"},
		}
		if err := conn.WriteJSON(resp); err != nil {
			return
		}
	}))
	defer ts.Close()

	wsURL := "ws" + ts.URL[4:] + "/ws"
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	defer conn.Close()

	manager := NewWSManager()
	manager.wsConn["abc123"] = conn

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/tunnel/abc123/api", nil)

	manager.TunnelHandler(c)

	time.Sleep(100 * time.Millisecond)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestTunnelHandlerWithBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		var payload RequestPayload
		if err := conn.ReadJSON(&payload); err == nil {
			resp := ResponsePayload{
				StatusCode: 201,
				Headers:    map[string][]string{},
				Body:       map[string]interface{}{"received": string(payload.Body)},
			}
			conn.WriteJSON(resp)
		}
	}))
	defer ts.Close()

	wsURL := "ws" + ts.URL[4:] + "/ws"
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	defer conn.Close()

	manager := NewWSManager()
	manager.wsConn["postClient"] = conn

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := bytes.NewReader([]byte(`{"input":"data"}`))
	c.Request = httptest.NewRequest("POST", "/tunnel/postClient/submit", body)

	manager.TunnelHandler(c)

	time.Sleep(100 * time.Millisecond)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
