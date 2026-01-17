package com.docflow.config;

import java.util.Arrays;
import java.util.List;

/**
 * Configuration for text chunking.
 */
public class ChunkingConfig {

    public enum SplitBy {
        PARAGRAPH("paragraph"),
        SENTENCE("sentence"),
        TOKEN("token"),
        CHARACTER("character"),
        HEADING("heading");

        private final String value;

        SplitBy(String value) {
            this.value = value;
        }

        public String getValue() {
            return value;
        }
    }

    private int chunkSize = 1000;
    private int chunkOverlap = 200;
    private int minChunkSize = 100;
    private int maxChunkSize = 2000;
    private SplitBy splitBy = SplitBy.PARAGRAPH;
    private List<String> separators = Arrays.asList("\n\n", "\n", ". ", " ");
    private String tokenizer = "cl100k_base";
    private boolean respectHeadings = true;
    private boolean keepTablesTogether = true;
    private boolean keepCodeTogether = true;
    private boolean addChunkMarkers = true;
    private String markerFormat = "[CHUNK %d]";

    public ChunkingConfig() {
    }

    public static ChunkingConfig defaultConfig() {
        return new ChunkingConfig();
    }

    public void validate() {
        if (chunkSize <= 0) {
            throw new IllegalArgumentException("chunkSize must be positive");
        }
        if (chunkOverlap < 0) {
            throw new IllegalArgumentException("chunkOverlap cannot be negative");
        }
        if (chunkOverlap >= chunkSize) {
            throw new IllegalArgumentException("chunkOverlap must be less than chunkSize");
        }
    }

    // Getters and setters
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

    public int getMinChunkSize() {
        return minChunkSize;
    }

    public void setMinChunkSize(int minChunkSize) {
        this.minChunkSize = minChunkSize;
    }

    public int getMaxChunkSize() {
        return maxChunkSize;
    }

    public void setMaxChunkSize(int maxChunkSize) {
        this.maxChunkSize = maxChunkSize;
    }

    public SplitBy getSplitBy() {
        return splitBy;
    }

    public void setSplitBy(SplitBy splitBy) {
        this.splitBy = splitBy;
    }

    public List<String> getSeparators() {
        return separators;
    }

    public void setSeparators(List<String> separators) {
        this.separators = separators;
    }

    public String getTokenizer() {
        return tokenizer;
    }

    public void setTokenizer(String tokenizer) {
        this.tokenizer = tokenizer;
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

    public boolean isKeepCodeTogether() {
        return keepCodeTogether;
    }

    public void setKeepCodeTogether(boolean keepCodeTogether) {
        this.keepCodeTogether = keepCodeTogether;
    }

    public boolean isAddChunkMarkers() {
        return addChunkMarkers;
    }

    public void setAddChunkMarkers(boolean addChunkMarkers) {
        this.addChunkMarkers = addChunkMarkers;
    }

    public String getMarkerFormat() {
        return markerFormat;
    }

    public void setMarkerFormat(String markerFormat) {
        this.markerFormat = markerFormat;
    }
}
