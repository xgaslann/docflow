package com.docflow.models;

import java.util.List;
import java.util.Map;

/**
 * Comprehensive document metadata.
 */
public class DocumentMetadata {
    // Basic info
    private String title;
    private String author;
    private String createdDate;
    private String modifiedDate;

    // Structure
    private List<HeadingInfo> headings;
    private Map<String, Object> headingTree;
    private List<TOCItem> tableOfContents;

    // Statistics
    private int wordCount;
    private int charCount;
    private int pageCount;
    private int imageCount;
    private int tableCount;

    // Language & Content
    private String language;
    private List<String> keywords;
    private List<String> entities;

    // LLM-generated
    private String summary;
    private List<String> keyPoints;

    // Constructors
    public DocumentMetadata() {
    }

    // Getters and Setters
    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getAuthor() {
        return author;
    }

    public void setAuthor(String author) {
        this.author = author;
    }

    public String getCreatedDate() {
        return createdDate;
    }

    public void setCreatedDate(String createdDate) {
        this.createdDate = createdDate;
    }

    public String getModifiedDate() {
        return modifiedDate;
    }

    public void setModifiedDate(String modifiedDate) {
        this.modifiedDate = modifiedDate;
    }

    public List<HeadingInfo> getHeadings() {
        return headings;
    }

    public void setHeadings(List<HeadingInfo> headings) {
        this.headings = headings;
    }

    public Map<String, Object> getHeadingTree() {
        return headingTree;
    }

    public void setHeadingTree(Map<String, Object> headingTree) {
        this.headingTree = headingTree;
    }

    public List<TOCItem> getTableOfContents() {
        return tableOfContents;
    }

    public void setTableOfContents(List<TOCItem> tableOfContents) {
        this.tableOfContents = tableOfContents;
    }

    public int getWordCount() {
        return wordCount;
    }

    public void setWordCount(int wordCount) {
        this.wordCount = wordCount;
    }

    public int getCharCount() {
        return charCount;
    }

    public void setCharCount(int charCount) {
        this.charCount = charCount;
    }

    public int getPageCount() {
        return pageCount;
    }

    public void setPageCount(int pageCount) {
        this.pageCount = pageCount;
    }

    public int getImageCount() {
        return imageCount;
    }

    public void setImageCount(int imageCount) {
        this.imageCount = imageCount;
    }

    public int getTableCount() {
        return tableCount;
    }

    public void setTableCount(int tableCount) {
        this.tableCount = tableCount;
    }

    public String getLanguage() {
        return language;
    }

    public void setLanguage(String language) {
        this.language = language;
    }

    public List<String> getKeywords() {
        return keywords;
    }

    public void setKeywords(List<String> keywords) {
        this.keywords = keywords;
    }

    public List<String> getEntities() {
        return entities;
    }

    public void setEntities(List<String> entities) {
        this.entities = entities;
    }

    public String getSummary() {
        return summary;
    }

    public void setSummary(String summary) {
        this.summary = summary;
    }

    public List<String> getKeyPoints() {
        return keyPoints;
    }

    public void setKeyPoints(List<String> keyPoints) {
        this.keyPoints = keyPoints;
    }
}
