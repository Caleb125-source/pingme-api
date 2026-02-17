# PingMe API 🚀

A lightweight, beginner-friendly REST API built with Go (1.21+) that demonstrates core backend development concepts using only the standard `net/http` library.

## 📋 Overview

PingMe API is designed for beginners and students to learn real-world backend patterns. It provides three simple but production-quality endpoints that showcase proper HTTP handling, JSON validation, and error management.

### Features

- ✅ **Pure Go** - Uses only the standard library
- ✅ **Three RESTful Endpoints** - Greeting, Health Check, and Echo
- ✅ **Proper HTTP Handling** - Method validation and status codes
- ✅ **JSON Validation** - Strict input validation with error handling
- ✅ **Docker Ready** - Multi-stage build for minimal image size
- ✅ **Production Patterns** - Timeouts, health checks, and logging
- ✅ **Clean Code** - Readable and well-documented
- ✅ **Automated Testing** - Go unit tests with >80% coverage
- ✅ **CI/CD Pipeline** - GitHub Actions for automated testing and deployment

## 🎯 Endpoints

### 1. Greeting Endpoint
**`GET /`**

Returns a welcome message with timestamp.

**Response:**
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

### 2. Health Check Endpoint
**`GET /healthz`**

Standard health check endpoint for monitoring and orchestration tools.

**Response:**
```json
{
  "success": true,
  "message": "Service is healthy",
  "data": {
    "status": "healthy",
    "time": "2024-02-15T10:30:00Z"
  }
}
```

### 3. Echo Endpoint
**`POST /echo`**

Accepts JSON input and echoes it back with metadata.

**Request:**
```json
{
  "message": "Hello, World!"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Echo processed successfully",
  "data": {
    "original": "Hello, World!",
    "echoed": "Echo: Hello, World!",
    "length": 13,
    "timestamp": "2024-02-15T10:30:00Z"
  }
}
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Docker (optional, for containerized deployment)

### Running Locally

1. **Clone the repository:**
```bash
git clone https://github.com/yourusername/pingme-api.git
cd pingme-api
```

2. **Run directly:**
```bash
go run main.go
```

The API will start on `http://localhost:8080`

3. **Test the endpoints:**

```bash
# Greeting endpoint
curl http://localhost:8080/

# Health check
curl http://localhost:8080/healthz

# Echo endpoint
curl -X POST http://localhost:8080/echo \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello, PingMe!"}'
```

### Running with Docker

```bash
docker build -t pingme-api:latest .
docker run -p 8080:8080 pingme-api:latest
```

### Running with Docker Compose

```bash
docker-compose up -d
docker-compose down  # to stop
```

## 📁 Project Structure

```
pingme-api/
├── .github/
│   └── workflows/
│       ├── ci.yml               # Automated testing on every push
│       └── deploy.yml           # Automated deployment to production
├── documentation/
│   ├── API_DOCUMENTATION.md     # Full API reference
│   ├── CONTRIBUTING.md          # Contribution guidelines
│   ├── DEPLOYMENT.md            # Deployment instructions
│   ├── QUICKSTART.md            # Quick start guide
│   ├── README.md                # This file
│   └── ROADMAP.md               # Project roadmap
├── tests/
│   └── api-tests.sh             # Bash integration tests
├── .dockerignore                # Docker build optimization
├── .gitignore                   # Git ignore rules
├── CHECKLIST.md                 # Deployment checklist
├── docker-compose.yml           # Docker Compose configuration
├── Dockerfile                   # Multi-stage Docker build
├── go.mod                       # Go module definition
├── go.sum                       # Go dependency checksums
├── main.go                      # Main application entry point
├── main_test.go                 # Go unit tests ← NEW
├── Makefile                     # Build and task automation
├── setup.sh                     # Environment setup script
└── TESTING.md                   # Complete testing guide ← NEW
```

## 🧪 Testing

### Unit Tests (Go) — Run Automatically in CI/CD

```bash
# Run all unit tests
go test -v

# Run with coverage
go test -v -cover

# Generate HTML coverage report
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

```bash
chmod +x tests/api-tests.sh
./tests/api-tests.sh
```

### Test Cases Covered

- ✅ GET request to greeting endpoint
- ✅ GET request to health check endpoint
- ✅ POST request with valid JSON to echo endpoint
- ✅ POST request with empty message (validation)
- ✅ POST request with invalid JSON
- ✅ POST request with unknown fields (strict validation)
- ✅ Wrong HTTP method handling (405)
- ✅ Wrong Content-Type handling (415)
- ✅ Empty body handling (400)

See [TESTING.md](../TESTING.md) for the complete testing guide.

## 🤖 CI/CD Pipeline

This project uses GitHub Actions for automated testing and deployment.

### How It Works

```
git push origin main
        ↓
