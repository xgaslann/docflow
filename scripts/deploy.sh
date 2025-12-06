#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Default values
DEPLOY_USER="${DEPLOY_USER:-docflow}"
DEPLOY_HOST="${DEPLOY_HOST:-}"
DEPLOY_PATH="${DEPLOY_PATH:-/opt/docflow}"
SSH_KEY="${SSH_KEY:-}"

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

usage() {
    echo "Usage: $0 --host <hostname> [options]"
    echo ""
    echo "Options:"
    echo "  --host     Target server hostname (required)"
    echo "  --user     SSH user (default: docflow)"
    echo "  --path     Deployment path (default: /opt/docflow)"
    echo "  --key      SSH key file path"
    echo "  --help     Show this help"
    echo ""
    echo "Environment variables:"
    echo "  DEPLOY_HOST  Target server hostname"
    echo "  DEPLOY_USER  SSH user"
    echo "  DEPLOY_PATH  Deployment path"
    echo "  SSH_KEY      SSH key file path"
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --host)
            DEPLOY_HOST="$2"
            shift 2
            ;;
        --user)
            DEPLOY_USER="$2"
            shift 2
            ;;
        --path)
            DEPLOY_PATH="$2"
            shift 2
            ;;
        --key)
            SSH_KEY="$2"
            shift 2
            ;;
        --help|-h)
            usage
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            usage
            exit 1
            ;;
    esac
done

# Validate
if [[ -z "$DEPLOY_HOST" ]]; then
    echo -e "${RED}Error: --host is required${NC}"
    usage
    exit 1
fi

# Build SSH command
SSH_CMD="ssh"
SCP_CMD="scp"
if [[ -n "$SSH_KEY" ]]; then
    SSH_CMD="ssh -i $SSH_KEY"
    SCP_CMD="scp -i $SSH_KEY"
fi

SSH_TARGET="${DEPLOY_USER}@${DEPLOY_HOST}"

echo -e "${GREEN}DocFlow Deploy Script${NC}"
echo "====================="
echo "Target: ${SSH_TARGET}:${DEPLOY_PATH}"

# Check if build exists
if [[ ! -f "${PROJECT_ROOT}/bin/docflow-server" ]]; then
    echo -e "${YELLOW}Build not found. Running build script...${NC}"
    "${PROJECT_ROOT}/scripts/build.sh" all
fi

if [[ ! -d "${PROJECT_ROOT}/frontend/dist" ]]; then
    echo -e "${RED}Frontend build not found${NC}"
    exit 1
fi

# Create directories on server
echo -e "\n${YELLOW}Creating directories on server...${NC}"
$SSH_CMD $SSH_TARGET << EOF
    sudo mkdir -p ${DEPLOY_PATH}/{bin,static,data/temp,data/output}
    sudo chown -R ${DEPLOY_USER}:${DEPLOY_USER} ${DEPLOY_PATH}
EOF

# Stop service if running
echo -e "\n${YELLOW}Stopping service...${NC}"
$SSH_CMD $SSH_TARGET "sudo systemctl stop docflow 2>/dev/null || true"

# Upload binary
echo -e "\n${YELLOW}Uploading backend...${NC}"
$SCP_CMD "${PROJECT_ROOT}/bin/docflow-server" "${SSH_TARGET}:${DEPLOY_PATH}/bin/"

# Upload frontend
echo -e "\n${YELLOW}Uploading frontend...${NC}"
$SCP_CMD -r "${PROJECT_ROOT}/frontend/dist/"* "${SSH_TARGET}:${DEPLOY_PATH}/static/"

# Create/update env file if not exists
echo -e "\n${YELLOW}Checking config...${NC}"
$SSH_CMD $SSH_TARGET << EOF
    if [[ ! -f ${DEPLOY_PATH}/.env ]]; then
        cat > ${DEPLOY_PATH}/.env << 'ENVFILE'
SERVER_HOST=127.0.0.1
SERVER_PORT=8080
STORAGE_TEMP_DIR=${DEPLOY_PATH}/data/temp
STORAGE_OUTPUT_DIR=${DEPLOY_PATH}/data/output
ENVFILE
        echo "Created default .env file"
    else
        echo ".env file exists, keeping current config"
    fi
EOF

# Set permissions
$SSH_CMD $SSH_TARGET "chmod +x ${DEPLOY_PATH}/bin/docflow-server"

# Install systemd service
echo -e "\n${YELLOW}Installing systemd service...${NC}"
$SSH_CMD $SSH_TARGET << EOF
    sudo tee /etc/systemd/system/docflow.service > /dev/null << 'SERVICE'
[Unit]
Description=DocFlow Document Converter
After=network.target

[Service]
Type=simple
User=${DEPLOY_USER}
Group=${DEPLOY_USER}
WorkingDirectory=${DEPLOY_PATH}
EnvironmentFile=${DEPLOY_PATH}/.env
ExecStart=${DEPLOY_PATH}/bin/docflow-server
Restart=always
RestartSec=5

# Security
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=${DEPLOY_PATH}/data

[Install]
WantedBy=multi-user.target
SERVICE

    sudo systemctl daemon-reload
EOF

# Start service
echo -e "\n${YELLOW}Starting service...${NC}"
$SSH_CMD $SSH_TARGET "sudo systemctl enable docflow && sudo systemctl start docflow"

# Wait and check status
sleep 2
echo -e "\n${YELLOW}Checking service status...${NC}"
$SSH_CMD $SSH_TARGET "sudo systemctl status docflow --no-pager"

# Health check
echo -e "\n${YELLOW}Running health check...${NC}"
$SSH_CMD $SSH_TARGET "curl -s http://localhost:8080/api/health | head -c 100"

echo -e "\n\n${GREEN}Deployment complete!${NC}"
echo -e "Server running at: http://${DEPLOY_HOST}:8080"
echo -e "\nNext steps:"
echo "  1. Configure nginx reverse proxy"
echo "  2. Set up SSL with certbot"
echo "  3. Configure firewall"
