# DocFlow

A bidirectional document converter for Markdown and PDF files.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.25+-00ADD8.svg)
![Node Version](https://img.shields.io/badge/node-22+-339933.svg)

## Why This Exists

I got frustrated with existing MD to PDF converters. Most of them either:

- Don't let you merge multiple markdown files into a single PDF
- Don't give you control over the merge order
- Have clunky UIs that make simple tasks painful
- Charge money for basic features

So I built this. It's simple, it works, and it does exactly what I needed.

## What It Does

### MD → PDF
- Upload multiple `.md` files at once
- Drag and drop to set the order
- Preview before converting
- Merge all files into one PDF or convert separately
- Clean A4 output with proper formatting

### PDF → MD
- Extract text from PDF files
- Auto-detect headers and lists
- Get clean markdown output
- Preview extraction before downloading

### General
- Light/Dark theme
- No account needed
- No file size limits (within reason)
- Self-hostable

## Quick Start

### Prerequisites

- Go 1.25+
- Node.js 22+
- Chrome/Chromium (for PDF generation)
- poppler-utils (for PDF extraction)

```bash
# Ubuntu/Debian
sudo apt install chromium-browser poppler-utils

# macOS
brew install --cask google-chrome
brew install poppler

# Arch
sudo pacman -S chromium poppler
```

### Development

```bash
# Clone
git clone https://github.com/yourusername/docflow.git
cd docflow

# Backend
cd backend
go mod tidy
go run cmd/server/main.go

# Frontend (new terminal)
cd frontend
npm install
npm run dev
```

Open `http://localhost:3000`

### Production Build

```bash
# Build everything
make build

# Or use the build script
./scripts/build.sh

# Output:
# - bin/docflow-server (backend binary)
# - frontend/dist/ (static files)
```

### Docker

```bash
docker-compose up -d
```

## Project Structure

```
docflow/
├── backend/
│   ├── cmd/server/          # Entry point
│   ├── internal/
│   │   ├── config/          # Configuration
│   │   ├── handler/         # HTTP handlers
│   │   ├── middleware/      # CORS, logging, etc.
│   │   ├── model/           # Data structures
│   │   └── service/         # Business logic
│   └── pkg/pdf/             # PDF template
├── frontend/
│   ├── src/
│   │   ├── components/      # React components
│   │   ├── hooks/           # Custom hooks
│   │   ├── services/        # API client
│   │   └── utils/           # Helpers
│   └── package.json
├── scripts/                 # Build & deploy scripts
├── docker-compose.yml
└── Makefile
```

## API

### MD → PDF

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/health` | GET | Health check |
| `/api/preview` | POST | MD to HTML preview |
| `/api/preview/merge` | POST | Merged preview |
| `/api/convert` | POST | Convert to PDF |

### PDF → MD

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/pdf/preview` | POST | Extraction preview |
| `/api/pdf/extract` | POST | Extract to MD |

### Examples

**Convert MD to PDF:**
```bash
curl -X POST http://localhost:8080/api/convert \
  -H "Content-Type: application/json" \
  -d '{
    "files": [{"id": "1", "name": "doc.md", "content": "# Hello", "order": 0}],
    "mergeMode": "separate"
  }'
```

**Extract PDF to MD:**
```bash
curl -X POST http://localhost:8080/api/pdf/extract \
  -H "Content-Type: application/json" \
  -d '{
    "fileName": "doc.pdf",
    "content": "<base64-encoded-pdf>"
  }'
```

## Deployment

See [deployment guide](docs/DEPLOYMENT.md) for detailed instructions on:

- Bare metal deployment
- Docker deployment
- Nginx configuration
- SSL setup
- Systemd service

Quick version:

```bash
# Build
./scripts/build.sh

# Copy to server
scp -r bin/ user@server:/opt/docflow/
scp -r frontend/dist/ user@server:/opt/docflow/static/

# On server
/opt/docflow/bin/docflow-server
```

## Configuration

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_HOST` | `0.0.0.0` | Server bind address |
| `SERVER_PORT` | `8080` | Server port |
| `STORAGE_TEMP_DIR` | `./temp` | Temp file directory |
| `STORAGE_OUTPUT_DIR` | `./output` | Output directory |

## Roadmap

See [ROADMAP.md](ROADMAP.md) for planned features.

**Coming soon:**
- [ ] Batch processing CLI
- [ ] Custom PDF templates
- [ ] DOCX support
- [ ] API rate limiting
- [ ] File encryption

**Maybe later:**
- [ ] OCR for scanned PDFs
- [ ] Collaborative editing
- [ ] Cloud storage integration

## Contributing

Contributions are welcome. Please read [CONTRIBUTING.md](CONTRIBUTING.md) first.

**Important:** All contributions must include tests. PRs without tests will not be merged.

## Tech Stack

**Backend:** Go, Fiber, Goldmark, chromedp, Zap  
**Frontend:** React, Vite, Tailwind CSS, Framer Motion  
**PDF:** Headless Chrome, poppler-utils


## To Run
### Dependent packages
- chromium-browser
- poppler-utils
```bash
brew install --cask google-chrome
```

#### Frontend
```bash
cd frontend && npm install && npm run dev   
```

#### Backend
```bash
 cd backend && go mod tidy && go run cmd/server/main.go  
```

## License

MIT - do whatever you want with it.
[Licence](LICENSE)

## Acknowledgments

- [Goldmark](https://github.com/yuin/goldmark) - Markdown parser
- [chromedp](https://github.com/chromedp/chromedp) - Chrome DevTools Protocol
- [Fiber](https://github.com/gofiber/fiber) - Web framework
