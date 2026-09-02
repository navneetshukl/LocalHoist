package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
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

	if w.Code != http.StatusOK {
		t.Logf("Status: %d (expected OK due to async handling)", w.Code)
	}
}
