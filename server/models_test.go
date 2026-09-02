package server

import "testing"

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
