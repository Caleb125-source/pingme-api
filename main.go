package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Response represents the standard JSON response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// EchoRequest represents the expected JSON input for the echo endpoint
type EchoRequest struct {
	Message string `json:"message"`
}

// EchoData represents the data returned by the echo endpoint
type EchoData struct {
	Original  string    `json:"original"`
	Echoed    string    `json:"echoed"`
	Length    int       `json:"length"`
	Timestamp time.Time `json:"timestamp"`
}

// GreetingData represents the data returned by the greeting endpoint
type GreetingData struct {
	Greeting  string    `json:"greeting"`
	Timestamp time.Time `json:"timestamp"`
}

// HealthData represents the data returned by the health check endpoint
type HealthData struct {
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
}

// respondJSON sends a JSON response with the specified status code
func respondJSON(w http.ResponseWriter, statusCode int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// greetingHandler handles GET requests to the root endpoint
func greetingHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, Response{
			Success: false,
			Error:   "Method not allowed. Use GET.",
		})
		return
	}

	// Create greeting response
	data := GreetingData{
		Greeting:  "Welcome to PingMe API!",
		Timestamp: time.Now().UTC(),
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Message: "Greeting retrieved successfully",
		Data:    data,
	})
}

// healthHandler handles GET requests to the /healthz endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		respondJSON(w, http.StatusMethodNotAllowed, Response{
			Success: false,
			Error:   "Method not allowed. Use GET.",
		})
		return
	}

	// Return health status
	data := HealthData{
		Status: "healthy",
		Time:   time.Now().UTC(),
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Message: "Service is healthy",
		Data:    data,
	})
}

// echoHandler handles POST requests to the /echo endpoint
func echoHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, Response{
			Success: false,
			Error:   "Method not allowed. Use POST.",
		})
		return
	}

	// Verify Content-Type is application/json
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		respondJSON(w, http.StatusUnsupportedMediaType, Response{
			Success: false,
			Error:   "Content-Type must be application/json",
		})
		return
	}

	// Decode JSON request body with strict validation
	var req EchoRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Reject unexpected fields

	if err := decoder.Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   fmt.Sprintf("Invalid JSON: %v", err),
		})
		return
	}

	// Validate that message is not empty
	if req.Message == "" {
		respondJSON(w, http.StatusBadRequest, Response{
			Success: false,
			Error:   "Message field cannot be empty",
		})
		return
	}

	// Create echo response
	data := EchoData{
		Original:  req.Message,
		Echoed:    fmt.Sprintf("Echo: %s", req.Message),
		Length:    len(req.Message),
		Timestamp: time.Now().UTC(),
	}

	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Message: "Echo processed successfully",
		Data:    data,
	})
}

// newServer creates and configures the HTTP server - extracted for testability
func newServer(port string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", greetingHandler)
	mux.HandleFunc("/healthz", healthHandler)
	mux.HandleFunc("/echo", echoHandler)

	return &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

// getPort returns the port from environment variable or default
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

func main() {
	port := getPort()
	server := newServer(port)

	// Start server
	log.Printf("PingMe API starting on port %s...", port)
	log.Printf("Endpoints available:")
	log.Printf("  GET  / - Greeting endpoint")
	log.Printf("  GET  /healthz - Health check endpoint")
	log.Printf("  POST /echo - Echo endpoint")

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
