package com.docflow.models;

import java.util.*;

/**
 * Represents a document processed for RAG.
 */
public class RAGDocument {
    private String id;
    private String filename;
    private String content;
    private List<Chunk> chunks = new ArrayList<>();
    private List<ExtractedImage> images = new ArrayList<>();
    private List<ExtractedTable> tables = new ArrayList<>();
    private DocumentMetadata metadata;

    public RAGDocument() {
    }

    public RAGDocument(String id, String filename, String content) {
        this.id = id;
        this.filename = filename;
        this.content = content;
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getFilename() {
        return filename;
    }

    public void setFilename(String filename) {
        this.filename = filename;
    }

    public String getContent() {
        return content;
    }

    public void setContent(String content) {
        this.content = content;
    }

    public List<Chunk> getChunks() {
        return chunks;
    }

    public void setChunks(List<Chunk> chunks) {
        this.chunks = chunks;
    }

    public List<ExtractedImage> getImages() {
        return images;
    }

    public void setImages(List<ExtractedImage> images) {
        this.images = images;
    }

    public List<ExtractedTable> getTables() {
        return tables;
    }

    public void setTables(List<ExtractedTable> tables) {
        this.tables = tables;
    }

    public DocumentMetadata getMetadata() {
        return metadata;
    }

    public void setMetadata(DocumentMetadata metadata) {
        this.metadata = metadata;
    }
}
