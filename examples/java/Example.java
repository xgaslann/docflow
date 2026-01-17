package com.docflow.examples;

import com.docflow.*;
import com.docflow.models.*;
import com.docflow.storage.*;

import java.io.IOException;
import java.nio.file.*;
import java.util.List;

/**
 * DocFlow Java SDK Examples
 *
 * Demonstrates standalone usage of the DocFlow library.
 * No server required.
 */
public class Example {

    public static void main(String[] args) throws IOException {
        // Create output directory
        Path outputDir = Path.of("./output");
        Files.createDirectories(outputDir);

        System.out.println("=== Example 1: Basic Conversion ===");
        basicConversion(outputDir);

        System.out.println("\n=== Example 2: Merge Multiple Files ===");
        mergeFiles(outputDir);

        System.out.println("\n=== Example 3: PDF Extraction ===");
        extractPdf();

        System.out.println("\n=== Example 4: Get PDF as Bytes ===");
        getPdfBytes(outputDir);

        System.out.println("\n=== Example 5: Preview Markdown ===");
        previewMarkdown();
    }

    static void basicConversion(Path outputDir) throws IOException {
        // Create storage
        Storage storage = new LocalStorage(outputDir.toString());
        Converter converter = new Converter(storage);

        // Create markdown file
        List<MDFile> files = List.of(
                new MDFile("hello.md", """
                        # Hello World

                        This is a **bold** statement and this is *italic*.

                        ## Features

                        - Easy to use
                        - Standalone library
                        - No server required

                        ## Code Example

                        ```java
                        System.out.println("Hello, World!");
                        ```
                        """));

        // Convert to PDF
        PDFResult result = converter.convertToPdf(files, ConvertOptions.separate());

        if (result.isSuccess()) {
            System.out.println("✓ PDF created: " + result.getFilePaths());
        } else {
            System.out.println("✗ Error: " + result.getError());
        }
    }

    static void mergeFiles(Path outputDir) throws IOException {
        Storage storage = new LocalStorage(outputDir.toString());
        Converter converter = new Converter(storage);

        List<MDFile> files = List.of(
                MDFile.builder().name("chapter1.md").content("# Chapter 1\n\nIntroduction.").order(0).build(),
                MDFile.builder().name("chapter2.md").content("# Chapter 2\n\nDeeper exploration.").order(1).build(),
                MDFile.builder().name("chapter3.md").content("# Chapter 3\n\nConclusion.").order(2).build());

        PDFResult result = converter.convertToPdf(files, ConvertOptions.merged("combined_document"));

        if (result.isSuccess()) {
            System.out.println("✓ Merged PDF created: " + result.getFilePaths());
        }
    }

    static void extractPdf() {
        Path samplePdf = Path.of("./sample.pdf");
        if (!Files.exists(samplePdf)) {
            System.out.println("⊙ Skipping: sample.pdf not found");
            System.out.println("  Create a sample.pdf file to test extraction");
            return;
        }

        try {
            Extractor extractor = new Extractor();
            MDResult result = extractor.extractFromFile(samplePdf.toString());

            if (result.isSuccess()) {
                System.out.println("✓ Extraction completed");
                String preview = result.getMarkdown();
                if (preview.length() > 200) {
                    preview = preview.substring(0, 200) + "...";
                }
                System.out.println("  Preview: " + preview);
            }
        } catch (IOException e) {
            System.out.println("✗ Error: " + e.getMessage());
        }
    }

    static void getPdfBytes(Path outputDir) throws IOException {
        Converter converter = new Converter(outputDir.toString());

        List<MDFile> files = List.of(
                new MDFile("inline.md", "# Inline PDF\n\nThis PDF is generated as bytes."));

        byte[] pdfBytes = converter.convertToBytes(files);

        // Save bytes to file
        Path outputPath = outputDir.resolve("from_bytes.pdf");
        Files.write(outputPath, pdfBytes);

        System.out.println("✓ PDF bytes saved: " + outputPath + " (" + pdfBytes.length + " bytes)");
    }

    static void previewMarkdown() throws IOException {
        Converter converter = new Converter("./output");
        String html = converter.preview("# Preview\n\nThis is a **preview** of the markdown.");

        System.out.println("✓ HTML Preview (" + html.length() + " bytes):");
        String preview = html.length() > 100 ? html.substring(0, 100) + "..." : html;
        System.out.println("  " + preview);
    }
}
