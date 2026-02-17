package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// TestGreetingHandler tests the GET / endpoint
func TestGreetingHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	greetingHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	contentType := res.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	var response Response
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !response.Success {
		t.Error("expected success to be true")
	}

	if response.Message == "" {
		t.Error("expected non-empty message")
	}

	if response.Data == nil {
		t.Fatal("expected data field to be present")
	}

	dataMap, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatal("expected data to be a map")
	}

	if dataMap["greeting"] == nil {
		t.Error("expected 'greeting' field in data")
	}

	if dataMap["timestamp"] == nil {
		t.Error("expected 'timestamp' field in data")
	}

	if timestamp, ok := dataMap["timestamp"].(string); ok {
		if _, err := time.Parse(time.RFC3339, timestamp); err != nil {
			t.Errorf("invalid timestamp format: %v", err)
		}
	}
}

// TestGreetingHandlerWrongMethod tests wrong HTTP method on greeting endpoint
func TestGreetingHandlerWrongMethod(t *testing.T) {
	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/", nil)
			w := httptest.NewRecorder()

			greetingHandler(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != http.StatusMethodNotAllowed {
				t.Errorf("expected status 405 for method %s, got %d", method, res.StatusCode)
			}

			var response Response
			if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if response.Success {
				t.Error("expected success to be false for wrong method")
			}

			if response.Error == "" {
				t.Error("expected error message for wrong method")
			}
		})
	}
}

// TestHealthHandler tests the GET /healthz endpoint
func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !response.Success {
		t.Error("expected success to be true")
	}

	if response.Data == nil {
		t.Fatal("expected data field to be present")
	}

	dataMap, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatal("expected data to be a map")
	}

	status, ok := dataMap["status"].(string)
	if !ok || status != "healthy" {
		t.Errorf("expected status 'healthy', got %v", dataMap["status"])
	}
}

// TestHealthHandlerWrongMethod tests wrong HTTP method on health endpoint
func TestHealthHandlerWrongMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/healthz", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", res.StatusCode)
	}
}

// TestEchoHandlerValidJSON tests POST /echo with valid JSON
func TestEchoHandlerValidJSON(t *testing.T) {
	payload := EchoRequest{Message: "Hello, World!"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	echoHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !response.Success {
		t.Error("expected success to be true")
	}

	if response.Data == nil {
		t.Fatal("expected data field to be present")
	}

	dataMap, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatal("expected data to be a map")
	}

	if dataMap["original"] != "Hello, World!" {
		t.Errorf("expected original message 'Hello, World!', got %v", dataMap["original"])
	}

	if dataMap["echoed"] != "Echo: Hello, World!" {
		t.Errorf("expected echoed message 'Echo: Hello, World!', got %v", dataMap["echoed"])
	}

	if length, ok := dataMap["length"].(float64); !ok || int(length) != 13 {
		t.Errorf("expected length 13, got %v", dataMap["length"])
	}
}

// TestEchoHandlerEmptyMessage tests validation for empty message
func TestEchoHandlerEmptyMessage(t *testing.T) {
	payload := EchoRequest{Message: ""}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	echoHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", res.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Success {
		t.Error("expected success to be false for empty message")
	}

	if response.Error == "" {
		t.Error("expected error message for empty message")
	}
}

// TestEchoHandlerInvalidJSON tests handling of malformed JSON
func TestEchoHandlerInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewBufferString("{invalid json}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	echoHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", res.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Success {
		t.Error("expected success to be false for invalid JSON")
	}
}

// TestEchoHandlerUnknownFields tests strict JSON validation
func TestEchoHandlerUnknownFields(t *testing.T) {
	invalidPayload := `{"message": "test", "extra": "field"}`

	req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewBufferString(invalidPayload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	echoHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", res.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Success {
		t.Error("expected success to be false for unknown fields")
	}
}

