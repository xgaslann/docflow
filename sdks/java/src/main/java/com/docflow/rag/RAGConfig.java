package com.docflow.rag;

import com.docflow.config.LLMConfig;

/**
 * RAG configuration options.
 */
public class RAGConfig {
    private boolean enabled = true;
    private int chunkSize = 1000;
    private int chunkOverlap = 200;

    private boolean extractImages = true;
    private boolean extractTables = true;
    private boolean preserveMetadata = true;

    private boolean describeImages = false;
    private boolean summarizeTables = false;
    private LLMConfig llmConfig;

    private boolean respectHeadings = true;
    private boolean keepTablesTogether = true;
    private boolean addChunkMarkers = true;

    public static RAGConfig defaultConfig() {
        return new RAGConfig();
    }

    // Getters and setters
    public boolean isEnabled() {
        return enabled;
    }

    public void setEnabled(boolean enabled) {
        this.enabled = enabled;
    }

    public int getChunkSize() {
        return chunkSize;
    }

    public void setChunkSize(int chunkSize) {
        this.chunkSize = chunkSize;
    }

    public int getChunkOverlap() {
        return chunkOverlap;
    }

    public void setChunkOverlap(int chunkOverlap) {
        this.chunkOverlap = chunkOverlap;
    }

    public boolean isExtractImages() {
        return extractImages;
    }

    public void setExtractImages(boolean extractImages) {
        this.extractImages = extractImages;
    }

    public boolean isExtractTables() {
        return extractTables;
    }

    public void setExtractTables(boolean extractTables) {
        this.extractTables = extractTables;
    }

    public boolean isPreserveMetadata() {
        return preserveMetadata;
    }

    public void setPreserveMetadata(boolean preserveMetadata) {
        this.preserveMetadata = preserveMetadata;
    }

    public boolean isDescribeImages() {
        return describeImages;
    }

    public void setDescribeImages(boolean describeImages) {
        this.describeImages = describeImages;
    }

    public LLMConfig getLlmConfig() {
        return llmConfig;
    }

    public void setLlmConfig(LLMConfig llmConfig) {
        this.llmConfig = llmConfig;
    }

    public boolean isSummarizeTables() {
        return summarizeTables;
    }

    public void setSummarizeTables(boolean summarizeTables) {
        this.summarizeTables = summarizeTables;
    }

    public boolean isRespectHeadings() {
        return respectHeadings;
    }

    public void setRespectHeadings(boolean respectHeadings) {
        this.respectHeadings = respectHeadings;
    }

    public boolean isKeepTablesTogether() {
        return keepTablesTogether;
    }

    public void setKeepTablesTogether(boolean keepTablesTogether) {
        this.keepTablesTogether = keepTablesTogether;
    }

    public boolean isAddChunkMarkers() {
        return addChunkMarkers;
    }

    public void setAddChunkMarkers(boolean addChunkMarkers) {
        this.addChunkMarkers = addChunkMarkers;
    }
}
