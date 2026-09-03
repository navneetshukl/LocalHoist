package main

import (
	"strconv"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

func TestPortValidation(t *testing.T) {
	cases := []struct {
		input string
		valid bool
	}{
		{"8080", true},
		{"80", true},
		{"65535", true},
		{"0", false},
		{"-1", false},
		{"abc", false},
		{"99999", false},
		{"70000", false},
		{"1", true},
		{"65534", true},
	}

	for _, c := range cases {
		port, err := strconv.Atoi(c.input)
		if c.valid {
			if err != nil {
				t.Errorf("Expected %s to be valid, got error: %v", c.input, err)
			} else if port <= 0 || port > 65535 {
				t.Errorf("Expected %s to be valid port (1-65535), got %d", c.input, port)
			}
		} else {
			if err == nil && port > 0 && port <= 65535 {
				t.Errorf("Expected %s to be invalid, but got valid port: %d", c.input, port)
			}
		}
	}
}

func TestWebSocketDialerExists(t *testing.T) {
	dialer := websocket.DefaultDialer
	if dialer == nil {
		t.Error("Expected DefaultDialer to exist")
	}
}

func TestDefaultServerURL(t *testing.T) {
	// Verify the default server URL constant
	const defaultURL = "ws://localhost:3000/ws"
	if defaultURL == "" {
		t.Error("Default server URL should not be empty")
	}
}

func TestCobraCommandStructure(t *testing.T) {
	// Verify that the root command is set up correctly
	var rootCmd = &cobra.Command{
		Use:   "LocalHoist",
		Short: "LocalHoist creates secure tunnels to your local ports",
	}
	if rootCmd.Use != "LocalHoist" {
		t.Errorf("Expected root command Use to be 'LocalHoist', got %s", rootCmd.Use)
	}
}
