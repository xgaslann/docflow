package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gorkem/md-to-pdf/internal/config"
	"github.com/gorkem/md-to-pdf/internal/model"
	"github.com/gorkem/md-to-pdf/internal/service"
	"go.uber.org/zap"
)

func setupTestHandler() (*Handler, *fiber.App) {
	logger, _ := zap.NewDevelopment()
	cfg := &config.Config{
		Storage: config.StorageConfig{
			TempDir:   "/tmp/test-temp",
			OutputDir: "/tmp/test-output",
		},
	}

	markdownSvc := service.NewMarkdownService()
	converterSvc := service.NewConverterService(cfg, markdownSvc, logger)
	pdfExtractorSvc := service.NewPDFExtractorService(cfg, logger)

	h := NewHandler(markdownSvc, converterSvc, pdfExtractorSvc, logger)

	app := fiber.New()
	h.RegisterRoutes(app)

	return h, app
}

func TestHealthCheck(t *testing.T) {
	_, app := setupTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var result model.HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if result.Status != "healthy" {
		t.Errorf("expected status 'healthy', got %q", result.Status)
	}

	if result.Version == "" {
		t.Error("expected non-empty version")
	}

	if result.Timestamp == 0 {
		t.Error("expected non-zero timestamp")
	}
}

func TestPreview(t *testing.T) {
	_, app := setupTestHandler()

	tests := []struct {
		name           string
		body           map[string]interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "valid markdown",
			body: map[string]interface{}{
				"content": "# Hello World",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var result model.PreviewResponse
				if err := json.Unmarshal(body, &result); err != nil {
					t.Fatalf("failed to decode: %v", err)
				}
				if result.HTML == "" {
					t.Error("expected non-empty HTML")
				}
				if !bytes.Contains(body, []byte("Hello World")) {
					t.Error("expected HTML to contain 'Hello World'")
				}
			},
		},
		{
			name: "empty content",
			body: map[string]interface{}{
				"content": "",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing content field",
			body:           map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "complex markdown",
			body: map[string]interface{}{
				"content": "# Title\n\n- Item 1\n- Item 2\n\n```go\nfunc main() {}\n```",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				if !bytes.Contains(body, []byte("<h1>")) {
					t.Error("expected h1 tag")
				}
				if !bytes.Contains(body, []byte("<li>")) {
					t.Error("expected li tag")
				}
				if !bytes.Contains(body, []byte("<code>")) {
					t.Error("expected code tag")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/preview", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("failed to make request: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			if tt.checkResponse != nil {
				buf := new(bytes.Buffer)
				buf.ReadFrom(resp.Body)
				tt.checkResponse(t, buf.Bytes())
			}
		})
	}
}

func TestMergePreview(t *testing.T) {
	_, app := setupTestHandler()

	tests := []struct {
		name           string
		body           map[string]interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "multiple files",
			body: map[string]interface{}{
				"files": []map[string]interface{}{
					{"id": "1", "name": "a.md", "content": "# First", "order": 0},
					{"id": "2", "name": "b.md", "content": "# Second", "order": 1},
				},
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var result model.MergePreviewResponse
				if err := json.Unmarshal(body, &result); err != nil {
					t.Fatalf("failed to decode: %v", err)
				}
				if result.TotalFiles != 2 {
					t.Errorf("expected 2 files, got %d", result.TotalFiles)
				}
				if result.EstimatedPages < 1 {
					t.Error("expected at least 1 page")
				}
				if result.HTML == "" {
					t.Error("expected non-empty HTML")
				}
			},
		},
		{
			name: "single file",
			body: map[string]interface{}{
				"files": []map[string]interface{}{
					{"id": "1", "name": "only.md", "content": "# Only", "order": 0},
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "empty files array",
			body: map[string]interface{}{
				"files": []map[string]interface{}{},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing files field",
			body:           map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/preview/merge", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("failed to make request: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			if tt.checkResponse != nil {
				buf := new(bytes.Buffer)
				buf.ReadFrom(resp.Body)
				tt.checkResponse(t, buf.Bytes())
			}
		})
	}
}

func TestConvert_Validation(t *testing.T) {
	_, app := setupTestHandler()

	tests := []struct {
		name           string
		body           map[string]interface{}
		expectedStatus int
	}{
		{
			name: "empty files",
			body: map[string]interface{}{
				"files":     []map[string]interface{}{},
				"mergeMode": "separate",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing files",
			body:           map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "valid request structure",
			body: map[string]interface{}{
				"files": []map[string]interface{}{
					{"id": "1", "name": "test.md", "content": "# Test", "order": 0},
				},
				"mergeMode": "separate",
			},
			// This will fail because chromedp isn't available in test
			// but it validates the request structure is correct
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/convert", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1) // -1 for no timeout
			if err != nil {
				t.Fatalf("failed to make request: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				buf := new(bytes.Buffer)
				buf.ReadFrom(resp.Body)
				t.Errorf("expected status %d, got %d. Body: %s", tt.expectedStatus, resp.StatusCode, buf.String())
			}
		})
	}
}

func TestPDFPreview_Validation(t *testing.T) {
	_, app := setupTestHandler()

	tests := []struct {
		name           string
		body           map[string]interface{}
		expectedStatus int
	}{
		{
			name: "missing content",
			body: map[string]interface{}{
				"fileName": "test.pdf",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing filename",
			body: map[string]interface{}{
				"content": "dGVzdA==", // base64 "test"
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty body",
			body:           map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/pdf/preview", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("failed to make request: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func TestPDFExtract_Validation(t *testing.T) {
	_, app := setupTestHandler()

	tests := []struct {
		name           string
		body           map[string]interface{}
		expectedStatus int
	}{
		{
			name: "missing content",
			body: map[string]interface{}{
				"fileName": "test.pdf",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing filename",
			body: map[string]interface{}{
				"content": "dGVzdA==",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty body",
			body:           map[string]interface{}{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/pdf/extract", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("failed to make request: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func TestInvalidJSON(t *testing.T) {
	_, app := setupTestHandler()

	endpoints := []string{
		"/api/preview",
		"/api/preview/merge",
		"/api/convert",
		"/api/pdf/preview",
		"/api/pdf/extract",
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, endpoint, bytes.NewReader([]byte("not json")))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("failed to make request: %v", err)
			}

			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("expected status 400 for invalid JSON, got %d", resp.StatusCode)
			}
		})
	}
}
