package docflow

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/xgaslan/docflow/sdks/go/docflow/storage"
)

// Converter handles Markdown to PDF conversion.
type Converter struct {
	options  Options
	storage  storage.Storage
	parser   *MarkdownParser
	template *Template
}

// NewConverter creates a new Converter instance.
func NewConverter(opts ...ConverterOption) *Converter {
	c := &Converter{
		options:  DefaultOptions(),
		parser:   NewMarkdownParser(),
		template: NewTemplate(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// ConverterOption is a function that configures a Converter.
type ConverterOption func(*Converter)

// WithOptions sets the converter options.
func WithOptions(opts Options) ConverterOption {
	return func(c *Converter) {
		c.options = opts
	}
}

// WithStorage sets the storage backend.
func WithStorage(s storage.Storage) ConverterOption {
	return func(c *Converter) {
		c.storage = s
	}
}

// WithLocalStorage creates and sets a local storage backend.
func WithLocalStorage(path string) ConverterOption {
	return func(c *Converter) {
		s, err := storage.NewLocalStorage(path)
		if err != nil {
			// Log error but continue - storage is optional
			fmt.Fprintf(os.Stderr, "docflow: failed to create local storage at %s: %v\n", path, err)
			return
		}
		c.storage = s
	}
}

// ConvertToPDF converts markdown files to PDF and saves to storage.
func (c *Converter) ConvertToPDF(ctx context.Context, files []MDFile, opts ConvertOptions) (*Result, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("at least one file is required")
	}

	// Ensure temp directory exists
	if err := os.MkdirAll(c.options.TempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	timestamp := time.Now().Unix()

	// Set default merge mode
	if opts.MergeMode == "" {
		opts.MergeMode = "separate"
	}

	var results []string
	var resultBytes []byte

	switch opts.MergeMode {
	case "merged":
		path, data, err := c.convertMerged(ctx, files, opts.OutputName, timestamp)
		if err != nil {
			return &Result{Success: false, Error: err}, nil
		}
		results = append(results, path)
		resultBytes = data

	case "separate":
		for i, file := range files {
			file.Order = i
			path, _, err := c.convertSingle(ctx, file, timestamp)
			if err != nil {
				return &Result{Success: false, Error: err}, nil
			}
			results = append(results, path)
		}

	default:
		return nil, fmt.Errorf("invalid merge mode: %s", opts.MergeMode)
	}

	return &Result{
		Success:   true,
		FilePaths: results,
		Bytes:     resultBytes,
	}, nil
}

// ConvertToBytes converts markdown to PDF and returns the bytes.
func (c *Converter) ConvertToBytes(ctx context.Context, files []MDFile) ([]byte, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("at least one file is required")
	}

	// Ensure temp directory exists
	if err := os.MkdirAll(c.options.TempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	timestamp := time.Now().Unix()

	_, data, err := c.convertMerged(ctx, files, "", timestamp)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Preview generates an HTML preview of markdown content.
func (c *Converter) Preview(content string) (string, error) {
	return c.parser.ToHTML(content)
}

func (c *Converter) convertMerged(ctx context.Context, files []MDFile, outputName string, timestamp int64) (string, []byte, error) {
	mergedContent := c.parser.MergeFiles(files)

	if outputName == "" {
		outputName = fmt.Sprintf("merged_%d", timestamp)
	}
	outputName = sanitizeFilename(outputName)

	return c.generatePDF(ctx, mergedContent, outputName)
}

func (c *Converter) convertSingle(ctx context.Context, file MDFile, timestamp int64) (string, []byte, error) {
	baseName := strings.TrimSuffix(file.Name, filepath.Ext(file.Name))
	outputName := fmt.Sprintf("%s_%d", sanitizeFilename(baseName), timestamp)

	return c.generatePDF(ctx, file.Content, outputName)
}

func (c *Converter) generatePDF(ctx context.Context, mdContent, outputName string) (string, []byte, error) {
	// Convert markdown to HTML
	htmlContent, err := c.parser.ToHTML(mdContent)
	if err != nil {
		return "", nil, fmt.Errorf("markdown conversion failed: %w", err)
	}

	// Generate full HTML document
	fullHTML := c.template.Generate(htmlContent)

	// Write HTML to temp file
	tempHTMLPath := filepath.Join(c.options.TempDir, outputName+".html")
	if err := os.WriteFile(tempHTMLPath, []byte(fullHTML), 0644); err != nil {
		return "", nil, fmt.Errorf("failed to write temp file: %w", err)
	}
	defer os.Remove(tempHTMLPath)

	// Generate PDF using Chrome
	pdfData, err := c.generateWithChrome(ctx, tempHTMLPath)
	if err != nil {
		return "", nil, err
	}

	// Save to storage if configured
	if c.storage != nil {
		outputPath := outputName + ".pdf"
		if err := c.storage.Save(outputPath, pdfData); err != nil {
			return "", nil, fmt.Errorf("failed to save PDF: %w", err)
		}
		return c.storage.GetURL(outputPath), pdfData, nil
	}

	// Save to temp directory if no storage configured
	outputPath := filepath.Join(c.options.TempDir, outputName+".pdf")
	if err := os.WriteFile(outputPath, pdfData, 0644); err != nil {
		return "", nil, fmt.Errorf("failed to write PDF: %w", err)
	}

	return outputPath, pdfData, nil
}

func (c *Converter) generateWithChrome(ctx context.Context, htmlPath string) ([]byte, error) {
	// Get absolute path
	absPath, err := filepath.Abs(htmlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Create chromedp context with options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-software-rasterizer", true),
	)

	if c.options.ChromePath != "" {
		opts = append(opts, chromedp.ExecPath(c.options.ChromePath))
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	defer allocCancel()

	taskCtx, taskCancel := chromedp.NewContext(allocCtx)
	defer taskCancel()

	// Set timeout
	taskCtx, cancel := context.WithTimeout(taskCtx, c.options.Timeout)
	defer cancel()

	// Generate PDF
	var pdfBuf []byte
	fileURL := "file://" + absPath

	if err := chromedp.Run(taskCtx,
		chromedp.Navigate(fileURL),
		chromedp.WaitReady("body"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuf, _, err = page.PrintToPDF().
				WithPaperWidth(8.27).      // A4 width in inches
				WithPaperHeight(11.69).    // A4 height in inches
				WithMarginTop(0.79).       // 20mm in inches
				WithMarginBottom(0.79).    // 20mm in inches
				WithMarginLeft(0.79).      // 20mm in inches
				WithMarginRight(0.79).     // 20mm in inches
				WithPrintBackground(true). // Print background colors/images
				WithScale(1.0).
				WithPreferCSSPageSize(false).
				Do(ctx)
			return err
		}),
	); err != nil {
		return nil, fmt.Errorf("PDF generation failed: %w", err)
	}

	return pdfBuf, nil
}

func sanitizeFilename(name string) string {
	replacer := strings.NewReplacer(
		"/", "_", "\\", "_", ":", "_", "*", "_",
		"?", "_", "\"", "_", "<", "_", ">", "_",
		"|", "_", " ", "_",
	)
	return replacer.Replace(name)
}
