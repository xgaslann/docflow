package com.docflow.rag;

import com.docflow.formats.*;
import com.docflow.models.*;

import java.nio.file.*;
import java.util.*;

/**
 * RAGProcessor - Main orchestrator for RAG document processing.
 * Handles document conversion, chunking, and optional LLM enrichment.
 */
public class RAGProcessor {

    private final RAGConfig config;
    private final RAGChunker chunker;
    private final LLMProcessor llmProcessor;

    public RAGProcessor(RAGConfig config) {
        this.config = config;
        this.chunker = new RAGChunker(config);

        if (config.getLlmConfig() != null) {
            this.llmProcessor = new LLMProcessor(config.getLlmConfig());
        } else {
            this.llmProcessor = null;
        }
    }

    /**
     * Process a file and return a RAGDocument.
     */
    public RAGDocument processFile(String filePath) throws Exception {
        byte[] fileBytes = Files.readAllBytes(Paths.get(filePath));
        String filename = Paths.get(filePath).getFileName().toString();
        return process(fileBytes, filename);
    }

    /**
     * Process bytes and return a RAGDocument.
     */
    public RAGDocument process(byte[] data, String filename) throws Exception {
        // Convert to markdown
        ConvertResult convertResult = convertToMarkdown(data, filename);

        // Create document
        RAGDocument doc = new RAGDocument();
        doc.setId(UUID.randomUUID().toString());
        doc.setFilename(filename);
        doc.setContent(convertResult.getContent());

        // Chunk the content
        List<Chunk> ragChunks = chunker.chunk(convertResult.getContent());

        // Convert rag.Chunk to models.Chunk
        List<com.docflow.models.Chunk> modelChunks = new ArrayList<>();
        for (Chunk ragChunk : ragChunks) {
            com.docflow.models.Chunk mc = new com.docflow.models.Chunk();
            mc.setIndex(ragChunk.getIndex());
            mc.setContent(ragChunk.getContent());
            modelChunks.add(mc);
        }
        doc.setChunks(modelChunks);

        // Extract images if available
        if (convertResult.getImages() != null) {
            doc.setImages(convertResult.getImages());

            // Describe images with LLM if configured
            if (llmProcessor != null && config.isDescribeImages()) {
                describeImages(doc);
            }
        }

        // Extract tables if available
        if (convertResult.getTables() != null) {
            doc.setTables(convertResult.getTables());

            // Summarize tables with LLM if configured
            if (llmProcessor != null && config.isSummarizeTables()) {
                summarizeTables(doc);
            }
        }

        return doc;
    }

    private ConvertResult convertToMarkdown(byte[] data, String filename) throws Exception {
        String ext = getFileExtension(filename).toLowerCase();

        switch (ext) {
            case "csv":
                return new CSVConverter().toMarkdown(data, filename);
            case "xlsx":
            case "xls":
                return new ExcelConverter().toMarkdown(data, filename);
            case "docx":
                return new DOCXConverter().toMarkdown(data, filename);
            case "txt":
                return new TXTConverter().toMarkdown(new String(data), filename);
            case "pdf":
                // For PDF, return as plain text - would need PDF library
                return new ConvertResult(true, new String(data), "pdf", new HashMap<>());
            case "md":
            case "markdown":
                return new ConvertResult(true, new String(data), "markdown", new HashMap<>());
            default:
                // Try as plain text
                return new ConvertResult(true, new String(data), "text", new HashMap<>());
        }
    }

    private void describeImages(RAGDocument doc) {
        for (ExtractedImage image : doc.getImages()) {
            try {
                if (image.getData() != null) {
                    String description = llmProcessor.describeImage(image.getData(), image.getFilename());
                    image.setDescription(description);
                }
            } catch (Exception e) {
                // Continue with other images
            }
        }
    }

    private void summarizeTables(RAGDocument doc) {
        for (ExtractedTable table : doc.getTables()) {
            try {
                if (table.getMarkdown() != null) {
                    String summary = llmProcessor.summarizeTable(table.getMarkdown());
                    table.setSummary(summary);
                }
            } catch (Exception e) {
                // Continue with other tables
            }
        }
    }

    private String getFileExtension(String filename) {
        int lastDot = filename.lastIndexOf('.');
        if (lastDot == -1 || lastDot == filename.length() - 1) {
            return "";
        }
        return filename.substring(lastDot + 1);
    }
}
