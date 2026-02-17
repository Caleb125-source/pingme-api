package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

	// Check status code
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	// Check Content-Type
	contentType := res.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	// Decode response
	var response Response
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Verify response structure
	if !response.Success {
		t.Error("expected success to be true")
	}

	if response.Message == "" {
		t.Error("expected non-empty message")
	}

	if response.Data == nil {
		t.Fatal("expected data field to be present")
	}

	// Verify data structure (type assertion)
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

	// Verify timestamp is valid RFC3339 format
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

	// Check status code
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	// Decode response
	var response Response
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Verify response structure
	if !response.Success {
		t.Error("expected success to be true")
	}

	if response.Data == nil {
		t.Fatal("expected data field to be present")
	}

	// Verify health status
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

	// Check status code
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	// Decode response
	var response Response
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Verify response structure
	if !response.Success {
		t.Error("expected success to be true")
	}

	if response.Data == nil {
		t.Fatal("expected data field to be present")
	}

	// Verify echo data
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

	// Verify length is correct
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

	// Verify status code
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	// Verify Content-Type
	contentType := res.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	// Verify response can be decoded
	var decoded Response
	if err := json.NewDecoder(res.Body).Decode(&decoded); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if decoded.Success != response.Success {
		t.Errorf("expected success %v, got %v", response.Success, decoded.Success)
	}
}
