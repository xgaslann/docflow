package model

import (
	"encoding/json"
	"testing"
)

func TestFileDataJSON(t *testing.T) {
	original := FileData{
		ID:      "test-123",
		Name:    "document.md",
		Content: "# Hello\n\nWorld",
		Order:   5,
	}

	// Marshal
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Unmarshal
	var decoded FileData
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Compare
	if decoded.ID != original.ID {
		t.Errorf("ID mismatch: got %q, want %q", decoded.ID, original.ID)
	}
	if decoded.Name != original.Name {
		t.Errorf("Name mismatch: got %q, want %q", decoded.Name, original.Name)
	}
	if decoded.Content != original.Content {
		t.Errorf("Content mismatch: got %q, want %q", decoded.Content, original.Content)
	}
	if decoded.Order != original.Order {
		t.Errorf("Order mismatch: got %d, want %d", decoded.Order, original.Order)
	}
}

func TestConvertRequestJSON(t *testing.T) {
	jsonStr := `{
		"files": [
			{"id": "1", "name": "a.md", "content": "# A", "order": 0},
			{"id": "2", "name": "b.md", "content": "# B", "order": 1}
		],
		"mergeMode": "merged",
		"outputName": "combined"
	}`

	var req ConvertRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(req.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(req.Files))
	}
	if req.MergeMode != MergeModeMerged {
		t.Errorf("expected merge mode %q, got %q", MergeModeMerged, req.MergeMode)
	}
	if req.OutputName != "combined" {
		t.Errorf("expected output name 'combined', got %q", req.OutputName)
	}
}

func TestMergeModeConstants(t *testing.T) {
	if MergeModeSeparate != "separate" {
		t.Errorf("MergeModeSeparate should be 'separate'")
	}
	if MergeModeMerged != "merged" {
		t.Errorf("MergeModeMerged should be 'merged'")
	}
}

func TestConvertResponseJSON(t *testing.T) {
	resp := ConvertResponse{
		Success: true,
		Files:   []string{"/output/doc1.pdf", "/output/doc2.pdf"},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Check JSON structure
	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded["success"] != true {
		t.Error("expected success to be true")
	}

	files, ok := decoded["files"].([]interface{})
	if !ok {
		t.Fatal("expected files to be an array")
	}
	if len(files) != 2 {
		t.Errorf("expected 2 files, got %d", len(files))
	}
}

func TestErrorResponseJSON(t *testing.T) {
	resp := ErrorResponse{
		Success: false,
		Error:   "something went wrong",
		Code:    "INTERNAL_ERROR",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded ErrorResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Success != false {
		t.Error("expected success to be false")
	}
	if decoded.Error != resp.Error {
		t.Errorf("error mismatch: got %q, want %q", decoded.Error, resp.Error)
	}
	if decoded.Code != resp.Code {
		t.Errorf("code mismatch: got %q, want %q", decoded.Code, resp.Code)
	}
}

func TestPDFExtractRequestJSON(t *testing.T) {
	jsonStr := `{
		"fileName": "document.pdf",
		"content": "base64encodedcontent"
	}`

	var req PDFExtractRequest
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if req.FileName != "document.pdf" {
		t.Errorf("expected fileName 'document.pdf', got %q", req.FileName)
	}
	if req.Content != "base64encodedcontent" {
		t.Errorf("expected content 'base64encodedcontent', got %q", req.Content)
	}
}

func TestPDFExtractResponseJSON(t *testing.T) {
	resp := PDFExtractResponse{
		Success:  true,
		Markdown: "# Extracted Content",
		FilePath: "/output/doc.md",
		FileName: "doc.md",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded PDFExtractResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Success != true {
		t.Error("expected success to be true")
	}
	if decoded.Markdown != resp.Markdown {
		t.Errorf("markdown mismatch")
	}
	if decoded.FilePath != resp.FilePath {
		t.Errorf("filePath mismatch")
	}
}

func TestHealthResponseJSON(t *testing.T) {
	resp := HealthResponse{
		Status:    "healthy",
		Version:   "1.0.0",
		Timestamp: 1234567890,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded HealthResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.Status != "healthy" {
		t.Errorf("expected status 'healthy', got %q", decoded.Status)
	}
	if decoded.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got %q", decoded.Version)
	}
	if decoded.Timestamp != 1234567890 {
		t.Errorf("expected timestamp 1234567890, got %d", decoded.Timestamp)
	}
}
