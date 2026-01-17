package com.docflow.config;

import java.util.*;

/**
 * Configuration for metadata extraction.
 */
public class MetadataConfig {

    private List<String> includeFields = Arrays.asList("title", "headings", "table_of_contents", "word_count");
    private List<String> excludeFields = new ArrayList<>();
    private Map<String, Object> customFields = new HashMap<>();

    private boolean extractTitle = true;
    private boolean extractAuthor = true;
    private boolean extractHeadings = true;
    private boolean extractToc = true;
    private boolean extractWordCount = true;
    private boolean extractPageCount = true;
    private boolean extractEntities = false;
    private boolean extractSummary = false;
    private boolean extractKeyPoints = false;
    private boolean extractLanguage = false;

    private int maxHeadingLevel = 6;
    private boolean buildHeadingTree = true;
    private int tocMaxDepth = 3;
    private boolean tocIncludePageNumbers = true;

    public MetadataConfig() {
    }

    public static MetadataConfig defaultConfig() {
        return new MetadataConfig();
    }

    public boolean shouldExtract(String fieldName) {
        if (excludeFields.contains(fieldName)) {
            return false;
        }
        if (!includeFields.isEmpty() && !includeFields.contains(fieldName)) {
            return false;
        }
        return true;
    }

    public void addCustomField(String name, Object value) {
        customFields.put(name, value);
    }

    public void removeField(String name) {
        if (!excludeFields.contains(name)) {
            excludeFields.add(name);
        }
    }

    // Getters and setters
    public List<String> getIncludeFields() {
        return includeFields;
    }

    public void setIncludeFields(List<String> includeFields) {
        this.includeFields = includeFields;
    }

    public List<String> getExcludeFields() {
        return excludeFields;
    }

    public void setExcludeFields(List<String> excludeFields) {
        this.excludeFields = excludeFields;
    }

    public Map<String, Object> getCustomFields() {
        return customFields;
    }

    public void setCustomFields(Map<String, Object> customFields) {
        this.customFields = customFields;
    }

    public boolean isExtractTitle() {
        return extractTitle;
    }

    public void setExtractTitle(boolean extractTitle) {
        this.extractTitle = extractTitle;
    }

    public boolean isExtractAuthor() {
        return extractAuthor;
    }

    public void setExtractAuthor(boolean extractAuthor) {
        this.extractAuthor = extractAuthor;
    }

    public boolean isExtractHeadings() {
        return extractHeadings;
    }

    public void setExtractHeadings(boolean extractHeadings) {
        this.extractHeadings = extractHeadings;
    }

    public boolean isExtractToc() {
        return extractToc;
    }

    public void setExtractToc(boolean extractToc) {
        this.extractToc = extractToc;
    }

    public boolean isExtractWordCount() {
        return extractWordCount;
    }

    public void setExtractWordCount(boolean extractWordCount) {
        this.extractWordCount = extractWordCount;
    }

    public boolean isExtractPageCount() {
        return extractPageCount;
    }

    public void setExtractPageCount(boolean extractPageCount) {
        this.extractPageCount = extractPageCount;
    }

    public boolean isExtractEntities() {
        return extractEntities;
    }

    public void setExtractEntities(boolean extractEntities) {
        this.extractEntities = extractEntities;
    }

    public boolean isExtractSummary() {
        return extractSummary;
    }

    public void setExtractSummary(boolean extractSummary) {
        this.extractSummary = extractSummary;
    }

    public boolean isExtractKeyPoints() {
        return extractKeyPoints;
    }

    public void setExtractKeyPoints(boolean extractKeyPoints) {
        this.extractKeyPoints = extractKeyPoints;
    }

    public boolean isExtractLanguage() {
        return extractLanguage;
    }

    public void setExtractLanguage(boolean extractLanguage) {
        this.extractLanguage = extractLanguage;
    }

    public int getMaxHeadingLevel() {
        return maxHeadingLevel;
    }

    public void setMaxHeadingLevel(int maxHeadingLevel) {
        this.maxHeadingLevel = maxHeadingLevel;
    }

    public boolean isBuildHeadingTree() {
        return buildHeadingTree;
    }

    public void setBuildHeadingTree(boolean buildHeadingTree) {
        this.buildHeadingTree = buildHeadingTree;
    }

    public int getTocMaxDepth() {
        return tocMaxDepth;
    }

    public void setTocMaxDepth(int tocMaxDepth) {
        this.tocMaxDepth = tocMaxDepth;
    }

    public boolean isTocIncludePageNumbers() {
        return tocIncludePageNumbers;
    }

    public void setTocIncludePageNumbers(boolean tocIncludePageNumbers) {
        this.tocIncludePageNumbers = tocIncludePageNumbers;
    }
}
