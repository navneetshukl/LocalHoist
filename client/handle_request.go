package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
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
	host := fmt.Sprintf("http://localhost:%d", r.Port)
	url := fmt.Sprintf("%s%s", host, r.URL)
	url, err := CleanTunnelURL(url)
	if err != nil {
		return nil, fmt.Errorf("error in parsing url: %w", err)
	}

	log.Println("URL is ", url)

	req, err := http.NewRequest(r.Method, url, bodyReader)
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

func CleanTunnelURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	cleanPath := strings.ReplaceAll(u.Path, "//", "/")

	re := regexp.MustCompile(`^/tunnel/[^/]+`)
	newPath := re.ReplaceAllString(cleanPath, "")

	if newPath == "" {
		newPath = "/"
	}
	u.Path = newPath
	return u.String(), nil
}
