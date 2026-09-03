package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRootRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRouter()
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["route"] != "/" {
		t.Errorf("Expected route '/', got '%v'", resp["route"])
	}
}

func TestRouteA_GET(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRouter()
	req, _ := http.NewRequest("GET", "/a", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["route"] != "/a" {
		t.Errorf("Expected route '/a', got '%v'", resp["route"])
	}
}

func TestNestedRoute_GET(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRouter()
	req, _ := http.NewRequest("GET", "/a/b", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["route"] != "/a/b" {
		t.Errorf("Expected route '/a/b', got '%v'", resp["route"])
	}
}

func TestRouteA_POST_WithJSONBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRouter()
	body := map[string]interface{}{"name": "Navneet", "age": 28}
	bodyJSON, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/a", bytes.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	receivedData := resp["received_data"].(map[string]interface{})
	if receivedData["name"] != "Navneet" {
		t.Errorf("Expected name 'Navneet', got '%v'", receivedData["name"])
	}
}

func TestRouteA_POST_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRouter()
	req, _ := http.NewRequest("POST", "/a", strings.NewReader("not valid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestRouteA_POST_EmptyBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRouter()
	req, _ := http.NewRequest("POST", "/a", bytes.NewReader([]byte("")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for empty body, got %d", w.Code)
	}
}

func TestNestedRoute_POST(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRouter()
	req, _ := http.NewRequest("POST", "/a/b", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestNotFoundRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRouter()
	req, _ := http.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestWrongMethodOnRouteA(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRouter()
	req, _ := http.NewRequest("DELETE", "/a", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestSetupRouter_NotNil(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := SetupRouter()
	if router == nil {
		t.Fatal("Expected non-nil router")
	}
}
