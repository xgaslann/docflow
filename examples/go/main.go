package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/xgaslan/docflow/sdks/go/docflow"
	"github.com/xgaslan/docflow/sdks/go/docflow/storage"
)

func main() {
	// Create output directory
	outputDir := "./output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatal(err)
	}

	// Example 1: Basic MD to PDF conversion
	fmt.Println("=== Example 1: Basic Conversion ===")
	basicConversion(outputDir)

	// Example 2: Merge multiple files
	fmt.Println("\n=== Example 2: Merge Multiple Files ===")
	mergeFiles(outputDir)

	// Example 3: PDF to Markdown extraction
	fmt.Println("\n=== Example 3: PDF Extraction ===")
	extractPDF()

	// Example 4: Get PDF as bytes
	fmt.Println("\n=== Example 4: Get PDF as Bytes ===")
	getPDFBytes(outputDir)

	// Example 5: Preview markdown
	fmt.Println("\n=== Example 5: Preview Markdown ===")
	previewMarkdown()
}

func basicConversion(outputDir string) {
	// Create storage
	store, err := storage.NewLocalStorage(outputDir)
	if err != nil {
		log.Fatal(err)
	}

	// Create converter
	converter := docflow.NewConverter(
		docflow.WithStorage(store),
	)

	// Create markdown file
	files := []docflow.MDFile{
		docflow.NewMDFile("hello.md", `# Hello World

This is a **bold** statement and this is *italic*.

## Features

- Easy to use
- Standalone library
- No server required

## Code Example

`+"`"+`go
fmt.Println("Hello, World!")
`+"`"+`
`),
	}

	// Convert to PDF
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := converter.ConvertToPDF(ctx, files, docflow.ConvertOptions{
		MergeMode: "separate",
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	if result.Success {
		fmt.Printf("✓ PDF created: %v\n", result.FilePaths)
	} else {
		fmt.Printf("✗ Error: %v\n", result.Error)
	}
}

func mergeFiles(outputDir string) {
	store, _ := storage.NewLocalStorage(outputDir)
	converter := docflow.NewConverter(docflow.WithStorage(store))

	files := []docflow.MDFile{
		docflow.NewMDFileWithOrder("chapter1.md", "# Chapter 1\n\nIntroduction to the topic.", 0),
		docflow.NewMDFileWithOrder("chapter2.md", "# Chapter 2\n\nDeeper exploration.", 1),
		docflow.NewMDFileWithOrder("chapter3.md", "# Chapter 3\n\nConclusion.", 2),
	}

	ctx := context.Background()
	result, err := converter.ConvertToPDF(ctx, files, docflow.ConvertOptions{
		MergeMode:  "merged",
		OutputName: "combined_document",
	})

	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	if result.Success {
		fmt.Printf("✓ Merged PDF created: %v\n", result.FilePaths)
	}
}

func extractPDF() {
	// Check if sample PDF exists
	samplePDF := "./sample.pdf"
	if _, err := os.Stat(samplePDF); os.IsNotExist(err) {
		fmt.Println("⊙ Skipping: sample.pdf not found")
		fmt.Println("  Create a sample.pdf file to test extraction")
		return
	}

	extractor := docflow.NewExtractor()

	ctx := context.Background()
	result, err := extractor.ExtractFromFile(ctx, samplePDF)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	if result.Success {
		fmt.Printf("✓ Extracted %d pages\n", result.PageCount)
		// Show first 200 chars of markdown
		preview := result.Markdown
		if len(preview) > 200 {
			preview = preview[:200] + "..."
		}
		fmt.Printf("  Preview: %s\n", preview)
	}
}

func getPDFBytes(outputDir string) {
	converter := docflow.NewConverter()

	files := []docflow.MDFile{
		docflow.NewMDFile("inline.md", "# Inline PDF\n\nThis PDF is generated as bytes."),
	}

	ctx := context.Background()
	pdfBytes, err := converter.ConvertToBytes(ctx, files)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	// Save bytes to file
	outputPath := filepath.Join(outputDir, "from_bytes.pdf")
	if err := os.WriteFile(outputPath, pdfBytes, 0644); err != nil {
		log.Printf("Error saving: %v\n", err)
		return
	}

	fmt.Printf("✓ PDF bytes saved: %s (%d bytes)\n", outputPath, len(pdfBytes))
}

func previewMarkdown() {
	converter := docflow.NewConverter()

	html, err := converter.Preview("# Preview\n\nThis is a **preview** of the markdown.")
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("✓ HTML Preview (%d bytes):\n", len(html))
	// Show first 100 chars
	if len(html) > 100 {
		fmt.Printf("  %s...\n", html[:100])
	} else {
		fmt.Printf("  %s\n", html)
	}
}
