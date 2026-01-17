package com.docflow;

import com.docflow.models.*;
import com.docflow.storage.Storage;
import com.docflow.storage.LocalStorage;
import com.openhtmltopdf.pdfboxout.PdfRendererBuilder;

import java.io.*;
import java.nio.file.*;
import java.time.Instant;
import java.util.*;
import java.util.regex.Pattern;

/**
 * Converts Markdown files to PDF.
 *
 * <p>
 * Example usage:
 * 
 * <pre>{@code
 * Storage storage = new LocalStorage("./output");
 * Converter converter = new Converter(storage);
 * 
 * List<MDFile> files = List.of(
 *         new MDFile("doc.md", "# Hello World"));
 * 
 * PDFResult result = converter.convertToPdf(files, ConvertOptions.separate());
 * }</pre>
 */
public class Converter {

    private final Storage storage;
    private final Path tempDir;
    private final MarkdownParser parser;
    private final Template template;

    /**
     * Create a converter with local storage.
     *
     * @param outputPath Path for output files.
     * @throws IOException If storage creation fails.
     */
    public Converter(String outputPath) throws IOException {
        this(new LocalStorage(outputPath));
    }

    /**
     * Create a converter with custom storage.
     *
     * @param storage Storage backend.
     */
    public Converter(Storage storage) {
        this.storage = storage;
        this.tempDir = Path.of(System.getProperty("java.io.tmpdir"), "docflow");
        this.parser = new MarkdownParser();
        this.template = new Template();

        try {
            Files.createDirectories(tempDir);
        } catch (IOException e) {
            throw new RuntimeException("Failed to create temp directory", e);
        }
    }

    /**
     * Convert markdown files to PDF.
     *
     * @param files   List of MDFile objects.
     * @param options Conversion options.
     * @return PDFResult with success status and file paths.
     */
    public PDFResult convertToPdf(List<MDFile> files, ConvertOptions options) {
        if (files == null || files.isEmpty()) {
            PDFResult result = new PDFResult(false, "At least one file is required");
            return result;
        }

        if (options == null) {
            options = ConvertOptions.separate();
        }

        long timestamp = Instant.now().getEpochSecond();

        try {
            if ("merged".equals(options.getMergeMode())) {
                String outputName = options.getOutputName();
                if (outputName == null || outputName.isEmpty()) {
                    outputName = "merged_" + timestamp;
                }

                byte[] pdfData = convertMerged(files, outputName);
                String path = outputName + ".pdf";
                storage.save(path, pdfData);

                PDFResult result = new PDFResult();
                result.setSuccess(true);
                result.setFilePaths(List.of(storage.getUrl(path).orElse(path)));
                return result;
            } else {
                List<String> paths = new ArrayList<>();
                for (int i = 0; i < files.size(); i++) {
                    MDFile file = files.get(i);
                    String baseName = file.getName().replaceAll("\\.[^.]+$", "");
                    String outputName = sanitizeFilename(baseName) + "_" + timestamp;

                    byte[] pdfData = convertSingle(file);
                    String path = outputName + ".pdf";
                    storage.save(path, pdfData);
                    paths.add(storage.getUrl(path).orElse(path));
                }

                PDFResult result = new PDFResult();
                result.setSuccess(true);
                result.setFilePaths(paths);
                return result;
            }
        } catch (Exception e) {
            PDFResult result = new PDFResult(false, e.getMessage());
            return result;
        }
    }

    /**
     * Convert markdown to PDF and return bytes.
     *
     * @param files List of MDFile objects.
     * @return PDF bytes.
     * @throws IOException If conversion fails.
     */
    public byte[] convertToBytes(List<MDFile> files) throws IOException {
        return convertMerged(files, "output");
    }

    /**
     * Generate HTML preview of markdown content.
     *
     * @param content Markdown content.
     * @return HTML string.
     */
    public String preview(String content) {
        return parser.toHtml(content);
    }

    private byte[] convertMerged(List<MDFile> files, String outputName) throws IOException {
        String mergedContent = parser.mergeFiles(files);
        return generatePdf(mergedContent);
    }

    private byte[] convertSingle(MDFile file) throws IOException {
        return generatePdf(file.getContent());
    }

    private byte[] generatePdf(String mdContent) throws IOException {
        // Convert to HTML
        String htmlContent = parser.toHtml(mdContent);
        String fullHtml = template.generate(htmlContent);

        // Generate PDF using OpenHTMLToPDF
        ByteArrayOutputStream baos = new ByteArrayOutputStream();
        PdfRendererBuilder builder = new PdfRendererBuilder();
        builder.useFastMode();
        builder.withHtmlContent(fullHtml, null);
        builder.toStream(baos);
        builder.run();

        return baos.toByteArray();
    }

    private String sanitizeFilename(String name) {
        return Pattern.compile("[/\\\\:*?\"<>| ]").matcher(name).replaceAll("_");
    }
}
