package com.docflow;

import com.docflow.models.MDResult;
import com.docflow.storage.Storage;
import org.apache.pdfbox.pdmodel.PDDocument;
import org.apache.pdfbox.text.PDFTextStripper;

import java.io.*;
import java.nio.file.*;
import java.time.Instant;
import java.util.regex.Pattern;

/**
 * Extracts text from PDF files and converts to Markdown.
 *
 * <p>
 * Uses Apache PDFBox for text extraction.
 *
 * <p>
 * Example usage:
 * 
 * <pre>{@code
 * Extractor extractor = new Extractor();
 * byte[] pdfData = Files.readAllBytes(Path.of("document.pdf"));
 * MDResult result = extractor.extractToMarkdown(pdfData, "document.pdf");
 * System.out.println(result.getMarkdown());
 * }</pre>
 */
public class Extractor {

    private final Storage storage;
    private final Path tempDir;

    /**
     * Create an extractor without storage.
     */
    public Extractor() {
        this(null);
    }

    /**
     * Create an extractor with storage.
     *
     * @param storage Storage backend for saving output files.
     */
    public Extractor(Storage storage) {
        this.storage = storage;
        this.tempDir = Path.of(System.getProperty("java.io.tmpdir"), "docflow");

        try {
            Files.createDirectories(tempDir);
        } catch (IOException e) {
            throw new RuntimeException("Failed to create temp directory", e);
        }
    }

    /**
     * Extract text from PDF and convert to Markdown.
     *
     * @param pdfData  PDF file bytes.
     * @param filename Original filename.
     * @return MDResult with markdown content.
     */
    public MDResult extractToMarkdown(byte[] pdfData, String filename) {
        if (pdfData == null || pdfData.length == 0) {
            MDResult result = new MDResult();
            result.setSuccess(false);
            result.setError("PDF data is required");
            return result;
        }

        long timestamp = Instant.now().getEpochSecond();
        String baseName = filename.replaceAll("\\.[^.]+$", "");
        String safeName = sanitizeFilename(baseName);

        try {
            // Extract text
            String text = extractText(pdfData);
            int pageCount = getPageCount(pdfData);

            // Convert to markdown
            String markdown = textToMarkdown(text, baseName);

            // Save if storage configured
            String outputPath = null;
            if (storage != null) {
                outputPath = safeName + "_" + timestamp + ".md";
                storage.save(outputPath, markdown.getBytes());
                outputPath = storage.getUrl(outputPath).orElse(outputPath);
            }

            MDResult result = new MDResult();
            result.setSuccess(true);
            result.setMarkdown(markdown);
            result.setFilePath(outputPath);
            result.setFileName(safeName + ".md");
            result.setPageCount(pageCount);
            return result;

        } catch (Exception e) {
            MDResult result = new MDResult();
            result.setSuccess(false);
            result.setError(e.getMessage());
            return result;
        }
    }

    /**
     * Extract markdown from a PDF file path.
     *
     * @param path Path to PDF file.
     * @return MDResult with markdown content.
     * @throws IOException If reading fails.
     */
    public MDResult extractFromFile(String path) throws IOException {
        byte[] data = Files.readAllBytes(Path.of(path));
        return extractToMarkdown(data, Path.of(path).getFileName().toString());
    }

    /**
     * Get page count from PDF.
     *
     * @param pdfData PDF bytes.
     * @return Number of pages.
     * @throws IOException If parsing fails.
     */
    public int getPageCount(byte[] pdfData) throws IOException {
        try (PDDocument doc = PDDocument.load(pdfData)) {
            return doc.getNumberOfPages();
        }
    }

    /**
     * Get preview of extracted content.
     *
     * @param pdfData  PDF bytes.
     * @param filename Original filename.
     * @return MDResult with truncated markdown.
     */
    public MDResult preview(byte[] pdfData, String filename) {
        MDResult result = extractToMarkdown(pdfData, filename);
        if (result.isSuccess() && result.getMarkdown().length() > 2000) {
            result.setMarkdown(result.getMarkdown().substring(0, 2000) + "\n\n... (continued)");
        }
        return result;
    }

    private String extractText(byte[] pdfData) throws IOException {
        try (PDDocument doc = PDDocument.load(pdfData)) {
            PDFTextStripper stripper = new PDFTextStripper();
            stripper.setSortByPosition(true);
            return stripper.getText(doc);
        }
    }

    private String textToMarkdown(String text, String title) {
        StringBuilder result = new StringBuilder();
        result.append("# ").append(title).append("\n\n");

        String[] lines = text.split("\n");
        boolean inParagraph = false;

        for (String line : lines) {
            line = line.trim();

            if (line.isEmpty()) {
                if (inParagraph) {
                    result.append("\n\n");
                    inParagraph = false;
                }
                continue;
            }

            // Detect headers (ALL CAPS)
            if (isPotentialHeader(line)) {
                if (inParagraph) {
                    result.append("\n\n");
                }
                result.append("## ").append(toTitleCase(line)).append("\n");
                inParagraph = false;
                continue;
            }

            // Detect bullet points
            if (line.startsWith("•") || line.startsWith("●") || line.startsWith("○")) {
                if (inParagraph) {
                    result.append("\n\n");
                }
                result.append("- ").append(line.substring(1).trim()).append("\n");
                inParagraph = false;
                continue;
            }

            // Regular paragraph
            if (inParagraph) {
                result.append(" ");
            }
            result.append(line);
            inParagraph = true;
        }

        result.append("\n");
        return result.toString();
    }

    private boolean isPotentialHeader(String line) {
        return line.length() > 3 && line.length() < 60
                && line.equals(line.toUpperCase())
                && line.split("\\s+").length >= 2;
    }

    private String toTitleCase(String text) {
        StringBuilder result = new StringBuilder();
        boolean nextTitleCase = true;

        for (char c : text.toLowerCase().toCharArray()) {
            if (Character.isSpaceChar(c)) {
                nextTitleCase = true;
            } else if (nextTitleCase) {
                c = Character.toTitleCase(c);
                nextTitleCase = false;
            }
            result.append(c);
        }

        return result.toString();
    }

    private String sanitizeFilename(String name) {
        return Pattern.compile("[/\\\\:*?\"<>| ]").matcher(name).replaceAll("_");
    }
}
