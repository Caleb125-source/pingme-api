#!/bin/bash

# PingMe API - Complete Setup and Deployment Script
# This script helps you set up testing and deployment

set -e  # Exit on error

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  PingMe API - Setup & Deployment${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Function to print colored output
print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

# Check prerequisites
check_prerequisites() {
    print_info "Checking prerequisites..."
    
    # Check Go
    if command -v go &> /dev/null; then
        GO_VERSION=$(go version | awk '{print $3}')
        print_success "Go installed: $GO_VERSION"
    else
        print_error "Go is not installed"
        exit 1
    fi
    
    # Check Docker
    if command -v docker &> /dev/null; then
        DOCKER_VERSION=$(docker --version | awk '{print $3}')
        print_success "Docker installed: $DOCKER_VERSION"
    else
        print_warning "Docker not installed (optional for local testing)"
    fi
    
    # Check Git
    if command -v git &> /dev/null; then
        print_success "Git installed"
    else
        print_error "Git is not installed"
        exit 1
    fi
    
    echo ""
}

# Run tests
run_tests() {
    print_info "Running unit tests..."
    
    if go test -v -cover; then
        print_success "All tests passed!"
    else
        print_error "Tests failed"
        exit 1
    fi
    
    echo ""
    
    print_info "Generating coverage report..."
    go test -coverprofile=coverage.out
    go tool cover -html=coverage.out -o coverage.html
    print_success "Coverage report generated: coverage.html"
    
    echo ""
}

# Build Docker image
build_docker() {
    print_info "Building Docker image..."
    
    if docker build -t pingme-api:latest .; then
        print_success "Docker image built successfully"
        
        # Show image size
        SIZE=$(docker images pingme-api:latest --format "{{.Size}}")
        print_info "Image size: $SIZE"
    else
        print_error "Docker build failed"
        exit 1
    fi
    
    echo ""
}

# Test Docker container
test_docker() {
    print_info "Testing Docker container..."
    
    # Run container
    docker run -d -p 8080:8080 --name pingme-test pingme-api:latest
    
    # Wait for container to start
    sleep 3
    
    # Test health endpoint
    if curl -f http://localhost:8080/healthz > /dev/null 2>&1; then
        print_success "Docker container is running and healthy"
    else
        print_error "Docker container health check failed"
        docker logs pingme-test
        docker stop pingme-test
        docker rm pingme-test
        exit 1
    fi
    
    # Cleanup
    docker stop pingme-test > /dev/null 2>&1
    docker rm pingme-test > /dev/null 2>&1
    
    echo ""
}

# Push to Docker Hub
push_docker_hub() {
    print_info "Pushing to Docker Hub..."
    
    read -p "Enter your Docker Hub username: " DOCKER_USER
    
    # Login
    print_info "Logging in to Docker Hub..."
    docker login
    
    # Tag image
    docker tag pingme-api:latest $DOCKER_USER/pingme-api:latest
    docker tag pingme-api:latest $DOCKER_USER/pingme-api:v1.0.0
    
    # Push
    print_info "Pushing images..."
    docker push $DOCKER_USER/pingme-api:latest
    docker push $DOCKER_USER/pingme-api:v1.0.0
    
    print_success "Images pushed to Docker Hub"
    print_info "View at: https://hub.docker.com/r/$DOCKER_USER/pingme-api"
    
    echo ""
}

# Deploy to Fly.io
deploy_flyio() {
    print_info "Deploying to Fly.io..."
    
    # Check if flyctl is installed
    if ! command -v flyctl &> /dev/null; then
        print_warning "flyctl not installed. Install it from: https://fly.io/docs/hands-on/install-flyctl/"
        return
    fi
    
    # Check if already initialized
    if [ ! -f "fly.toml" ]; then
        print_info "Initializing Fly.io app..."
        flyctl launch --no-deploy
    fi
    
    # Deploy
    print_info "Deploying..."
    flyctl deploy
    
    print_success "Deployed to Fly.io!"
    print_info "Opening app in browser..."
    flyctl open
    
    echo ""
}

# Setup GitHub Actions
setup_github_actions() {
    print_info "Setting up GitHub Actions..."
    
    # Create workflows directory if it doesn't exist
    mkdir -p .github/workflows
    
    # Copy workflow files
    if [ -f ".github/workflows/ci.yml" ]; then
        print_success "CI workflow already exists"
    else
        print_warning "CI workflow not found. Please add .github/workflows/ci.yml"
    fi
    
    print_info "Don't forget to add GitHub secrets:"
    echo "  - DOCKER_USERNAME"
    echo "  - DOCKER_PASSWORD"
    echo "  - FLY_API_TOKEN (if using Fly.io)"
    
    echo ""
}

# Main menu
show_menu() {
    echo -e "${BLUE}What would you like to do?${NC}"
    echo "1) Run all tests"
    echo "2) Build Docker image"
    echo "3) Test Docker container"
    echo "4) Push to Docker Hub"
    echo "5) Deploy to Fly.io"
    echo "6) Full setup (tests + Docker + push)"
    echo "7) Exit"
    echo ""
    read -p "Enter choice [1-7]: " choice
    
    case $choice in
        1)
            run_tests
            show_menu
            ;;
        2)
            build_docker
            show_menu
            ;;
        3)
            test_docker
            show_menu
            ;;
        4)
            push_docker_hub
            show_menu
            ;;
        5)
            deploy_flyio
            show_menu
            ;;
        6)
            run_tests
            build_docker
            test_docker
            push_docker_hub
            print_success "Full setup complete!"
            echo ""
            show_menu
            ;;
        7)
            print_info "Goodbye!"
            exit 0
            ;;
        *)
            print_error "Invalid choice"
            show_menu
            ;;
    esac
}

# Main execution
main() {
    check_prerequisites
    show_menu
}

main