package com.docflow.config;

import java.util.*;

/**
 * Configuration for LLM integration.
 */
public class LLMConfiguration {

    public enum Provider {
        OPENAI("openai"),
        AZURE_OPENAI("azure_openai"),
        ANTHROPIC("anthropic"),
        OLLAMA("ollama"),
        GOOGLE("google");

        private final String value;

        Provider(String value) {
            this.value = value;
        }

        public String getValue() {
            return value;
        }
    }

    private Provider provider = Provider.OPENAI;
    private String model = "gpt-4-vision-preview";
    private String apiKey = "";
    private LLMPrompts prompts = new LLMPrompts();

    // OpenAI
    private String organization = "";
    private String baseUrl = null;

    // Azure OpenAI
    private String azureEndpoint = "";
    private String azureDeployment = "";
    private String apiVersion = "2024-02-01";

    // Ollama
    private String ollamaBaseUrl = "http://localhost:11434";

    // Generation Parameters
    private double temperature = 0.7;
    private int maxTokens = 1000;
    private double topP = 1.0;
    private double frequencyPenalty = 0.0;
    private double presencePenalty = 0.0;
    private List<String> stopSequences = new ArrayList<>();

    // Vision Parameters
    private String detail = "auto";
    private long maxImageSize = 20 * 1024 * 1024;
    private List<String> supportedFormats = Arrays.asList("png", "jpg", "jpeg", "gif", "webp");

    // Retry & Timeout
    private int timeout = 60;
    private int retryCount = 3;
    private double retryDelay = 1.0;

    // Batch Processing
    private int batchSize = 5;
    private int concurrentRequests = 3;

    private String responseFormat = "text";

    public LLMConfiguration() {
    }

    public static LLMConfiguration defaultConfig() {
        return new LLMConfiguration();
    }

    public void validate() {
        if (apiKey == null || apiKey.isEmpty()) {
            if (provider != Provider.OLLAMA) {
                throw new IllegalArgumentException("API key is required for " + provider.getValue());
            }
        }
        if (provider == Provider.AZURE_OPENAI) {
            if (azureEndpoint == null || azureEndpoint.isEmpty()) {
                throw new IllegalArgumentException("Azure endpoint is required");
            }
        }
    }

    // Getters and setters
    public Provider getProvider() {
        return provider;
    }

    public void setProvider(Provider provider) {
        this.provider = provider;
    }

    public String getModel() {
        return model;
    }

    public void setModel(String model) {
        this.model = model;
    }

    public String getApiKey() {
        return apiKey;
    }

    public void setApiKey(String apiKey) {
        this.apiKey = apiKey;
    }

    public LLMPrompts getPrompts() {
        return prompts;
    }

    public void setPrompts(LLMPrompts prompts) {
        this.prompts = prompts;
    }

    public String getOrganization() {
        return organization;
    }

    public void setOrganization(String organization) {
        this.organization = organization;
    }

    public String getBaseUrl() {
        return baseUrl;
    }

    public void setBaseUrl(String baseUrl) {
        this.baseUrl = baseUrl;
    }

    public String getAzureEndpoint() {
        return azureEndpoint;
    }

    public void setAzureEndpoint(String azureEndpoint) {
        this.azureEndpoint = azureEndpoint;
    }

    public String getAzureDeployment() {
        return azureDeployment;
    }

    public void setAzureDeployment(String azureDeployment) {
        this.azureDeployment = azureDeployment;
    }

    public String getApiVersion() {
        return apiVersion;
    }

    public void setApiVersion(String apiVersion) {
        this.apiVersion = apiVersion;
    }

    public String getOllamaBaseUrl() {
        return ollamaBaseUrl;
    }

    public void setOllamaBaseUrl(String ollamaBaseUrl) {
        this.ollamaBaseUrl = ollamaBaseUrl;
    }

    public double getTemperature() {
        return temperature;
    }

    public void setTemperature(double temperature) {
        this.temperature = temperature;
    }

    public int getMaxTokens() {
        return maxTokens;
    }

    public void setMaxTokens(int maxTokens) {
        this.maxTokens = maxTokens;
    }

    public double getTopP() {
        return topP;
    }

    public void setTopP(double topP) {
        this.topP = topP;
    }

    public double getFrequencyPenalty() {
        return frequencyPenalty;
    }

    public void setFrequencyPenalty(double frequencyPenalty) {
        this.frequencyPenalty = frequencyPenalty;
    }

    public double getPresencePenalty() {
        return presencePenalty;
    }

    public void setPresencePenalty(double presencePenalty) {
        this.presencePenalty = presencePenalty;
    }

    public List<String> getStopSequences() {
        return stopSequences;
    }

    public void setStopSequences(List<String> stopSequences) {
        this.stopSequences = stopSequences;
    }

    public String getDetail() {
        return detail;
    }

    public void setDetail(String detail) {
        this.detail = detail;
    }

    public long getMaxImageSize() {
        return maxImageSize;
    }

    public void setMaxImageSize(long maxImageSize) {
        this.maxImageSize = maxImageSize;
    }

    public List<String> getSupportedFormats() {
        return supportedFormats;
    }

    public void setSupportedFormats(List<String> supportedFormats) {
        this.supportedFormats = supportedFormats;
    }

    public int getTimeout() {
        return timeout;
    }

    public void setTimeout(int timeout) {
        this.timeout = timeout;
    }

    public int getRetryCount() {
        return retryCount;
    }

    public void setRetryCount(int retryCount) {
        this.retryCount = retryCount;
    }

    public double getRetryDelay() {
        return retryDelay;
    }

    public void setRetryDelay(double retryDelay) {
        this.retryDelay = retryDelay;
    }

    public int getBatchSize() {
        return batchSize;
    }

    public void setBatchSize(int batchSize) {
        this.batchSize = batchSize;
    }

    public int getConcurrentRequests() {
        return concurrentRequests;
    }

    public void setConcurrentRequests(int concurrentRequests) {
        this.concurrentRequests = concurrentRequests;
    }

    public String getResponseFormat() {
        return responseFormat;
    }

    public void setResponseFormat(String responseFormat) {
        this.responseFormat = responseFormat;
    }
}
