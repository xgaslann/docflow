package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/gorkem/md-to-pdf/internal/config"
	"github.com/gorkem/md-to-pdf/internal/model"
	"github.com/gorkem/md-to-pdf/pkg/pdf"
	"go.uber.org/zap"
)

// ConverterService handles PDF conversion
type ConverterService struct {
	cfg      *config.Config
	markdown *MarkdownService
	template *pdf.TemplateGenerator
	logger   *zap.Logger
}

// NewConverterService creates a new converter service instance
func NewConverterService(cfg *config.Config, markdown *MarkdownService, logger *zap.Logger) *ConverterService {
	return &ConverterService{
		cfg:      cfg,
		markdown: markdown,
		template: pdf.NewTemplateGenerator(),
		logger:   logger,
	}
}

// ConvertResult represents the result of a conversion operation
type ConvertResult struct {
	FilePath string
	FileName string
	Error    error
}

// Convert handles the conversion of files based on merge mode
func (s *ConverterService) Convert(ctx context.Context, req *model.ConvertRequest) (*model.ConvertResponse, error) {
	timestamp := time.Now().Unix()
	var results []ConvertResult

	// Sort files by order before processing
	sortedFiles := make([]model.FileData, len(req.Files))
	copy(sortedFiles, req.Files)
	sort.Slice(sortedFiles, func(i, j int) bool {
		return sortedFiles[i].Order < sortedFiles[j].Order
	})

	switch req.MergeMode {
	case model.MergeModeMerged:
		result := s.convertMerged(ctx, sortedFiles, req.OutputName, timestamp)
		results = append(results, result)
	case model.MergeModeSeparate:
		results = s.convertSeparate(ctx, sortedFiles, timestamp)
	default:
		return nil, fmt.Errorf("invalid merge mode: %s", req.MergeMode)
	}

	// Check for errors and collect successful files
	var outputFiles []string
	var errors []string

	for _, result := range results {
		if result.Error != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", result.FileName, result.Error))
			s.logger.Error("conversion failed",
				zap.String("file", result.FileName),
				zap.Error(result.Error),
			)
		} else {
			outputFiles = append(outputFiles, result.FilePath)
			s.logger.Info("conversion successful",
				zap.String("file", result.FileName),
				zap.String("path", result.FilePath),
			)
		}
	}

	if len(errors) > 0 && len(outputFiles) == 0 {
		return &model.ConvertResponse{
			Success: false,
			Error:   strings.Join(errors, "; "),
		}, nil
	}

	return &model.ConvertResponse{
		Success: true,
		Files:   outputFiles,
	}, nil
}

func (s *ConverterService) convertMerged(ctx context.Context, files []model.FileData, outputName string, timestamp int64) ConvertResult {
	mergedContent := s.markdown.MergeFiles(files)

	if outputName == "" {
		outputName = fmt.Sprintf("merged_%d", timestamp)
	}
	outputName = sanitizeFilename(outputName)

	pdfPath, err := s.generatePDF(ctx, mergedContent, outputName)
	return ConvertResult{
		FilePath: pdfPath,
		FileName: outputName,
		Error:    err,
	}
}

func (s *ConverterService) convertSeparate(ctx context.Context, files []model.FileData, timestamp int64) []ConvertResult {
	results := make([]ConvertResult, len(files))

	for i, file := range files {
		baseName := strings.TrimSuffix(file.Name, filepath.Ext(file.Name))
		outputName := fmt.Sprintf("%s_%d", sanitizeFilename(baseName), timestamp)

		pdfPath, err := s.generatePDF(ctx, file.Content, outputName)
		results[i] = ConvertResult{
			FilePath: pdfPath,
			FileName: file.Name,
			Error:    err,
		}
	}

	return results
}

func (s *ConverterService) generatePDF(ctx context.Context, mdContent, outputName string) (string, error) {
	htmlContent, err := s.markdown.ToHTML(mdContent)
	if err != nil {
		return "", fmt.Errorf("markdown conversion failed: %w", err)
	}

	fullHTML := s.template.Generate(htmlContent)

	// Write HTML to temp file
	tempHTMLPath := filepath.Join(s.cfg.Storage.TempDir, outputName+".html")
	absHTMLPath, err := filepath.Abs(tempHTMLPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	if err := os.WriteFile(absHTMLPath, []byte(fullHTML), 0644); err != nil {
		return "", fmt.Errorf("failed to write temp file: %w", err)
	}
	defer os.Remove(absHTMLPath)

	outputPath := filepath.Join(s.cfg.Storage.OutputDir, outputName+".pdf")
	absOutputPath, err := filepath.Abs(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to get output path: %w", err)
	}

	// Use chromedp for PDF generation
	if err := s.generateWithChromedp(ctx, absHTMLPath, absOutputPath); err != nil {
		s.logger.Error("chromedp PDF generation failed", zap.Error(err))
		return "", fmt.Errorf("PDF generation failed: %w", err)
	}

	return "/output/" + outputName + ".pdf", nil
}

func (s *ConverterService) generateWithChromedp(ctx context.Context, htmlPath, outputPath string) error {
	// Create chromedp context with options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-software-rasterizer", true),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	defer allocCancel()

	taskCtx, taskCancel := chromedp.NewContext(allocCtx)
	defer taskCancel()

	// Set timeout
	taskCtx, cancel := context.WithTimeout(taskCtx, 60*time.Second)
	defer cancel()

	// Navigate to the HTML file and generate PDF
	var pdfBuf []byte

	fileURL := "file://" + htmlPath

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
				WithPreferCSSPageSize(false). // Use our paper size, not CSS
				Do(ctx)
			return err
		}),
	); err != nil {
		return fmt.Errorf("chromedp execution failed: %w", err)
	}

	// Write PDF to file
	if err := os.WriteFile(outputPath, pdfBuf, 0644); err != nil {
		return fmt.Errorf("failed to write PDF file: %w", err)
	}

	return nil
}

func sanitizeFilename(name string) string {
	replacer := strings.NewReplacer(
		"/", "_", "\\", "_", ":", "_", "*", "_",
		"?", "_", "\"", "_", "<", "_", ">", "_",
		"|", "_", " ", "_",
	)
	return replacer.Replace(name)
}
