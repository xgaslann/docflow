package com.docflow.models;

/**
 * Enum for output formats.
 */
public enum OutputFormat {
    MARKDOWN("markdown"),
    PDF("pdf"),
    HTML("html");

    private final String value;

    OutputFormat(String value) {
        this.value = value;
    }

    public String getValue() {
        return value;
    }

    public static OutputFormat fromValue(String value) {
        for (OutputFormat format : values()) {
            if (format.value.equals(value)) {
                return format;
            }
        }
        throw new IllegalArgumentException("Unknown output format: " + value);
    }
}