// TestEchoHandlerWrongContentType tests Content-Type validation
func TestEchoHandlerWrongContentType(t *testing.T) {
	payload := EchoRequest{Message: "test"}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	echoHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnsupportedMediaType {
		t.Errorf("expected status 415, got %d", res.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Success {
		t.Error("expected success to be false for wrong Content-Type")
	}
}

// TestEchoHandlerWrongMethod tests wrong HTTP method
func TestEchoHandlerWrongMethod(t *testing.T) {
	methods := []string{http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/echo", nil)
			w := httptest.NewRecorder()

			echoHandler(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != http.StatusMethodNotAllowed {
				t.Errorf("expected status 405 for method %s, got %d", method, res.StatusCode)
			}
		})
	}
}

// TestEchoHandlerEmptyBody tests handling of empty request body
func TestEchoHandlerEmptyBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	echoHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", res.StatusCode)
	}
}

// TestRespondJSON tests the respondJSON helper function
func TestRespondJSON(t *testing.T) {
	w := httptest.NewRecorder()
	response := Response{
		Success: true,
		Message: "Test message",
		Data:    map[string]string{"key": "value"},
	}

	respondJSON(w, http.StatusOK, response)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	contentType := res.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	var decoded Response
	if err := json.NewDecoder(res.Body).Decode(&decoded); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if decoded.Success != response.Success {
		t.Errorf("expected success %v, got %v", response.Success, decoded.Success)
	}
}

// brokenWriter simulates a ResponseWriter that fails on Write
type brokenWriter struct {
	httptest.ResponseRecorder
}

func (b *brokenWriter) Write(p []byte) (int, error) {
	return 0, fmt.Errorf("simulated write error")
}

// TestRespondJSONEncodingError tests the error path in respondJSON
func TestRespondJSONEncodingError(t *testing.T) {
	w := &brokenWriter{*httptest.NewRecorder()}
	response := Response{
		Success: true,
		Message: "test",
	}
	// Should not panic even when write fails
	respondJSON(w, http.StatusOK, response)
}

// TestNewServer tests that newServer creates a properly configured server
func TestNewServer(t *testing.T) {
	server := newServer("9090")

	if server == nil {
		t.Fatal("expected server to be non-nil")
	}

	if server.Addr != ":9090" {
		t.Errorf("expected addr :9090, got %s", server.Addr)
	}

	if server.ReadTimeout != 10*time.Second {
		t.Errorf("expected ReadTimeout 10s, got %v", server.ReadTimeout)
	}

	if server.WriteTimeout != 10*time.Second {
		t.Errorf("expected WriteTimeout 10s, got %v", server.WriteTimeout)
	}

	if server.IdleTimeout != 60*time.Second {
		t.Errorf("expected IdleTimeout 60s, got %v", server.IdleTimeout)
	}

	if server.Handler == nil {
		t.Error("expected server handler to be set")
	}
}

// TestNewServerRoutes tests that newServer registers all routes correctly
func TestNewServerRoutes(t *testing.T) {
	server := newServer("8080")
	ts := httptest.NewServer(server.Handler)
	defer ts.Close()

	// Test greeting route
	res, err := http.Get(ts.URL + "/")
	if err != nil {
		t.Fatalf("failed to GET /: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200 from /, got %d", res.StatusCode)
	}

	// Test health route
	res, err = http.Get(ts.URL + "/healthz")
	if err != nil {
		t.Fatalf("failed to GET /healthz: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200 from /healthz, got %d", res.StatusCode)
	}

	// Test echo route
	payload := bytes.NewBufferString(`{"message": "test"}`)
	res, err = http.Post(ts.URL+"/echo", "application/json", payload)
	if err != nil {
		t.Fatalf("failed to POST /echo: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200 from /echo, got %d", res.StatusCode)
	}
}

// TestGetPort tests the getPort function
func TestGetPort(t *testing.T) {
	// Test default port
	os.Unsetenv("PORT")
	port := getPort()
	if port != "8080" {
		t.Errorf("expected default port 8080, got %s", port)
	}

	// Test custom port from environment
	os.Setenv("PORT", "3000")
	defer os.Unsetenv("PORT")
	port = getPort()
	if port != "3000" {
		t.Errorf("expected port 3000, got %s", port)
	}
}

// BenchmarkEchoHandler benchmarks the echo endpoint performance
func BenchmarkEchoHandler(b *testing.B) {
	payload := `{"message": "benchmark test"}`
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/echo",
			bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		echoHandler(w, req)
	}
}
