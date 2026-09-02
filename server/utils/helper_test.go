package utils

import "testing"

func TestGenerateClientID(t *testing.T) {
	id := GenerateClientID()
	if id == "" {
		t.Error("Expected non-empty client ID")
	}
	if len(id) != 6 {
		t.Errorf("Expected 6-char ID, got %d: %s", len(id), id)
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
