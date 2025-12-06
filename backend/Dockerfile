# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server cmd/server/main.go

# Runtime stage
FROM alpine:3.21

# Install chromium for PDF generation and poppler-utils for PDF extraction
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont \
    font-noto \
    font-noto-emoji \
    poppler-utils

WORKDIR /app

# Copy binary
COPY --from=builder /app/server .

# Create data directories
RUN mkdir -p /app/data/temp /app/data/output && \
    adduser -D -u 1000 appuser && \
    chown -R appuser:appuser /app

USER appuser

# Environment
ENV CHROME_PATH=/usr/bin/chromium-browser
ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=8080
ENV STORAGE_TEMP_DIR=/app/data/temp
ENV STORAGE_OUTPUT_DIR=/app/data/output

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget -q --spider http://localhost:8080/api/health || exit 1

CMD ["./server"]
