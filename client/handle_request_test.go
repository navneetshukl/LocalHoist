package main

import (
	"encoding/json"
	"testing"
)

func TestDecodeRequest(t *testing.T) {
	payload, err := DecodeRequest([]byte(`{"method":"POST","url":"/api","port":8080}`))
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	if payload.Method != "POST" || payload.Port != 8080 {
		t.Errorf("Got: %+v", payload)
	}
}

func TestCleanTunnelURL(t *testing.T) {
	result, err := CleanTunnelURL("http://localhost:3000/tunnel/abc/api/test")
	if err != nil || result == "" {
		t.Errorf("Error: %v, Result: %s", err, result)
	}
}

func TestPayloadRoundTrip(t *testing.T) {
	orig := RequestPayload{Method: "GET", URL: "/test", Port: 8080}
	data, _ := json.Marshal(orig)
	var restored RequestPayload
	json.Unmarshal(data, &restored)
	if restored.Method != orig.Method {
		t.Error("Round trip failed")
	}
}
