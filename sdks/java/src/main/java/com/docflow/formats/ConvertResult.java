package com.docflow.formats;

import com.docflow.models.ExtractedImage;
import com.docflow.models.ExtractedTable;
import java.util.*;

/**
 * Result of a format conversion.
 */
public class ConvertResult {

    private boolean success;
    private String content;
    private String format;
    private String error;
    private Map<String, Object> metadata = new HashMap<>();
    private List<ExtractedImage> images = new ArrayList<>();
    private List<ExtractedTable> tables = new ArrayList<>();

    public ConvertResult() {
    }

    /**
     * Create a result with success/error status.
     */
    public ConvertResult(boolean success, String errorOrContent) {
        this.success = success;
        if (success) {
            this.content = errorOrContent;
        } else {
            this.error = errorOrContent;
        }
    }

    /**
     * Create a successful result with content and format.
     */
    public ConvertResult(boolean success, String content, String format, Map<String, Object> metadata) {
        this.success = success;
        this.content = content;
        this.format = format;
        if (metadata != null) {
            this.metadata = metadata;
        }
    }

    public static ConvertResult success(String content, String format) {
        ConvertResult result = new ConvertResult();
        result.success = true;
        result.content = content;
        result.format = format;
        return result;
    }

    public static ConvertResult error(String error) {
        ConvertResult result = new ConvertResult();
        result.success = false;
        result.error = error;
        return result;
    }

    // Getters and setters
    public boolean isSuccess() {
        return success;
    }

    public void setSuccess(boolean success) {
        this.success = success;
    }

    public String getContent() {
        return content;
    }

    public void setContent(String content) {
        this.content = content;
    }

    public String getFormat() {
        return format;
    }

    public void setFormat(String format) {
        this.format = format;
    }

    public String getError() {
        return error;
    }

    public void setError(String error) {
        this.error = error;
    }

    public Map<String, Object> getMetadata() {
        return metadata;
    }

    public void setMetadata(Map<String, Object> metadata) {
        this.metadata = metadata;
    }

    public ConvertResult withMetadata(String key, Object value) {
        this.metadata.put(key, value);
        return this;
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
}
