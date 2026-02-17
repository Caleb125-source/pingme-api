# Testing Guide

This document explains how to test the PingMe API using different methods.

## Table of Contents
- [Unit Tests](#unit-tests)
- [Integration Tests](#integration-tests)
- [Manual Testing](#manual-testing)
- [CI/CD Testing](#cicd-testing)
- [Coverage Reports](#coverage-reports)

## Unit Tests

PingMe API includes comprehensive Go unit tests that test all endpoints and edge cases.

### Running Unit Tests

```bash
# Run all tests
go test -v

# Run tests with coverage
go test -v -cover

# Run tests with race detection
go test -v -race

# Run specific test
go test -v -run TestGreetingHandler
```

### Test Coverage

Generate and view test coverage:

```bash
# Generate coverage profile
go test -coverprofile=coverage.out

# View coverage in terminal
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Open in browser (macOS)
open coverage.html

# Open in browser (Linux)
xdg-open coverage.html

# Open in browser (Windows)
start coverage.html
```

### What's Tested

The unit tests cover:

1. **Greeting Endpoint (`/`)**
   - ✅ Successful GET request
   - ✅ Wrong HTTP methods (POST, PUT, DELETE, PATCH)
   - ✅ JSON response structure
   - ✅ Timestamp validation
   - ✅ Content-Type header

2. **Health Check Endpoint (`/healthz`)**
   - ✅ Successful GET request
   - ✅ Wrong HTTP methods
   - ✅ Health status validation
   - ✅ Response structure

3. **Echo Endpoint (`/echo`)**
   - ✅ Valid JSON input
   - ✅ Empty message validation
   - ✅ Invalid JSON handling
   - ✅ Unknown fields rejection
   - ✅ Wrong Content-Type handling
   - ✅ Wrong HTTP methods
   - ✅ Empty body handling
   - ✅ Message length calculation

4. **Helper Functions**
   - ✅ `respondJSON` function

### Expected Coverage

Target: **>80%** code coverage

Current coverage breakdown:
- Handlers: ~95%
- Helper functions: 100%
- Main function: Not tested (standard practice)

## Integration Tests

The `tests/api-tests.sh` script provides end-to-end integration testing.

### Running Integration Tests

```bash
# Make script executable
chmod +x tests/api-tests.sh

# Start the server first
go run main.go &

# Run tests
./tests/api-tests.sh

# Kill the server
pkill -f "go run main.go"
```

### What's Tested

Integration tests verify:

1. Server is running and accessible
2. All endpoints respond correctly
3. Error handling works end-to-end
4. JSON parsing and validation
5. HTTP status codes
6. Response format consistency

## Manual Testing

### Using curl

**Test Greeting Endpoint:**
```bash
curl -X GET http://localhost:8080/ | jq
```

Expected response:
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

**Test Health Check:**
```bash
curl -X GET http://localhost:8080/healthz | jq
```

**Test Echo Endpoint:**
```bash
curl -X POST http://localhost:8080/echo \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello, World!"}' | jq
```

**Test Error Handling:**
```bash
# Wrong method
curl -X POST http://localhost:8080/ | jq

# Empty message
curl -X POST http://localhost:8080/echo \
  -H "Content-Type: application/json" \
  -d '{"message": ""}' | jq

# Invalid JSON
curl -X POST http://localhost:8080/echo \
  -H "Content-Type: application/json" \
  -d '{invalid}' | jq

# Wrong Content-Type
curl -X POST http://localhost:8080/echo \
  -H "Content-Type: text/plain" \
  -d '{"message": "test"}' | jq

# Unknown fields
curl -X POST http://localhost:8080/echo \
  -H "Content-Type: application/json" \
  -d '{"message": "test", "extra": "field"}' | jq
```

### Using Postman

1. Import the collection: `tests/postman_collection.json`
2. Run the entire collection
3. Verify all tests pass

### Using Web Browser

Only GET endpoints work in the browser:

- http://localhost:8080/
- http://localhost:8080/healthz

For POST endpoints, use curl, Postman, or any API client.

## Docker Testing

### Test in Docker Container

```bash
# Build image
docker build -t pingme-api:latest .

# Run container
docker run -d -p 8080:8080 --name pingme-test pingme-api:latest

# Test endpoints
curl http://localhost:8080/healthz

# Run integration tests
./tests/api-tests.sh

# Check logs
docker logs pingme-test

# Stop and remove
docker stop pingme-test
docker rm pingme-test
```

### Test Docker Compose

```bash
# Start services
docker-compose up -d

# Run tests
./tests/api-tests.sh

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## CI/CD Testing

Tests run automatically on every push via GitHub Actions.

### Workflows

1. **CI Pipeline** (`.github/workflows/ci.yml`)
   - Runs unit tests
   - Checks test coverage (must be >80%)
   - Builds Docker image
   - Tests Docker container
   - Runs security scans
   - Runs linters

2. **Deploy Pipeline** (`.github/workflows/deploy.yml`)
   - Publishes to Docker Hub
   - Publishes to GitHub Container Registry
   - Deploys to Fly.io (if configured)

### Viewing CI Results

1. Go to your repository on GitHub
2. Click "Actions" tab
3. View workflow runs
4. Click on a specific run to see details

### Coverage Badge

Add to your README:

```markdown
![Tests](https://github.com/yourusername/pingme-api/workflows/CI%2FCD%20Pipeline/badge.svg)
![Coverage](https://codecov.io/gh/yourusername/pingme-api/branch/main/graph/badge.svg)
```

## Continuous Testing

### Watch Mode

For development, use a file watcher:

```bash
# Install watchexec
brew install watchexec  # macOS
# or
cargo install watchexec  # Cross-platform

# Auto-run tests on file changes
watchexec -e go -r "go test -v"
```

### Pre-commit Hook

Add to `.git/hooks/pre-commit`:

```bash
#!/bin/bash
echo "Running tests..."
go test -v || exit 1
echo "All tests passed!"
```

Make executable:
```bash
chmod +x .git/hooks/pre-commit
```

## Benchmarking

### Run Benchmarks

```bash
# Run all benchmarks
go test -bench=. -benchmem

# Run specific benchmark
go test -bench=BenchmarkEchoHandler -benchmem
```

### Load Testing

Use Apache Bench:

```bash
# Install ab
sudo apt-get install apache2-utils  # Ubuntu
brew install httpie  # macOS

# Test greeting endpoint
ab -n 1000 -c 10 http://localhost:8080/

# Test echo endpoint
echo '{"message":"test"}' > /tmp/payload.json
ab -n 1000 -c 10 -p /tmp/payload.json \
  -T application/json http://localhost:8080/echo
```

## Troubleshooting

### Tests Fail Locally

1. Ensure no other service is using port 8080
2. Check Go version: `go version` (should be 1.21+)
3. Update dependencies: `go mod download`
4. Clear cache: `go clean -testcache`

### Docker Tests Fail

1. Rebuild image: `docker build --no-cache -t pingme-api:latest .`
2. Check logs: `docker logs pingme-test`
3. Verify port mapping: `docker ps`
4. Test health endpoint directly: `curl http://localhost:8080/healthz`

### CI/CD Tests Fail

1. Check GitHub Actions logs
2. Verify secrets are set (DOCKER_USERNAME, DOCKER_PASSWORD, etc.)
3. Ensure `go.mod` and `go.sum` are committed
4. Check branch protection rules

## Best Practices

1. **Run tests before committing**
   ```bash
   go test -v && git commit
   ```

2. **Maintain >80% coverage**
   ```bash
   go test -cover
   ```

3. **Test edge cases**
   - Empty inputs
   - Invalid JSON
   - Wrong methods
   - Missing headers

4. **Use table-driven tests** for multiple scenarios

5. **Mock external dependencies** (when you add them)

6. **Test error paths** as thoroughly as success paths

## Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Go Coverage Tool](https://blog.golang.org/cover)
- [TestMain Pattern](https://golang.org/pkg/testing/#hdr-Main)
- [httptest Package](https://golang.org/pkg/net/http/httptest/)

---

**Questions?** Open an issue on GitHub!