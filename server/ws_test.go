package server

import "testing"

func TestNewWSManager(t *testing.T) {
	manager := NewWSManager()
	if manager == nil {
		t.Fatal("Expected non-nil manager")
	}
	if manager.wsConn == nil {
		t.Error("Expected non-nil wsConn map")
	}
}

func TestForwardRequestNoConnection(t *testing.T) {
	manager := NewWSManager()
	err := manager.ForwardRequest("nonexistent", nil)
	if err == nil {
		t.Error("Expected error for non-existent client ID")
	}
}
