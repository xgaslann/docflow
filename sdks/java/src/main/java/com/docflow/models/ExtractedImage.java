package com.docflow.models;

/**
 * Represents an extracted image from a document.
 */
public class ExtractedImage {
    private String filename;
    private byte[] data;
    private String format;
    private String description;
    private int pageNumber;

    public ExtractedImage() {
    }

    public ExtractedImage(String filename, byte[] data, String format) {
        this.filename = filename;
        this.data = data;
        this.format = format;
    }

    public String getFilename() {
        return filename;
    }

    public void setFilename(String filename) {
        this.filename = filename;
    }

    public byte[] getData() {
        return data;
    }

    public void setData(byte[] data) {
        this.data = data;
    }

    public String getFormat() {
        return format;
    }

    public void setFormat(String format) {
        this.format = format;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public int getPageNumber() {
        return pageNumber;
    }

    public void setPageNumber(int pageNumber) {
        this.pageNumber = pageNumber;
    }
}
