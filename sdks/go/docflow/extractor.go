package docflow

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/xgaslan/docflow/sdks/go/docflow/storage"
)

// Extractor handles PDF to Markdown extraction.
type Extractor struct {
	options Options
	storage storage.Storage
}

// NewExtractor creates a new Extractor instance.
func NewExtractor(opts ...ExtractorOption) *Extractor {
	e := &Extractor{
		options: DefaultOptions(),
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

// ExtractorOption is a function that configures an Extractor.
type ExtractorOption func(*Extractor)

// WithExtractorOptions sets the extractor options.
func WithExtractorOptions(opts Options) ExtractorOption {
	return func(e *Extractor) {
		e.options = opts
	}
}

// WithExtractorStorage sets the storage backend for the extractor.
func WithExtractorStorage(s storage.Storage) ExtractorOption {
	return func(e *Extractor) {
		e.storage = s
	}
}

// ExtractToMarkdown extracts text from PDF and converts to Markdown.
func (e *Extractor) ExtractToMarkdown(ctx context.Context, pdfData []byte, filename string) (*ExtractResult, error) {
	if len(pdfData) == 0 {
		return nil, fmt.Errorf("PDF data is required")
	}

	// Ensure temp directory exists
	if err := os.MkdirAll(e.options.TempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	timestamp := time.Now().Unix()
	baseName := strings.TrimSuffix(filename, filepath.Ext(filename))
	safeName := sanitizeFilename(baseName)

	// Write PDF to temp file
	tempPDFPath := filepath.Join(e.options.TempDir, fmt.Sprintf("%s_%d.pdf", safeName, timestamp))
	if err := os.WriteFile(tempPDFPath, pdfData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write temp PDF: %w", err)
	}
	defer os.Remove(tempPDFPath)

	// Extract text
	text, err := e.extractWithPdftotext(ctx, tempPDFPath)
	if err != nil {
		// Try alternative method
		text, err = e.extractBasic(ctx, tempPDFPath)
		if err != nil {
			return &ExtractResult{Success: false, Error: err}, nil
		}
	}

	// Get page count
	pageCount, _ := e.getPageCount(ctx, tempPDFPath)

	// Convert to markdown
	markdown := e.textToMarkdown(text, baseName)

	// Save to storage if configured
	var outputPath string
	if e.storage != nil {
		outputPath = fmt.Sprintf("%s_%d.md", safeName, timestamp)
		if err := e.storage.Save(outputPath, []byte(markdown)); err != nil {
			return nil, fmt.Errorf("failed to save markdown: %w", err)
		}
		outputPath = e.storage.GetURL(outputPath)
	}

	return &ExtractResult{
		Success:   true,
		Markdown:  markdown,
		FilePath:  outputPath,
		PageCount: pageCount,
	}, nil
}

// ExtractFromFile extracts markdown from a PDF file path.
func (e *Extractor) ExtractFromFile(ctx context.Context, path string) (*ExtractResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF file: %w", err)
	}

	return e.ExtractToMarkdown(ctx, data, filepath.Base(path))
}

// GetPageCount returns the number of pages in a PDF.
func (e *Extractor) GetPageCount(ctx context.Context, pdfData []byte) (int, error) {
	// Write to temp file
	tempPath := filepath.Join(e.options.TempDir, fmt.Sprintf("pagecount_%d.pdf", time.Now().UnixNano()))
	if err := os.WriteFile(tempPath, pdfData, 0644); err != nil {
		return 0, err
	}
	defer os.Remove(tempPath)

	return e.getPageCount(ctx, tempPath)
}

// Preview extracts a preview of the first page.
func (e *Extractor) Preview(ctx context.Context, pdfData []byte, filename string) (*ExtractResult, error) {
	if len(pdfData) == 0 {
		return nil, fmt.Errorf("PDF data is required")
	}

	// Ensure temp directory exists
	if err := os.MkdirAll(e.options.TempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	timestamp := time.Now().Unix()
	baseName := strings.TrimSuffix(filename, filepath.Ext(filename))
	safeName := sanitizeFilename(baseName)

	// Write PDF to temp file
	tempPDFPath := filepath.Join(e.options.TempDir, fmt.Sprintf("%s_%d_preview.pdf", safeName, timestamp))
	if err := os.WriteFile(tempPDFPath, pdfData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write temp PDF: %w", err)
	}
	defer os.Remove(tempPDFPath)

	// Get page count
	pageCount, _ := e.getPageCount(ctx, tempPDFPath)

	// Extract first page only
	text, err := e.extractFirstPage(ctx, tempPDFPath)
	if err != nil {
		text = "Preview not available"
	}

	markdown := e.textToMarkdown(text, baseName)

	// Truncate for preview
	if len(markdown) > 2000 {
		markdown = markdown[:2000] + "\n\n... (continued)"
	}

	return &ExtractResult{
		Success:   true,
		Markdown:  markdown,
		PageCount: pageCount,
	}, nil
}

func (e *Extractor) extractWithPdftotext(ctx context.Context, pdfPath string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, e.options.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "pdftotext", "-layout", "-enc", "UTF-8", pdfPath, "-")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("pdftotext failed: %w", err)
	}

	return string(output), nil
}

func (e *Extractor) extractBasic(ctx context.Context, pdfPath string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, e.options.Timeout)
	defer cancel()

	// Try mutool
	cmd := exec.CommandContext(ctx, "mutool", "draw", "-F", "txt", pdfPath)
	output, err := cmd.Output()
	if err == nil {
		return string(output), nil
	}

	return "", fmt.Errorf("no PDF extraction tool available (install poppler-utils or mupdf-tools)")
}

