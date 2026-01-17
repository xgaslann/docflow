package service

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gorkem/md-to-pdf/internal/config"
	"github.com/gorkem/md-to-pdf/internal/model"
	"go.uber.org/zap"
)

// PDFExtractorService handles PDF to Markdown conversion
type PDFExtractorService struct {
	cfg    *config.Config
	logger *zap.Logger
}

// NewPDFExtractorService creates a new PDF extractor service
func NewPDFExtractorService(cfg *config.Config, logger *zap.Logger) *PDFExtractorService {
	return &PDFExtractorService{
		cfg:    cfg,
		logger: logger,
	}
}

// ExtractToMarkdown extracts text from PDF and converts to Markdown
func (s *PDFExtractorService) ExtractToMarkdown(ctx context.Context, pdfData []byte, filename string) (*model.PDFExtractResponse, error) {
	timestamp := time.Now().Unix()
	baseName := strings.TrimSuffix(filename, filepath.Ext(filename))
	safeName := sanitizeFilename(baseName)

	// Write PDF to temp file
	tempPDFPath := filepath.Join(s.cfg.Storage.TempDir, fmt.Sprintf("%s_%d.pdf", safeName, timestamp))
	if err := os.WriteFile(tempPDFPath, pdfData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write temp PDF: %w", err)
	}
	defer os.Remove(tempPDFPath)

	// Extract text using pdftotext
	text, err := s.extractWithPdftotext(ctx, tempPDFPath)
	if err != nil {
		s.logger.Warn("pdftotext failed, trying alternative method", zap.Error(err))
		// Fallback to basic extraction
		text, err = s.extractBasic(ctx, tempPDFPath)
		if err != nil {
			return nil, fmt.Errorf("PDF extraction failed: %w", err)
		}
	}

	// Convert extracted text to Markdown
	markdown := s.textToMarkdown(text, baseName)

	// Save markdown file
	outputName := fmt.Sprintf("%s_%d.md", safeName, timestamp)
	outputPath := filepath.Join(s.cfg.Storage.OutputDir, outputName)
	if err := os.WriteFile(outputPath, []byte(markdown), 0644); err != nil {
		return nil, fmt.Errorf("failed to write markdown file: %w", err)
	}

	s.logger.Info("PDF extracted successfully",
		zap.String("input", filename),
		zap.String("output", outputName),
		zap.Int("textLength", len(text)),
	)

	return &model.PDFExtractResponse{
		Success:  true,
		Markdown: markdown,
		FilePath: "/output/" + outputName,
		FileName: outputName,
	}, nil
}

// extractWithPdftotext uses pdftotext command
func (s *PDFExtractorService) extractWithPdftotext(ctx context.Context, pdfPath string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// pdftotext -layout preserves layout better
	cmd := exec.CommandContext(ctx, "pdftotext", "-layout", "-enc", "UTF-8", pdfPath, "-")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%w: %s", err, stderr.String())
	}

	return stdout.String(), nil
}

// extractBasic tries alternative extraction methods
func (s *PDFExtractorService) extractBasic(ctx context.Context, pdfPath string) (string, error) {
	// Try mutool (MuPDF)
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "mutool", "draw", "-F", "txt", pdfPath)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err == nil {
		return stdout.String(), nil
	}

	// If all methods fail, return error
	return "", fmt.Errorf("no PDF extraction tool available (install poppler-utils or mupdf-tools)")
}

// textToMarkdown converts extracted text to Markdown format
func (s *PDFExtractorService) textToMarkdown(text, title string) string {
	var result strings.Builder

	// Add title
	result.WriteString("# ")
	result.WriteString(title)
	result.WriteString("\n\n")

	// Clean and process text
	lines := strings.Split(text, "\n")
	var processedLines []string

	for _, line := range lines {
		// Trim whitespace
		line = strings.TrimSpace(line)

		// Skip empty lines (will add them back appropriately)
		if line == "" {
			if len(processedLines) > 0 && processedLines[len(processedLines)-1] != "" {
				processedLines = append(processedLines, "")
			}
			continue
		}

		// Detect potential headers (ALL CAPS lines, short lines that look like titles)
		if s.isPotentialHeader(line) {
			line = "## " + strings.Title(strings.ToLower(line))
		}

		// Detect bullet points
		if strings.HasPrefix(line, "•") || strings.HasPrefix(line, "●") || strings.HasPrefix(line, "○") {
			line = "- " + strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(line, "•"), "●"), "○")
			line = strings.TrimSpace(line)
			line = "- " + line[2:]
		}

		// Detect numbered lists
		if matched, _ := regexp.MatchString(`^\d+[\.\)]\s`, line); matched {
			// Already looks like a numbered list, keep it
		}

		processedLines = append(processedLines, line)
	}

	// Join lines with proper spacing
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

// isPotentialHeader detects if a line might be a header
func (s *PDFExtractorService) isPotentialHeader(line string) bool {
	// All caps and short
	if len(line) > 3 && len(line) < 60 && line == strings.ToUpper(line) {
		// Check if it's not just an acronym
		words := strings.Fields(line)
		if len(words) >= 2 {
			return true
		}
	}
	return false
}

// PreviewExtraction returns preview of extracted content
func (s *PDFExtractorService) PreviewExtraction(ctx context.Context, pdfData []byte, filename string) (*model.PDFPreviewResponse, error) {
	timestamp := time.Now().Unix()
	baseName := strings.TrimSuffix(filename, filepath.Ext(filename))
	safeName := sanitizeFilename(baseName)

	// Write PDF to temp file
	tempPDFPath := filepath.Join(s.cfg.Storage.TempDir, fmt.Sprintf("%s_%d_preview.pdf", safeName, timestamp))
	if err := os.WriteFile(tempPDFPath, pdfData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write temp PDF: %w", err)
	}
	defer os.Remove(tempPDFPath)

	// Get page count
	pageCount, _ := s.getPageCount(ctx, tempPDFPath)

	// Extract first page only for preview
	text, err := s.extractFirstPage(ctx, tempPDFPath)
	if err != nil {
		text = "Preview not available"
	}

	// Convert to markdown (preview)
	markdown := s.textToMarkdown(text, baseName)

	// Truncate for preview
	if len(markdown) > 2000 {
		markdown = markdown[:2000] + "\n\n... (devamı var)"
	}

	return &model.PDFPreviewResponse{
		Preview:   markdown,
		PageCount: pageCount,
		FileName:  filename,
	}, nil
}

// getPageCount returns the number of pages in PDF
func (s *PDFExtractorService) getPageCount(ctx context.Context, pdfPath string) (int, error) {
	cmd := exec.CommandContext(ctx, "pdfinfo", pdfPath)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return 0, err
	}

	// Parse output for "Pages:" line
	lines := strings.Split(stdout.String(), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Pages:") {
			var pages int
			fmt.Sscanf(line, "Pages: %d", &pages)
			return pages, nil
		}
	}

	return 0, fmt.Errorf("could not determine page count")
}

// extractFirstPage extracts only the first page
func (s *PDFExtractorService) extractFirstPage(ctx context.Context, pdfPath string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "pdftotext", "-f", "1", "-l", "1", "-layout", "-enc", "UTF-8", pdfPath, "-")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return stdout.String(), nil
}
