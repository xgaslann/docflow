package com.docflow.models;

/**
 * Enum for LLM processing modes.
 */
public enum LLMProcessingMode {
    IMAGES("images"),
    TABLES("tables"),
    TEXT("text"),
    ALL("all");

    private final String value;

    LLMProcessingMode(String value) {
        this.value = value;
    }

    public String getValue() {
        return value;
    }

    public static LLMProcessingMode fromValue(String value) {
        for (LLMProcessingMode mode : values()) {
            if (mode.value.equals(value)) {
                return mode;
            }
        }
        throw new IllegalArgumentException("Unknown LLM processing mode: " + value);
    }
}