func (e *Extractor) extractFirstPage(ctx context.Context, pdfPath string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "pdftotext", "-f", "1", "-l", "1", "-layout", "-enc", "UTF-8", pdfPath, "-")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func (e *Extractor) getPageCount(ctx context.Context, pdfPath string) (int, error) {
	cmd := exec.CommandContext(ctx, "pdfinfo", pdfPath)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Pages:") {
			var pages int
			fmt.Sscanf(line, "Pages: %d", &pages)
			return pages, nil
		}
	}

	return 0, fmt.Errorf("could not determine page count")
}

func (e *Extractor) textToMarkdown(text, title string) string {
	var result strings.Builder

	// Add title
	result.WriteString("# ")
	result.WriteString(title)
	result.WriteString("\n\n")

	// Process text
	lines := strings.Split(text, "\n")
	var processedLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			if len(processedLines) > 0 && processedLines[len(processedLines)-1] != "" {
				processedLines = append(processedLines, "")
			}
			continue
		}

		// Detect headers
		if e.isPotentialHeader(line) {
			line = "## " + strings.Title(strings.ToLower(line))
		}

		// Detect bullet points
		if strings.HasPrefix(line, "•") || strings.HasPrefix(line, "●") || strings.HasPrefix(line, "○") {
			rest := strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(line, "•"), "●"), "○")
			line = "- " + strings.TrimSpace(rest)
		}

		// Detect numbered lists (already formatted, keep as is)
		if matched, _ := regexp.MatchString(`^\d+[\.]\s`, line); matched {
			// Keep numbered list format
		}

		processedLines = append(processedLines, line)
	}

	// Join with proper spacing
	inParagraph := false
	for i, line := range processedLines {
		if line == "" {
			if inParagraph {
				result.WriteString("\n\n")
				inParagraph = false
			}
			continue
		}

		if strings.HasPrefix(line, "##") || strings.HasPrefix(line, "- ") {
			if inParagraph {
				result.WriteString("\n\n")
			}
			result.WriteString(line)
			result.WriteString("\n")
			inParagraph = false
		} else {
			if inParagraph && i > 0 && processedLines[i-1] != "" {
				result.WriteString(" ")
			}
			result.WriteString(line)
			inParagraph = true
		}
	}

	result.WriteString("\n")
	return result.String()
}

func (e *Extractor) isPotentialHeader(line string) bool {
	if len(line) > 3 && len(line) < 60 && line == strings.ToUpper(line) {
		words := strings.Fields(line)
		if len(words) >= 2 {
			return true
		}
	}
	return false
}
