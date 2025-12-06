# Roadmap

This document outlines planned features and improvements. Things may change based on feedback and priorities.

## Current Version (v1.0)

What's working now:

- [x] MD to PDF conversion
- [x] Multiple file upload
- [x] Drag-and-drop ordering
- [x] Merge files into single PDF
- [x] PDF to MD extraction
- [x] Live preview
- [x] Light/Dark theme
- [x] Docker support

## Short Term (v1.1)

Target: Next 1-2 months

### Features
- [ ] **Custom PDF templates** - Let users choose different styles (academic, minimal, corporate)
- [ ] **CLI tool** - Command line interface for batch processing
- [ ] **Syntax highlighting** - Better code block rendering in PDFs
- [ ] **Table of contents** - Auto-generate TOC from headings

### Improvements
- [ ] **Better PDF extraction** - Improve header/list detection accuracy
- [ ] **Keyboard shortcuts** - Quick actions for power users
- [ ] **Drag files to reorder** - Visual feedback improvements

### Tech Debt
- [ ] Increase test coverage to 80%+
- [ ] Add E2E tests with Playwright
- [ ] Set up CI/CD pipeline

## Medium Term (v1.2)

Target: 3-6 months

### Features
- [ ] **DOCX support** - Convert to/from Word documents
- [ ] **Image handling** - Embed images in markdown, extract from PDF
- [ ] **Custom fonts** - Upload and use custom fonts in PDFs
- [ ] **Watermarks** - Add watermarks to generated PDFs
- [ ] **Password protection** - Encrypt output PDFs

### Improvements
- [ ] **API rate limiting** - Prevent abuse
- [ ] **File size optimization** - Compress output PDFs
- [ ] **Progress indicators** - Show conversion progress for large files

## Long Term (v2.0)

Target: 6-12 months

### Features
- [ ] **OCR support** - Extract text from scanned PDFs
- [ ] **Collaborative editing** - Real-time markdown editing (maybe)
- [ ] **Cloud storage** - Save to Google Drive, Dropbox, etc.
- [ ] **Templates marketplace** - Share and download templates
- [ ] **API keys** - For integrating with other services

### Platform
- [ ] **Desktop app** - Electron or Tauri wrapper
- [ ] **VS Code extension** - Convert directly from editor
- [ ] **Browser extension** - Quick conversion from any page

## Maybe Someday

Ideas that might happen if there's demand:

- Mobile app
- LaTeX support
- Presentation mode (MD to slides)
- Version history
- Team workspaces
- Self-hosted cloud version

## Won't Do

Things that are out of scope:

- Full word processor features
- Real-time collaboration (complex, many solutions exist)
- DRM/copy protection
- Paid tiers (keeping it free and open source)

## How Features Get Prioritized

1. **Community demand** - Most requested features get priority
2. **Complexity vs value** - Quick wins over complex features
3. **Maintainability** - Must be testable and maintainable
4. **Alignment** - Must fit the project's purpose

## Want to Suggest Something?

Open an issue with the `feature-request` label. Include:

- What you want
- Why you need it
- How you'd use it

Good suggestions with clear use cases get prioritized.

## Contributing to Roadmap Items

Want to work on something from this list? Great!

1. Check if there's an existing issue
2. If not, create one and mention you want to work on it
3. Wait for confirmation (to avoid duplicate work)
4. Start coding

See [CONTRIBUTING.md](CONTRIBUTING.md) for details.
