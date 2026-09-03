package main

import (
	"encoding/json"
	"testing"
)

func TestRequestPayloadJSON(t *testing.T) {
	orig := RequestPayload{
		Method:  "GET",
		URL:     "/test",
		Port:    8080,
	}
	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}
	if string(data) != `{"method":"GET","url":"/test","headers":null,"body":null,"port":8080}` {
		t.Errorf("Unexpected JSON: %s", string(data))
	}
}

func TestResponsePayloadJSON(t *testing.T) {
	orig := ResponsePayload{
		StatusCode: 201,
		Headers: map[string][]string{
			"X-Request-ID": {"req-123"},
		},
		Body: map[string]interface{}{
			"message": "success",
			"count": 42,
		},
	}
	data, err := json.Marshal(orig)
	if err != nil {
		t.Fatalf("JSON marshal failed: %v", err)
	}
	if string(data) != `{"status_code":201,"headers":{"X-Request-ID":["req-123"]},"body":{"count":42,"message":"success"}}` {
		t.Errorf("Unexpected JSON: %s", string(data))
	}
}