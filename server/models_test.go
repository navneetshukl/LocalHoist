package server

import (
	"encoding/json"
	"testing"
)

func TestRequestPayload(t *testing.T) {
	payload := RequestPayload{
		Method:  "GET",
		URL:     "/api/test",
		Headers: map[string][]string{"Content-Type": {"application/json"}},
		Body:    []byte("test"),
	}
	if payload.Method != "GET" {
		t.Error("Method mismatch")
	}
}

func TestResponsePayload(t *testing.T) {
	payload := ResponsePayload{
		StatusCode: 200,
		Headers:    map[string][]string{},
		Body:       map[string]interface{}{"status": "ok"},
	}
	if payload.StatusCode != 200 {
		t.Error("StatusCode mismatch")
	}
}

func TestWebSocketClientRequest(t *testing.T) {
	req := WebSocketClientRequest{
		Message:   "connected",
		URL:       "https://example.com/tunnel/abc",
		ErrorCode: 0,
	}
	if req.Message != "connected" || req.ErrorCode != 0 {
		t.Error("WebSocketClientRequest mismatch")
	}
}

func TestRequestPayloadJSON(t *testing.T) {
	payload := RequestPayload{
		Method:  "GET",
		URL:     "/api/test",
		Headers: map[string][]string{"Content-Type": {"application/json"}},
		Body:    []byte("test"),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}
	var restored RequestPayload
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}
	if restored.Method != payload.Method {
		t.Errorf("Method mismatch: got %s, want %s", restored.Method, payload.Method)
	}
}

func TestResponsePayloadJSON(t *testing.T) {
	payload := ResponsePayload{
		StatusCode: 200,
		Headers:    map[string][]string{"X-Custom": {"value"}},
		Body:       map[string]interface{}{"status": "ok", "count": 42},
	}
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}
	var restored ResponsePayload
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}
	if restored.StatusCode != payload.StatusCode {
		t.Errorf("StatusCode mismatch: got %d, want %d", restored.StatusCode, payload.StatusCode)
	}
}

func TestWebSocketClientRequestJSON(t *testing.T) {
	req := WebSocketClientRequest{
		Message:   "connected",
		URL:       "https://example.com/tunnel/abc",
		ErrorCode: 0,
	}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}
	var restored WebSocketClientRequest
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("JSON unmarshal failed: %v", err)
	}
	if restored.Message != req.Message {
		t.Errorf("Message mismatch: got %s, want %s", restored.Message, req.Message)
	}
}
