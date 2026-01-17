package com.docflow.models;

/**
 * Options for MD to PDF conversion.
 */
public class ConvertOptions {
    /**
     * Merge mode for conversion.
     * "merged" - combine all files into one PDF
     * "separate" - create separate PDFs for each file
     */
    private String mergeMode = "separate";

    /**
     * Name for the output file (used in merged mode).
     */
    private String outputName;

    public ConvertOptions() {
    }

    public ConvertOptions(String mergeMode) {
        this.mergeMode = mergeMode;
    }

    public ConvertOptions(String mergeMode, String outputName) {
        this.mergeMode = mergeMode;
        this.outputName = outputName;
    }

    // Getters and Setters
    public String getMergeMode() {
        return mergeMode;
    }

    public void setMergeMode(String mergeMode) {
        this.mergeMode = mergeMode;
    }

    public String getOutputName() {
        return outputName;
    }

    public void setOutputName(String outputName) {
        this.outputName = outputName;
    }

    public static ConvertOptions merged() {
        return new ConvertOptions("merged");
    }

    public static ConvertOptions merged(String outputName) {
        return new ConvertOptions("merged", outputName);
    }

    public static ConvertOptions separate() {
        return new ConvertOptions("separate");
    }
}
