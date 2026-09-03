package utils

import (
	"testing"
)

func TestGenerateClientID(t *testing.T) {
	id := GenerateClientID()
	if id == "" {
		t.Error("Expected non-empty client ID")
	}
	if len(id) != 6 {
		t.Errorf("Expected 6-char ID, got %d: %s", len(id), id)
	}
}

func TestGenerateClientIDUniqueness(t *testing.T) {
	// Generate multiple IDs and verify they're all non-empty
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := GenerateClientID()
		if id == "" {
			t.Error("Expected non-empty client ID")
		}
		ids[id] = true
	}
	if len(ids) < 50 {
		t.Errorf("Expected at least 50 unique IDs out of 100, got %d", len(ids))
	}
}

func TestGenerateClientIDFormat(t *testing.T) {
	id := GenerateClientID()
	// Hex encoding of 3 bytes = 6 hex chars
	for _, c := range id {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("Expected hex character, got: %c", c)
		}
	}
}

func TestGetClientIdFromRouteValid(t *testing.T) {
	err, id := GetClientIdFromRoute("/tunnel/abc123/api")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if id != "abc123" {
		t.Errorf("Expected id 'abc123', got '%s'", id)
	}
}

func TestGetClientIdFromRouteInvalid(t *testing.T) {
	err, _ := GetClientIdFromRoute("/invalid")
	if err == nil {
		t.Error("Expected error for invalid route")
	}
}

func TestGetClientIdFromRouteEmpty(t *testing.T) {
	err, _ := GetClientIdFromRoute("")
	if err == nil {
		t.Error("Expected error for empty route")
	}
}

func TestGetClientIdFromRouteOnlyTunnel(t *testing.T) {
	err, id := GetClientIdFromRoute("/tunnel/")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if id != "" {
		t.Errorf("Expected empty client ID, got %s", id)
	}
}

func TestGetClientIdFromRouteWithQuery(t *testing.T) {
	err, id := GetClientIdFromRoute("/tunnel/xyz789/api?key=value")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if id != "xyz789" {
		t.Errorf("Expected id 'xyz789', got '%s'", id)
	}
}
