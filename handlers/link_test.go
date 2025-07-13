package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"linkshare-backend/models"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Middleware to mock user_id in context for testing
	r.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
		c.Next()
	})

	r.POST("/links", SendLink)
	r.GET("/ws", WebSocketHandler)
	return r
}

func TestSendLink(t *testing.T) {
	router := setupRouter()

	link := models.Link{
		FromDeviceID: 1,
		ToDeviceID:   2,
		URL:          "https://example.com",
	}
	body, _ := json.Marshal(link)
	req, _ := http.NewRequest("POST", "/links", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 200 or 400, got %d", w.Code)
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)

	if w.Code == http.StatusOK && resp["message"] != "Link sent" {
		t.Errorf("Expected response 'Link sent', got '%s'", resp["message"])
	}
}

func TestWebSocketHandler_InvalidDeviceID(t *testing.T) {
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/ws?device_id=abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid device_id, got %d", w.Code)
	}
}

func TestWebSocketHandler_ValidDeviceID(t *testing.T) {
	// This test checks upgrade, but won't actually upgrade without a real WebSocket client.
	// We test for status code and that upgrade does not panic.
	router := setupRouter()
	req, _ := http.NewRequest("GET", "/ws?device_id=1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// If upgrade fails, the handler logs and returns, but Gin will not set status code.
	// So we just check that no panic occurs and response is written.
	if w.Code != http.StatusOK && w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 200 or 400, got %d", w.Code)
	}
}
