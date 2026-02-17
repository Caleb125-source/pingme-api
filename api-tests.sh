#!/bin/bash

# PingMe API Test Script
# Tests all endpoints with various scenarios

API_URL="http://localhost:8080"
PASSED=0
FAILED=0

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=========================================="
echo "PingMe API Test Suite"
echo "=========================================="
echo ""

# Function to print test results
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ PASS${NC}: $2"
        ((PASSED++))
    else
        echo -e "${RED}✗ FAIL${NC}: $2"
        ((FAILED++))
    fi
}

# Test 1: Greeting endpoint (GET /)
echo "Test 1: Greeting endpoint (GET /)"
response=$(curl -s -w "\n%{http_code}" $API_URL/)
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 200 ] && echo "$body" | grep -q "Welcome to PingMe API"; then
    print_result 0 "Greeting endpoint returns 200 and welcome message"
else
    print_result 1 "Greeting endpoint failed (HTTP $http_code)"
fi
echo ""

# Test 2: Health check endpoint (GET /healthz)
echo "Test 2: Health check endpoint (GET /healthz)"
response=$(curl -s -w "\n%{http_code}" $API_URL/healthz)
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 200 ] && echo "$body" | grep -q "healthy"; then
    print_result 0 "Health check returns 200 and healthy status"
else
    print_result 1 "Health check failed (HTTP $http_code)"
fi
echo ""

# Test 3: Echo endpoint with valid JSON (POST /echo)
echo "Test 3: Echo endpoint with valid JSON (POST /echo)"
response=$(curl -s -w "\n%{http_code}" -X POST $API_URL/echo \
    -H "Content-Type: application/json" \
    -d '{"message": "Hello, World!"}')
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 200 ] && echo "$body" | grep -q "Echo: Hello, World!"; then
    print_result 0 "Echo endpoint processes valid JSON correctly"
else
    print_result 1 "Echo endpoint failed with valid JSON (HTTP $http_code)"
fi
echo ""

# Test 4: Echo endpoint with empty message (validation test)
echo "Test 4: Echo endpoint with empty message (validation)"
response=$(curl -s -w "\n%{http_code}" -X POST $API_URL/echo \
    -H "Content-Type: application/json" \
    -d '{"message": ""}')
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 400 ] && echo "$body" | grep -q "cannot be empty"; then
    print_result 0 "Echo endpoint validates empty message (400)"
else
    print_result 1 "Echo endpoint validation failed (HTTP $http_code)"
fi
echo ""

# Test 5: Echo endpoint with invalid JSON
echo "Test 5: Echo endpoint with invalid JSON"
response=$(curl -s -w "\n%{http_code}" -X POST $API_URL/echo \
    -H "Content-Type: application/json" \
    -d '{"message": invalid}')
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 400 ]; then
    print_result 0 "Echo endpoint rejects invalid JSON (400)"
else
    print_result 1 "Echo endpoint should reject invalid JSON (HTTP $http_code)"
fi
echo ""

# Test 6: Echo endpoint with wrong Content-Type
echo "Test 6: Echo endpoint with wrong Content-Type"
response=$(curl -s -w "\n%{http_code}" -X POST $API_URL/echo \
    -H "Content-Type: text/plain" \
    -d '{"message": "test"}')
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 415 ]; then
    print_result 0 "Echo endpoint validates Content-Type (415)"
else
    print_result 1 "Echo endpoint should validate Content-Type (HTTP $http_code)"
fi
echo ""

# Test 7: Echo endpoint with GET instead of POST (wrong method)
echo "Test 7: Echo endpoint with wrong HTTP method"
response=$(curl -s -w "\n%{http_code}" -X GET $API_URL/echo)
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 405 ]; then
    print_result 0 "Echo endpoint rejects wrong method (405)"
else
    print_result 1 "Echo endpoint should reject wrong method (HTTP $http_code)"
fi
echo ""

# Test 8: Greeting endpoint with POST (wrong method)
echo "Test 8: Greeting endpoint with wrong HTTP method"
response=$(curl -s -w "\n%{http_code}" -X POST $API_URL/)
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 405 ]; then
    print_result 0 "Greeting endpoint rejects wrong method (405)"
else
    print_result 1 "Greeting endpoint should reject wrong method (HTTP $http_code)"
fi
echo ""

# Test 9: Echo endpoint with extra fields (strict validation)
echo "Test 9: Echo endpoint with unexpected fields"
response=$(curl -s -w "\n%{http_code}" -X POST $API_URL/echo \
    -H "Content-Type: application/json" \
    -d '{"message": "test", "extra": "field"}')
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

if [ "$http_code" -eq 400 ]; then
    print_result 0 "Echo endpoint rejects unknown fields (400)"
else
    print_result 1 "Echo endpoint should reject unknown fields (HTTP $http_code)"
fi
echo ""

# Summary
echo "=========================================="
echo "Test Summary"
echo "=========================================="
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"
echo "Total: $((PASSED + FAILED))"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}All tests passed! ✓${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed! ✗${NC}"
    exit 1
fi