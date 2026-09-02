package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"route": "/", "method": "GET"})
	})
	r.GET("/a", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"route": "/a", "method": "GET"})
	})
	r.GET("/a/b", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"route": "/a/b", "method": "GET"})
	})
	r.POST("/a", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"route": "/a", "method": "POST"})
	})
	return r
}

func TestRootRoute(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestRouteA(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/a", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["route"] != "/a" {
		t.Errorf("Expected route '/a', got '%v'", resp["route"])
	}
}

func TestNestedRoute(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/a/b", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["route"] != "/a/b" {
		t.Errorf("Expected route '/a/b', got '%v'", resp["route"])
	}
}
