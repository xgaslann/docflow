package com.docflow.config;

import java.util.*;

/**
 * Custom prompts for LLM processing.
 */
public class LLMPrompts {

    private String imageDescription = "Describe this image in detail for use in a document retrieval system.";
    private String imageOcr = "Extract all text visible in this image.";
    private String tableAnalysis = "Analyze this table and provide key insights.";
    private String textSummary = "Summarize the following content concisely.";
    private String textKeyPoints = "Extract the most important key points.";
    private String entityExtraction = "Extract all named entities from this content.";
    private String keywordExtraction = "Extract the most relevant keywords.";
    private Map<String, String> custom = new HashMap<>();

    public LLMPrompts() {
    }

    public String getPrompt(String name) {
        if (custom.containsKey(name)) {
            return custom.get(name);
        }
        return null;
    }

    public void setCustomPrompt(String name, String prompt) {
        custom.put(name, prompt);
    }

    // Getters and setters
    public String getImageDescription() {
        return imageDescription;
    }

    public void setImageDescription(String imageDescription) {
        this.imageDescription = imageDescription;
    }

    public String getImageOcr() {
        return imageOcr;
    }

    public void setImageOcr(String imageOcr) {
        this.imageOcr = imageOcr;
    }

    public String getTableAnalysis() {
        return tableAnalysis;
    }

    public void setTableAnalysis(String tableAnalysis) {
        this.tableAnalysis = tableAnalysis;
    }

    public String getTextSummary() {
        return textSummary;
    }

    public void setTextSummary(String textSummary) {
        this.textSummary = textSummary;
    }

    public String getTextKeyPoints() {
        return textKeyPoints;
    }

    public void setTextKeyPoints(String textKeyPoints) {
        this.textKeyPoints = textKeyPoints;
    }

    public String getEntityExtraction() {
        return entityExtraction;
    }

    public void setEntityExtraction(String entityExtraction) {
        this.entityExtraction = entityExtraction;
    }

    public String getKeywordExtraction() {
        return keywordExtraction;
    }

    public void setKeywordExtraction(String keywordExtraction) {
        this.keywordExtraction = keywordExtraction;
    }

    public Map<String, String> getCustom() {
        return custom;
    }

    public void setCustom(Map<String, String> custom) {
        this.custom = custom;
    }
}
