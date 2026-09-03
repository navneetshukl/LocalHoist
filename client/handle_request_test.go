package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
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

// TestExecuteWithTestServer tests Execute using a test HTTP server
func TestExecuteWithTestServer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer ts.Close()

	// Extract port from test server URL
	_, port, _ := net.SplitHostPort(ts.Listener.Addr().String())
	// Convert port string to int
	portInt := 0
	if port != "" {
		// Use fmt.Sscanf to parse port
		fmt.Sscanf(port, "%d", &portInt)
	}

	payload := &RequestPayload{
		Method: "GET",
		URL:    "/",
		Port:   portInt,
	}

	resp, err := Execute(payload)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if resp == nil {
		t.Fatal("Expected non-nil response")
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if string(body) != `{"status":"ok"}` {
		t.Errorf("Expected ok body, got %s", string(body))
	}
}

// TestCleanTunnelURLEdgeCases tests CleanTunnelURL with various edge cases
func TestCleanTunnelURLEdgeCases(t *testing.T) {
	tests := []struct {
		input    string
		want     string
		wantErr bool
	}{
		{"http://localhost:3000/tunnel/abc/api/test", "http://localhost:3000/api/test", false},
		{"http://localhost:3000/tunnel/", "http://localhost:3000/tunnel/", false},
		{"http://localhost:3000/", "http://localhost:3000/", false},
		{"http://localhost:3000//double", "http://localhost:3000/double", false},
		{"", "/", false},
		{"no/tunnel/here", "no/tunnel/here", false},
		{"http://localhost:3000/tunnel//api/test", "http://localhost:3000/test", false},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result, err := CleanTunnelURL(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tc.want {
					t.Errorf("Expected %q, got %q", tc.want, result)
				}
			}
		})
	}
}
