package com.docflow.models;

/**
 * Enum for chunking strategies.
 */
public enum ChunkingStrategy {
    SIMPLE("simple"),
    HEADING_AWARE("heading_aware"),
    DOCUMENT_INTELLIGENCE("doc_intel"),
    SEMANTIC("semantic");

    private final String value;

    ChunkingStrategy(String value) {
        this.value = value;
    }

    public String getValue() {
        return value;
    }

    public static ChunkingStrategy fromValue(String value) {
        for (ChunkingStrategy strategy : values()) {
            if (strategy.value.equals(value)) {
                return strategy;
            }
        }
        throw new IllegalArgumentException("Unknown chunking strategy: " + value);
    }
}
