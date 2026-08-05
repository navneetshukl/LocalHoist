package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DecodeRequest will read the response and convert it to http protocol payload
func DecodeRequest(data []byte) (*RequestPayload, error) {
	var payload RequestPayload

	err := json.Unmarshal(data, &payload)
	if err != nil {
		return nil, err
	}
	return &payload, nil
}

// Execute performs the HTTP request specified by the RequestPayload
func Execute(r *RequestPayload) (*http.Response, error) {
	var bodyReader io.Reader
	if len(r.Body) > 0 {
		bodyReader = bytes.NewReader(r.Body)
	}

	req, err := http.NewRequest(r.Method, r.URL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, values := range r.Headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return client.Do(req)
}
