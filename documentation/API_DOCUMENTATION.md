# PingMe API Documentation

## Base URL
```
http://localhost:8080
```

## Response Format

All endpoints return JSON responses with the following structure:

```json
{
  "success": true|false,
  "message": "Optional message",
  "data": { ... },
  "error": "Error message if success is false"
}
```

---

## Endpoints

### 1. Greeting Endpoint

Get a welcome message with server timestamp.

**Endpoint:** `GET /`

**Request:**
```bash
curl http://localhost:8080/
```

**Response:** `200 OK`
```json
{
  "success": true,
  "message": "Greeting retrieved successfully",
  "data": {
    "greeting": "Welcome to PingMe API!",
    "timestamp": "2024-02-15T10:30:00.000Z"
  }
}
```

**Error Responses:**
- `405 Method Not Allowed` - When using HTTP methods other than GET

---

### 2. Health Check Endpoint

Check if the service is running and healthy. Useful for monitoring, load balancers, and orchestration tools.

**Endpoint:** `GET /healthz`

**Request:**
```bash
curl http://localhost:8080/healthz
```

**Response:** `200 OK`
```json
{
  "success": true,
  "message": "Service is healthy",
  "data": {
    "status": "healthy",
    "time": "2024-02-15T10:30:00.000Z"
  }
}
```

**Error Responses:**
- `405 Method Not Allowed` - When using HTTP methods other than GET

---

### 3. Echo Endpoint

Send a message and receive it back with metadata (length, timestamp, etc.).

**Endpoint:** `POST /echo`

**Headers:**
- `Content-Type: application/json` (Required)

**Request Body:**
```json
{
  "message": "Your message here"
}
```

**Request Example:**
```bash
curl -X POST http://localhost:8080/echo \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello, PingMe!"}'
```

**Response:** `200 OK`
```json
{
  "success": true,
  "message": "Echo processed successfully",
  "data": {
    "original": "Hello, PingMe!",
    "echoed": "Echo: Hello, PingMe!",
    "length": 14,
    "timestamp": "2024-02-15T10:30:00.000Z"
  }
}
```

**Error Responses:**

1. **Missing or Wrong HTTP Method:** `405 Method Not Allowed`
```json
{
  "success": false,
  "error": "Method not allowed. Use POST."
}
```

2. **Wrong Content-Type:** `415 Unsupported Media Type`
```json
{
  "success": false,
  "error": "Content-Type must be application/json"
}
```

3. **Invalid JSON:** `400 Bad Request`
```json
{
  "success": false,
  "error": "Invalid JSON: <error details>"
}
```

4. **Empty Message:** `400 Bad Request`
```json
{
  "success": false,
  "error": "Message field cannot be empty"
}
```

5. **Unknown Fields:** `400 Bad Request`
```json
{
  "success": false,
  "error": "Invalid JSON: json: unknown field \"extraField\""
}
```

---

## HTTP Status Codes

The API uses the following HTTP status codes:

- `200 OK` - Request succeeded
- `400 Bad Request` - Invalid request body or validation error
- `405 Method Not Allowed` - Wrong HTTP method used
- `415 Unsupported Media Type` - Wrong Content-Type header

---

## Examples

### cURL Examples

**Greeting:**
```bash
curl http://localhost:8080/
```

**Health Check:**
```bash
curl http://localhost:8080/healthz
```

**Echo - Success:**
```bash
curl -X POST http://localhost:8080/echo \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello, World!"}'
```

**Echo - Empty Message (Error):**
```bash
curl -X POST http://localhost:8080/echo \
  -H "Content-Type: application/json" \
  -d '{"message": ""}'
```

**Echo - Invalid JSON (Error):**
```bash
curl -X POST http://localhost:8080/echo \
  -H "Content-Type: application/json" \
  -d '{"message": invalid}'
```

---

### JavaScript (Fetch API) Examples

**Greeting:**
```javascript
fetch('http://localhost:8080/')
  .then(response => response.json())
  .then(data => console.log(data));
```

**Echo:**
```javascript
fetch('http://localhost:8080/echo', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    message: 'Hello from JavaScript!'
  })
})
  .then(response => response.json())
  .then(data => console.log(data));
```

---

### Python (Requests) Examples

**Greeting:**
```python
import requests

response = requests.get('http://localhost:8080/')
print(response.json())
```

**Echo:**
```python
import requests

response = requests.post(
    'http://localhost:8080/echo',
    json={'message': 'Hello from Python!'}
)
print(response.json())
```

---

## Error Handling Best Practices

When consuming this API, always:

1. **Check the `success` field** in the response
2. **Handle different HTTP status codes** appropriately
3. **Display the `error` field** to users when `success` is `false`
4. **Set proper Content-Type headers** for POST requests
5. **Validate your JSON** before sending

---

## Rate Limiting

Currently, there is no rate limiting implemented. In a production environment, you would typically implement:

- Request rate limits (e.g., 100 requests per minute)
- IP-based throttling
- API key authentication

---

## CORS

CORS is not currently configured. To enable CORS for cross-origin requests, you would need to add appropriate headers in the response:

```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, OPTIONS
Access-Control-Allow-Headers: Content-Type
```

---

## Monitoring

The `/healthz` endpoint is designed for:

- **Kubernetes/Docker health checks**
- **Load balancer health probes**
- **Uptime monitoring services**
- **CI/CD pipeline smoke tests**

---

## Extending the API

To add new endpoints:

1. Create a handler function following the same pattern
2. Register it with `http.HandleFunc()`
3. Implement proper error handling and validation
4. Update this documentation
5. Add tests to the test suite

Example:
```go
func myNewHandler(w http.ResponseWriter, r *http.Request) {
    // Your implementation here
}

// In main():
http.HandleFunc("/my-endpoint", myNewHandler)
```