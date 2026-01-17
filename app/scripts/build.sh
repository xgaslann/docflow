#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Project root
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUILD_DIR="${PROJECT_ROOT}/build"
BIN_DIR="${PROJECT_ROOT}/bin"

echo -e "${GREEN}DocFlow Build Script${NC}"
echo "====================="

# Check dependencies
check_dependencies() {
    echo -e "\n${YELLOW}Checking dependencies...${NC}"
    
    if ! command -v go &> /dev/null; then
        echo -e "${RED}Error: Go is not installed${NC}"
        exit 1
    fi
    echo "✓ Go $(go version | awk '{print $3}')"
    
    if ! command -v node &> /dev/null; then
        echo -e "${RED}Error: Node.js is not installed${NC}"
        exit 1
    fi
    echo "✓ Node.js $(node --version)"
    
    if ! command -v npm &> /dev/null; then
        echo -e "${RED}Error: npm is not installed${NC}"
        exit 1
    fi
    echo "✓ npm $(npm --version)"
}

# Build backend
build_backend() {
    echo -e "\n${YELLOW}Building backend...${NC}"
    cd "${PROJECT_ROOT}/backend"
    
    # Download dependencies
    go mod tidy
    
    # Build binary
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
        -ldflags="-s -w -X main.version=$(git describe --tags --always 2>/dev/null || echo 'dev')" \
        -o "${BIN_DIR}/docflow-server" \
        cmd/server/main.go
    
    echo -e "${GREEN}✓ Backend built: ${BIN_DIR}/docflow-server${NC}"
}

# Build frontend
build_frontend() {
    echo -e "\n${YELLOW}Building frontend...${NC}"
    cd "${PROJECT_ROOT}/frontend"
    
    # Install dependencies
    npm ci --silent
    
    # Build
    npm run build
    
    echo -e "${GREEN}✓ Frontend built: ${PROJECT_ROOT}/frontend/dist/${NC}"
}

# Run tests
run_tests() {
    echo -e "\n${YELLOW}Running tests...${NC}"
    
    # Backend tests
    cd "${PROJECT_ROOT}/backend"
    echo "Running backend tests..."
    go test -v ./... || {
        echo -e "${RED}Backend tests failed${NC}"
        exit 1
    }
    
    # Frontend tests
    cd "${PROJECT_ROOT}/frontend"
    echo "Running frontend tests..."
    npm run test:run || {
        echo -e "${RED}Frontend tests failed${NC}"
        exit 1
    }
    
    echo -e "${GREEN}✓ All tests passed${NC}"
}

# Create release package
create_package() {
    echo -e "\n${YELLOW}Creating release package...${NC}"
    
    RELEASE_DIR="${BUILD_DIR}/docflow-release"
    rm -rf "${RELEASE_DIR}"
    mkdir -p "${RELEASE_DIR}"
    
    # Copy binary
    cp "${BIN_DIR}/docflow-server" "${RELEASE_DIR}/"
    
    # Copy frontend build
    cp -r "${PROJECT_ROOT}/frontend/dist" "${RELEASE_DIR}/static"
    
    # Copy config examples
    cat > "${RELEASE_DIR}/.env.example" << 'EOF'
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
STORAGE_TEMP_DIR=./data/temp
STORAGE_OUTPUT_DIR=./data/output
EOF
    
    # Copy systemd service file
    mkdir -p "${RELEASE_DIR}/systemd"
    cp "${PROJECT_ROOT}/scripts/docflow.service" "${RELEASE_DIR}/systemd/" 2>/dev/null || true
    
    # Create start script
    cat > "${RELEASE_DIR}/start.sh" << 'EOF'
#!/bin/bash
cd "$(dirname "$0")"
./docflow-server
EOF
    chmod +x "${RELEASE_DIR}/start.sh"
    
    # Create archive
    cd "${BUILD_DIR}"
    tar -czvf "docflow-$(date +%Y%m%d).tar.gz" "docflow-release"
    
    echo -e "${GREEN}✓ Package created: ${BUILD_DIR}/docflow-$(date +%Y%m%d).tar.gz${NC}"
}

# Clean build artifacts
clean() {
    echo -e "\n${YELLOW}Cleaning build artifacts...${NC}"
    rm -rf "${BUILD_DIR}"
    rm -rf "${BIN_DIR}"
    rm -rf "${PROJECT_ROOT}/frontend/dist"
    rm -rf "${PROJECT_ROOT}/frontend/node_modules"
    echo -e "${GREEN}✓ Cleaned${NC}"
}

# Print usage
usage() {
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  all       Build everything (default)"
    echo "  backend   Build backend only"
    echo "  frontend  Build frontend only"
    echo "  test      Run all tests"
    echo "  package   Create release package"
    echo "  clean     Clean build artifacts"
    echo "  help      Show this help"
}

# Main
main() {
    mkdir -p "${BUILD_DIR}"
    mkdir -p "${BIN_DIR}"
    
    case "${1:-all}" in
        all)
            check_dependencies
            build_backend
            build_frontend
            echo -e "\n${GREEN}Build complete!${NC}"
            ;;
        backend)
            check_dependencies
            build_backend
            ;;
        frontend)
            check_dependencies
            build_frontend
            ;;
        test)
            check_dependencies
            run_tests
            ;;
        package)
            check_dependencies
            build_backend
            build_frontend
            create_package
            ;;
        clean)
            clean
            ;;
        help|--help|-h)
            usage
            ;;
        *)
            echo -e "${RED}Unknown command: $1${NC}"
            usage
            exit 1
            ;;
    esac
}

main "$@"
