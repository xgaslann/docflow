package handler

import (
	"encoding/base64"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gorkem/md-to-pdf/internal/model"
	"github.com/gorkem/md-to-pdf/internal/service"
	"go.uber.org/zap"
)

const version = "1.0.0"

// Handler contains all HTTP handlers
type Handler struct {
	markdown     *service.MarkdownService
	converter    *service.ConverterService
	pdfExtractor *service.PDFExtractorService
	logger       *zap.Logger
}

// NewHandler creates a new handler instance
func NewHandler(
	markdown *service.MarkdownService,
	converter *service.ConverterService,
	pdfExtractor *service.PDFExtractorService,
	logger *zap.Logger,
) *Handler {
	return &Handler{
		markdown:     markdown,
		converter:    converter,
		pdfExtractor: pdfExtractor,
		logger:       logger,
	}
}

// RegisterRoutes registers all routes
func (h *Handler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Health
	api.Get("/health", h.HealthCheck)

	// MD to PDF
	api.Post("/preview", h.Preview)
	api.Post("/preview/merge", h.MergePreview)
	api.Post("/convert", h.Convert)

	// PDF to MD
	api.Post("/pdf/preview", h.PDFPreview)
	api.Post("/pdf/extract", h.PDFExtract)
}

// HealthCheck handles health check requests
func (h *Handler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(model.HealthResponse{
		Status:    "healthy",
		Version:   version,
		Timestamp: time.Now().Unix(),
	})
}

// Preview handles markdown preview requests
func (h *Handler) Preview(c *fiber.Ctx) error {
	var req model.PreviewRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn("invalid preview request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Success: false,
			Error:   "Invalid request body",
			Code:    "INVALID_REQUEST",
		})
	}

	if req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Success: false,
			Error:   "Content is required",
			Code:    "CONTENT_REQUIRED",
		})
	}

	html, err := h.markdown.ToHTML(req.Content)
	if err != nil {
		h.logger.Error("markdown conversion failed", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Success: false,
			Error:   "Failed to convert markdown",
			Code:    "CONVERSION_ERROR",
		})
	}

	return c.JSON(model.PreviewResponse{HTML: html})
}

// MergePreview handles merge preview requests - shows how merged document will look
func (h *Handler) MergePreview(c *fiber.Ctx) error {
	var req model.MergePreviewRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn("invalid merge preview request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Success: false,
			Error:   "Invalid request body",
			Code:    "INVALID_REQUEST",
		})
	}

	if len(req.Files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Success: false,
			Error:   "At least one file is required",
			Code:    "FILES_REQUIRED",
		})
	}

	html, err := h.markdown.MergeFilesToHTML(req.Files)
	if err != nil {
		h.logger.Error("merge preview failed", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Success: false,
			Error:   "Failed to generate merge preview",
			Code:    "PREVIEW_ERROR",
		})
	}

	mergedContent := h.markdown.MergeFiles(req.Files)
	estimatedPages := h.markdown.EstimatePageCount(mergedContent)

	return c.JSON(model.MergePreviewResponse{
		HTML:           html,
		TotalFiles:     len(req.Files),
		EstimatedPages: estimatedPages,
	})
}

// Convert handles PDF conversion requests
func (h *Handler) Convert(c *fiber.Ctx) error {
	var req model.ConvertRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn("invalid convert request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Success: false,
			Error:   "Invalid request body",
			Code:    "INVALID_REQUEST",
		})
	}

	if len(req.Files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Success: false,
			Error:   "At least one file is required",
			Code:    "FILES_REQUIRED",
		})
	}

	if req.MergeMode == "" {
		req.MergeMode = model.MergeModeSeparate
	}

	result, err := h.converter.Convert(c.Context(), &req)
	if err != nil {
		h.logger.Error("conversion failed", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    "CONVERSION_ERROR",
		})
	}

	return c.JSON(result)
}

// PDFPreview handles PDF preview requests
func (h *Handler) PDFPreview(c *fiber.Ctx) error {
	var req model.PDFExtractRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn("invalid PDF preview request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Success: false,
			Error:   "Invalid request body",
			Code:    "INVALID_REQUEST",
		})
	}

	if req.Content == "" || req.FileName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Success: false,
			Error:   "Content and fileName are required",
			Code:    "MISSING_FIELDS",
		})
	}

	// Decode base64 PDF
	pdfData, err := base64.StdEncoding.DecodeString(req.Content)
	if err != nil {
		h.logger.Warn("invalid base64 content", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Success: false,
			Error:   "Invalid PDF content",
			Code:    "INVALID_CONTENT",
		})
	}

	result, err := h.pdfExtractor.PreviewExtraction(c.Context(), pdfData, req.FileName)
	if err != nil {
		h.logger.Error("PDF preview failed", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    "PREVIEW_ERROR",
		})
	}

	return c.JSON(result)
}

// PDFExtract handles PDF to Markdown extraction requests
func (h *Handler) PDFExtract(c *fiber.Ctx) error {
	var req model.PDFExtractRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Warn("invalid PDF extract request", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Success: false,
			Error:   "Invalid request body",
			Code:    "INVALID_REQUEST",
		})
	}

	if req.Content == "" || req.FileName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Success: false,
			Error:   "Content and fileName are required",
			Code:    "MISSING_FIELDS",
		})
	}

	// Decode base64 PDF
	pdfData, err := base64.StdEncoding.DecodeString(req.Content)
	if err != nil {
		h.logger.Warn("invalid base64 content", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Success: false,
			Error:   "Invalid PDF content",
			Code:    "INVALID_CONTENT",
		})
	}

	result, err := h.pdfExtractor.ExtractToMarkdown(c.Context(), pdfData, req.FileName)
	if err != nil {
		h.logger.Error("PDF extraction failed", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    "EXTRACTION_ERROR",
		})
	}

	return c.JSON(result)
}
