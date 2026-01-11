#!/bin/bash
# Development Environment Setup Script for Lumbrera
# This script sets up the complete development environment

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

# Print colored message
print_step() {
    echo -e "\n${BLUE}==>${NC} ${1}"
}

print_success() {
    echo -e "${GREEN}✓${NC} ${1}"
}

print_error() {
    echo -e "${RED}✗${NC} ${1}"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} ${1}"
}

# Check if command exists
command_exists() {
    command -v "$1" &> /dev/null
}

# Print banner
echo -e "${GREEN}"
cat << "EOF"
╔═══════════════════════════════════════════╗
║                                           ║
║   Lumbrera Development Setup              ║
║   Gamified Bible Study App                ║
║                                           ║
╚═══════════════════════════════════════════╝
EOF
echo -e "${NC}"

print_step "Checking system requirements..."

# Check for required tools
REQUIRED_TOOLS=("git" "docker" "docker-compose" "go" "make")
OPTIONAL_TOOLS=("aws" "flutter" "jq")

for tool in "${REQUIRED_TOOLS[@]}"; do
    if command_exists "$tool"; then
        print_success "$tool is installed"
    else
        print_error "$tool is NOT installed (required)"
        echo "Please install $tool and try again"
        exit 1
    fi
done

for tool in "${OPTIONAL_TOOLS[@]}"; do
    if command_exists "$tool"; then
        print_success "$tool is installed"
    else
        print_warning "$tool is NOT installed (optional)"
    fi
done

# Check Go version
print_step "Verifying Go version..."
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_GO_VERSION="1.22"

if [ "$(printf '%s\n' "$REQUIRED_GO_VERSION" "$GO_VERSION" | sort -V | head -n1)" = "$REQUIRED_GO_VERSION" ]; then
    print_success "Go version $GO_VERSION meets requirement (>= $REQUIRED_GO_VERSION)"
else
    print_error "Go version $GO_VERSION is too old. Please upgrade to >= $REQUIRED_GO_VERSION"
    exit 1
fi

# Check Docker
print_step "Verifying Docker..."
if docker info &> /dev/null; then
    print_success "Docker is running"
else
    print_error "Docker is not running. Please start Docker and try again"
    exit 1
fi

# Create necessary directories
print_step "Creating project directories..."
mkdir -p "$ROOT_DIR/bin"
mkdir -p "$ROOT_DIR/logs"
mkdir -p "$ROOT_DIR/mobile"
mkdir -p "$ROOT_DIR/tests/integration"
mkdir -p "$ROOT_DIR/tests/e2e/features"
mkdir -p "$ROOT_DIR/tests/e2e/steps"
mkdir -p "$ROOT_DIR/docs"
print_success "Directories created"

# Create .env file if it doesn't exist
print_step "Setting up environment variables..."
if [ ! -f "$ROOT_DIR/.env" ]; then
    cat > "$ROOT_DIR/.env" << 'ENVEOF'
# Lumbrera Environment Configuration
# Copy this file to .env and update with your values

# AWS Configuration
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=DUMMYIDEXAMPLE
AWS_SECRET_ACCESS_KEY=DUMMYEXAMPLEKEY

# DynamoDB
DYNAMODB_ENDPOINT=http://localhost:8000
TABLE_NAME=lumbrera-local

# Auth0 Configuration
AUTH0_DOMAIN=your-tenant.auth0.com
AUTH0_CLIENT_ID=your-client-id
AUTH0_CLIENT_SECRET=your-client-secret
AUTH0_AUDIENCE=https://api.lumbrera.app

# Application
STAGE=local
API_BASE_URL=http://localhost:3000

# Mobile App (Flutter)
FLUTTER_APP_NAME=Lumbrera
FLUTTER_BUNDLE_ID=com.lumbrera.app

# Feature Flags
ENABLE_DEBUG_LOGGING=true
ENABLE_ANALYTICS=false
ENVEOF
    print_success ".env file created (please update with your Auth0 credentials)"
else
    print_success ".env file already exists"
fi

# Install Go dependencies
print_step "Installing Go dependencies..."
cd "$ROOT_DIR"
go mod download
go mod tidy
print_success "Go dependencies installed"

# Install Go development tools
print_step "Installing Go development tools..."
GOTOOLS=(
    "github.com/cucumber/godog/cmd/godog@latest"
    "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
)

for tool in "${GOTOOLS[@]}"; do
    tool_name=$(echo "$tool" | cut -d'@' -f1 | awk -F'/' '{print $NF}')
    if command_exists "$tool_name"; then
        print_warning "$tool_name already installed, skipping"
    else
        go install "$tool"
        print_success "Installed $tool_name"
    fi
done

# Build backend
print_step "Building backend..."
cd "$ROOT_DIR"
make clean
make build
print_success "Backend built successfully"

# Stop any existing containers
print_step "Cleaning up existing Docker containers..."
docker-compose down 2>/dev/null || true
print_success "Cleanup complete"

# Start Docker services
print_step "Starting Docker services..."
docker-compose up -d dynamodb
sleep 5  # Wait for DynamoDB to be ready
print_success "DynamoDB started"

# Initialize database
print_step "Initializing DynamoDB tables..."
docker-compose up init-db
print_success "Database initialized"

# Start API services
print_step "Starting API services..."
docker-compose up -d api-lessons-create api-lessons-get dynamodb-admin
print_success "API services started"

# Run tests
print_step "Running tests..."
make test || print_warning "Some tests failed, please review"

# Print service URLs
print_step "Development environment is ready!"
echo ""
echo -e "${GREEN}Services:${NC}"
echo "  - DynamoDB:       http://localhost:8000"
echo "  - DynamoDB Admin: http://localhost:8001"
echo "  - Create Lesson:  http://localhost:8080"
echo "  - Get Lesson:     http://localhost:8081"
echo ""
echo -e "${GREEN}Useful commands:${NC}"
echo "  - make build             Build backend"
echo "  - make test              Run tests"
echo "  - make run               Start services"
echo "  - make stop              Stop services"
echo "  - ./scripts/version.sh   Manage versions"
echo "  - docker-compose logs -f View logs"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "  1. Update .env file with your Auth0 credentials"
echo "  2. Set up Flutter mobile app: cd mobile && flutter pub get"
echo "  3. Review documentation in docs/"
echo ""
echo -e "${GREEN}Setup complete! Happy coding! 🚀${NC}"
