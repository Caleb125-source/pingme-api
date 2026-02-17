# PingMe API - Quick Start Guide

Get up and running with PingMe API in less than 5 minutes!

## 🚀 Option 1: Run Locally (Recommended for Development)

### Prerequisites
- Go 1.21 or higher installed
- Terminal/Command Prompt

### Steps

1. **Clone or download the project:**
```bash
git clone https://github.com/yourusername/pingme-api.git
cd pingme-api
```

2. **Run the API:**
```bash
go run main.go
```

You should see:
```
2024/02/15 10:30:00 PingMe API starting on port 8080...
2024/02/15 10:30:00 Endpoints available:
2024/02/15 10:30:00   GET  / - Greeting endpoint
2024/02/15 10:30:00   GET  /healthz - Health check endpoint
2024/02/15 10:30:00   POST /echo - Echo endpoint
```

3. **Test it works:**

Open a new terminal and run:
```bash
curl http://localhost:8080/
```

You should get:
```json
{
  "success": true,
  "message": "Greeting retrieved successfully",
  "data": {
    "greeting": "Welcome to PingMe API!",
    "timestamp": "2024-02-15T10:30:00Z"
  }
}
```

**🎉 Success! Your API is running!**

---

## 🐳 Option 2: Run with Docker

### Prerequisites
- Docker installed

### Steps

1. **Clone or download the project:**
```bash
git clone https://github.com/yourusername/pingme-api.git
cd pingme-api
```

2. **Build the Docker image:**
```bash
docker build -t pingme-api:latest .
```

3. **Run the container:**
```bash
docker run -p 8080:8080 pingme-api:latest
```

4. **Test it works:**
```bash
curl http://localhost:8080/healthz
```

**🎉 Your containerized API is running!**

---

## 🔧 Option 3: Run with Docker Compose (Easiest)

### Prerequisites
- Docker and Docker Compose installed

### Steps

1. **Clone or download the project:**
```bash
git clone https://github.com/yourusername/pingme-api.git
cd pingme-api
```

2. **Start everything:**
```bash
docker-compose up -d
```

3. **Test it works:**
```bash
curl http://localhost:8080/
```

4. **Stop when done:**
```bash
docker-compose down
```

**🎉 The easiest way to run the API!**

---

## 🧪 Testing All Endpoints

Once the API is running, test all endpoints:

### 1. Greeting Endpoint
```bash
curl http://localhost:8080/
```

### 2. Health Check
```bash
curl http://localhost:8080/healthz
```

### 3. Echo Endpoint
```bash
curl -X POST http://localhost:8080/echo \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello, PingMe!"}'
```

---

## ✅ Run the Full Test Suite

### Unit Tests (Go) — Automated

These are the primary tests. Run them before every commit:

```bash
# Run all unit tests
go test -v

# Run with coverage report
go test -v -cover

# Generate HTML coverage report (opens in browser)
go test -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

Expected output:
```
=== RUN   TestGreetingHandler
--- PASS: TestGreetingHandler (0.00s)
=== RUN   TestHealthHandler
--- PASS: TestHealthHandler (0.00s)
=== RUN   TestEchoHandlerValidJSON
--- PASS: TestEchoHandlerValidJSON (0.00s)
...
PASS
coverage: 85.2% of statements
```

### Integration Tests (Bash) — Manual

Run these to test your live running server:

```bash
# Make the script executable (first time only)
chmod +x tests/api-tests.sh

# Start the server first
go run main.go &

# Run integration tests
./tests/api-tests.sh
```

> **Note:** Unit tests (`go test -v`) run without needing the server running.
> Integration tests (`api-tests.sh`) require the server to be running.

---

## 📱 Using Postman or Thunder Client

1. Set base URL to `http://localhost:8080`
2. Test these requests:

| Method | URL | Body |
|--------|-----|------|
| GET | `/` | None |
| GET | `/healthz` | None |
| POST | `/echo` | `{"message": "Hello!"}` |

---

## 🛠️ Using the Makefile (Optional)

If you have `make` installed:

```bash
# See all available commands
make help

# Run locally
make run

# Build binary
make build

# Run tests
make test

# Docker operations
make docker-build
make docker-run
```

---

## 🎯 What's Next?

1. **Explore the code** - Check out `main.go` to see how it works
2. **Read the tests** - Check `main_test.go` to understand what's tested
3. **Read the docs** - See `API_DOCUMENTATION.md` for detailed API reference
4. **Set up CI/CD** - Push to GitHub and watch automated testing run
5. **Deploy it** - Follow `DEPLOYMENT.md` to go live

---

## ❓ Troubleshooting

### Port 8080 already in use?

```bash
# Find what's using port 8080
lsof -i :8080

# Kill it
kill -9 <PID>

# Or change the port in main.go
port := "3000"
```

### Tests failing?

```bash
# Make sure you're in the project root
ls main.go main_test.go

# Clear test cache and retry
go clean -testcache
go test -v
```

### Permission denied on test script?

```bash
chmod +x tests/api-tests.sh
```

### Docker build fails?

```bash
# Make sure Docker is running
docker --version
docker ps
```

### Can't reach the API?

1. Check if it's running: `curl http://localhost:8080/healthz`
2. Check the logs: `docker logs pingme-api` (if using Docker)
3. Make sure no firewall is blocking port 8080

---

## 💡 Tips

- Run `go test -v` before every `git push`
- The API logs all requests to console
- Use the health check endpoint for monitoring
- All responses are in JSON format
- Errors include helpful messages

---

## 📚 Learn More

- **Full Documentation:** See `README.md`
- **API Reference:** See `documentation/API_DOCUMENTATION.md`
- **Testing Guide:** See `TESTING.md`
- **Deployment Guide:** See `documentation/DEPLOYMENT.md`

---

**Happy coding! 🚀**

If you found this helpful, please star the repo on GitHub!