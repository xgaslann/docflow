package com.docflow.models;

import java.util.List;

/**
 * Result of a PDF conversion.
 */
public class PDFResult {
    private boolean success;
    private List<String> filePaths;
    private String error;

    public PDFResult() {
    }

    public PDFResult(boolean success, List<String> filePaths) {
        this.success = success;
        this.filePaths = filePaths;
    }

    public PDFResult(boolean success, String error) {
        this.success = success;
        this.error = error;
    }

    // Getters and Setters
    public boolean isSuccess() {
        return success;
    }

    public void setSuccess(boolean success) {
        this.success = success;
    }

    public List<String> getFilePaths() {
        return filePaths;
    }

    public void setFilePaths(List<String> filePaths) {
        this.filePaths = filePaths;
    }

    public String getError() {
        return error;
    }

    public void setError(String error) {
        this.error = error;
    }
}
