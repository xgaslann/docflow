# Deployment Guide

This guide covers different ways to deploy DocFlow.

## Table of Contents

- [Requirements](#requirements)
- [Quick Deploy](#quick-deploy)
- [Manual Deployment](#manual-deployment)
- [Docker Deployment](#docker-deployment)
- [Nginx Configuration](#nginx-configuration)
- [SSL Setup](#ssl-setup)
- [Monitoring](#monitoring)
- [Troubleshooting](#troubleshooting)

## Requirements

### System
- Linux server (Ubuntu 22.04+ recommended)
- 1 CPU core minimum (2+ recommended)
- 1GB RAM minimum (2GB+ recommended)
- 10GB disk space

### Software
- Chrome/Chromium (for PDF generation)
- poppler-utils (for PDF extraction)

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install -y chromium-browser poppler-utils

# Verify
chromium-browser --version
pdftotext -v
```

## Quick Deploy

If you just want to get it running:

```bash
# On your local machine
./scripts/build.sh package

# Copy to server
scp build/docflow-*.tar.gz user@server:/tmp/

# On server
cd /opt
sudo tar -xzf /tmp/docflow-*.tar.gz
sudo mv docflow-release docflow

# Create user
sudo useradd -r -m -s /bin/bash docflow
sudo chown -R docflow:docflow /opt/docflow

# Start
cd /opt/docflow
./start.sh
```

## Manual Deployment

### Step 1: Prepare Server

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install dependencies
sudo apt install -y \
    chromium-browser \
    poppler-utils \
    nginx \
    certbot \
    python3-certbot-nginx

# Create user
sudo useradd -r -m -s /bin/bash docflow
```

### Step 2: Build

On your development machine:

```bash
# Build everything
./scripts/build.sh all

# Or build release package
./scripts/build.sh package
```

### Step 3: Transfer Files

```bash
# Create directories
sudo mkdir -p /opt/docflow/{bin,static,data/temp,data/output}

# Copy binary
scp bin/docflow-server user@server:/opt/docflow/bin/

# Copy frontend
scp -r frontend/dist/* user@server:/opt/docflow/static/

# Set permissions
sudo chown -R docflow:docflow /opt/docflow
sudo chmod +x /opt/docflow/bin/docflow-server
```

### Step 4: Configure

```bash
# Create config
sudo -u docflow cat > /opt/docflow/.env << 'EOF'
SERVER_HOST=127.0.0.1
SERVER_PORT=8080
STORAGE_TEMP_DIR=/opt/docflow/data/temp
STORAGE_OUTPUT_DIR=/opt/docflow/data/output
EOF
```

### Step 5: Systemd Service

```bash
# Copy service file
sudo cp scripts/docflow.service /etc/systemd/system/

# Enable and start
sudo systemctl daemon-reload
sudo systemctl enable docflow
sudo systemctl start docflow

# Check status
sudo systemctl status docflow
```

### Step 6: Verify

```bash
# Check if running
curl http://localhost:8080/api/health

# Check logs
sudo journalctl -u docflow -f
```

## Docker Deployment

### Using Docker Compose

```bash
# Clone repo
git clone https://github.com/yourusername/docflow.git
cd docflow

# Start
docker-compose up -d

# Check logs
docker-compose logs -f
```

### Custom Docker Build

```bash
# Build images
docker build -t docflow-backend ./backend
docker build -t docflow-frontend ./frontend

# Run
docker network create docflow-net

docker run -d \
    --name docflow-api \
    --network docflow-net \
    -v docflow-data:/app/data \
    docflow-backend

docker run -d \
    --name docflow-web \
    --network docflow-net \
    -p 80:80 \
    docflow-frontend
```

## Nginx Configuration

### Basic Setup

```nginx
# /etc/nginx/sites-available/docflow
server {
    listen 80;
    server_name your-domain.com;

    # Frontend static files
    root /opt/docflow/static;
    index index.html;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;

    # Gzip
    gzip on;
    gzip_types text/plain text/css application/json application/javascript;

    # Frontend routing
    location / {
        try_files $uri $uri/ /index.html;
    }

    # API proxy
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Timeout for PDF generation
        proxy_read_timeout 120s;
        proxy_send_timeout 120s;
        
        # File upload limit
        client_max_body_size 50M;
    }

    # PDF files
    location /output/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
    }

    # Static file caching
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2)$ {
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
}
```

Enable site:

```bash
sudo ln -sf /etc/nginx/sites-available/docflow /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default
sudo nginx -t
sudo systemctl reload nginx
```

## SSL Setup

Using Let's Encrypt:

```bash
# Get certificate
sudo certbot --nginx -d your-domain.com

# Auto-renewal test
sudo certbot renew --dry-run
```

## Monitoring

### Logs

```bash
# Application logs
sudo journalctl -u docflow -f

# Nginx access logs
sudo tail -f /var/log/nginx/access.log

# Nginx error logs
sudo tail -f /var/log/nginx/error.log
```

### Health Check

```bash
# Simple check
curl http://localhost:8080/api/health

# With monitoring tool
cat > /opt/docflow/healthcheck.sh << 'EOF'
#!/bin/bash
response=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/health)
if [ "$response" != "200" ]; then
    echo "Health check failed: $response"
    systemctl restart docflow
fi
EOF

chmod +x /opt/docflow/healthcheck.sh

# Add to crontab
echo "*/5 * * * * /opt/docflow/healthcheck.sh" | sudo crontab -
```

### Cleanup Old Files

```bash
# Add to crontab
cat > /etc/cron.daily/docflow-cleanup << 'EOF'
#!/bin/bash
find /opt/docflow/data/output -type f -mtime +7 -delete
find /opt/docflow/data/temp -type f -mtime +1 -delete
EOF

chmod +x /etc/cron.daily/docflow-cleanup
```

## Troubleshooting

### Service Won't Start

```bash
# Check logs
sudo journalctl -u docflow -n 50 --no-pager

# Check permissions
ls -la /opt/docflow/
ls -la /opt/docflow/bin/

# Verify binary
file /opt/docflow/bin/docflow-server
```

### PDF Generation Fails

```bash
# Check Chrome
which chromium-browser
chromium-browser --version

# Test Chrome headless
chromium-browser --headless --disable-gpu --print-to-pdf=/tmp/test.pdf https://example.com

# Check sandbox issues (if needed)
# Add to systemd service:
# Environment="CHROME_DEVEL_SANDBOX="
```

### PDF Extraction Fails

```bash
# Check poppler
which pdftotext
pdftotext -v

# Test extraction
pdftotext /path/to/test.pdf -
```

### Permission Issues

```bash
# Fix ownership
sudo chown -R docflow:docflow /opt/docflow

# Fix permissions
sudo chmod -R 755 /opt/docflow
sudo chmod -R 775 /opt/docflow/data
```

### Memory Issues

```bash
# Check memory
free -h

# Increase limit in systemd
sudo systemctl edit docflow

# Add:
[Service]
MemoryMax=2G

# Restart
sudo systemctl restart docflow
```

### Port Conflicts

```bash
# Check what's using port 8080
sudo lsof -i :8080

# Change port in .env
echo "SERVER_PORT=8081" >> /opt/docflow/.env
sudo systemctl restart docflow
```

## Updating

```bash
# Stop service
sudo systemctl stop docflow

# Backup current version
sudo cp /opt/docflow/bin/docflow-server /opt/docflow/bin/docflow-server.bak

# Upload new version
scp bin/docflow-server user@server:/opt/docflow/bin/
scp -r frontend/dist/* user@server:/opt/docflow/static/

# Start
sudo systemctl start docflow

# Verify
curl http://localhost:8080/api/health
```

## Rollback

```bash
# Stop service
sudo systemctl stop docflow

# Restore backup
sudo mv /opt/docflow/bin/docflow-server.bak /opt/docflow/bin/docflow-server

# Start
sudo systemctl start docflow
```