GitHub Actions runs automatically:
  ✅ All unit tests
  ✅ Test coverage check (>80%)
  ✅ Security scan
  ✅ Code quality lint
  ✅ Docker image build
  ✅ Container endpoint tests
        ↓
All pass? → Auto-deploy to production!
```

### View Pipeline Status

Go to your GitHub repo → **Actions** tab to see all workflow runs.

See [DEPLOYMENT.md](DEPLOYMENT.md) for setting up automated deployment.

## 🎓 Learning Objectives

This project teaches:

1. **HTTP Fundamentals**
   - Method handling (GET, POST)
   - Status codes (200, 400, 405, 415)
   - Headers (Content-Type, etc.)

2. **JSON Processing**
   - Encoding and decoding
   - Strict validation
   - Error handling

3. **REST API Design**
   - Consistent response structure
   - Proper endpoint naming
   - Resource-oriented patterns

4. **Production Patterns**
   - Health check endpoints
   - Request timeouts
   - Structured logging
   - Error responses

5. **Testing**
   - Unit tests with Go's `testing` package
   - Integration tests with bash
   - Test coverage reporting

6. **Docker & Deployment**
   - Multi-stage builds
   - Security (non-root user)
   - Health checks
   - Minimal image size

7. **CI/CD**
   - GitHub Actions workflows
   - Automated testing on push
   - Automated deployment

## 🔧 Configuration

The API uses sensible defaults:

- **Port:** 8080
- **Read Timeout:** 10 seconds
- **Write Timeout:** 10 seconds
- **Idle Timeout:** 60 seconds

To modify, edit the `main()` function in `main.go`.

## 📊 Error Handling

The API provides clear error messages:

**Method Not Allowed (405):**
```json
{
  "success": false,
  "error": "Method not allowed. Use POST."
}
```

**Bad Request (400):**
```json
{
  "success": false,
  "error": "Message field cannot be empty"
}
```

**Unsupported Media Type (415):**
```json
{
  "success": false,
  "error": "Content-Type must be application/json"
}
```

## 📚 Documentation

Additional documentation is available in the [`documentation/`](.) directory:

- [API_DOCUMENTATION.md](API_DOCUMENTATION.md) - Full API reference
- [QUICKSTART.md](QUICKSTART.md) - Getting started fast
- [DEPLOYMENT.md](DEPLOYMENT.md) - Deployment guide with CI/CD setup
- [CONTRIBUTING.md](CONTRIBUTING.md) - How to contribute
- [ROADMAP.md](ROADMAP.md) - Planned features

And in the project root:

- [TESTING.md](../TESTING.md) - Complete testing guide
- [CHECKLIST.md](../CHECKLIST.md) - Deployment checklist

## 🎯 Use Cases

Perfect for:

- 📚 Learning Go backend development
- 🎓 Teaching REST API concepts
- 💼 Portfolio projects
- 🔬 Interview demonstrations
- 🚀 Microservice templates
- 📖 Documentation examples

## 🛠️ Extending the API

Want to add more features? Here are some ideas:

1. **Add middleware** for logging, authentication, or CORS
2. **Add a database** (PostgreSQL, MySQL, MongoDB)
3. **Add more endpoints** (CRUD operations)
4. **Add rate limiting** to prevent abuse
5. **Add Swagger documentation** using go-swagger
6. **Add metrics** with Prometheus

## 📝 Best Practices Demonstrated

- ✅ Single Responsibility Principle (separate handlers)
- ✅ Consistent error handling
- ✅ Input validation
- ✅ Structured responses
- ✅ Proper HTTP status codes
- ✅ Security (timeouts, non-root user)
- ✅ Comprehensive documentation
- ✅ Containerization
- ✅ Unit testing with Go
- ✅ CI/CD automation

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) first, then:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Write tests for your changes
4. Commit your changes (`git commit -m 'Add some amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👨‍💻 Author

Your Name - [@yourhandle](https://twitter.com/yourhandle)

Project Link: [https://github.com/yourusername/pingme-api](https://github.com/yourusername/pingme-api)

## 🙏 Acknowledgments

- Built with Go standard library
- Inspired by real-world production APIs
- Designed for learning and teaching

---

**Happy Coding! 🎉**

If you found this project helpful, please give it a ⭐️ on GitHub!