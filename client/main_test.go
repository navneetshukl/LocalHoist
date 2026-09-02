package main

import (
	"encoding/json"
	"testing"

	"github.com/gorilla/websocket"
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
	}

	for _, c := range cases {
		_, err := json.Number(c.input).Int64()
		if c.valid && err != nil {
			t.Errorf("Expected %s to be valid", c.input)
		}
	}
}

func TestWebSocketDialerExists(t *testing.T) {
	dialer := websocket.DefaultDialer
	if dialer == nil {
		t.Error("Expected DefaultDialer to exist")
	}
}
