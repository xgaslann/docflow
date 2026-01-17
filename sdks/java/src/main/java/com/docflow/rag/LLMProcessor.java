package com.docflow.rag;

import com.docflow.config.LLMConfig;

import java.io.*;
import java.net.URI;
import java.net.http.*;
import java.nio.charset.StandardCharsets;
import java.util.*;

/**
 * LLM Processor for image description and table summarization.
 * Supports OpenAI, Anthropic, Azure OpenAI, and Ollama providers.
 */
public class LLMProcessor {

    private final LLMConfig config;
    private final HttpClient httpClient;

    public LLMProcessor(LLMConfig config) {
        this.config = config;
        this.httpClient = HttpClient.newBuilder()
                .connectTimeout(java.time.Duration.ofSeconds(config.getTimeout()))
                .build();
    }

    /**
     * Describe an image using a Vision LLM.
     *
     * @param imageData Raw image bytes
     * @param context   Optional context about the image
     * @return Text description of the image
     */
    public String describeImage(byte[] imageData, String context) throws Exception {
        String base64Image = Base64.getEncoder().encodeToString(imageData);
        String prompt = buildImagePrompt(context);

        switch (config.getProvider().toLowerCase()) {
            case "openai":
            case "azure":
                return callOpenAIVision(base64Image, prompt);
            case "anthropic":
                return callAnthropicVision(base64Image, prompt);
            case "ollama":
                return callOllamaVision(base64Image, prompt);
            default:
                throw new IllegalArgumentException("Unsupported provider: " + config.getProvider());
        }
    }

    /**
     * Summarize a markdown table.
     *
     * @param tableMarkdown Table in markdown format
     * @return Summary of the table contents
     */
    public String summarizeTable(String tableMarkdown) throws Exception {
        String prompt = String.format("""
                Analyze this table and provide a concise summary of its contents.
                Focus on key data points, trends, and notable values.

                Table:
                %s

                Summary:""", tableMarkdown);

        return complete(prompt);
    }

    /**
     * Send a completion request to the LLM.
     *
     * @param prompt The prompt to send
     * @return The completion response
     */
    public String complete(String prompt) throws Exception {
        switch (config.getProvider().toLowerCase()) {
            case "openai":
                return callOpenAI(prompt);
            case "azure":
                return callAzureOpenAI(prompt);
            case "anthropic":
                return callAnthropic(prompt);
            case "ollama":
                return callOllama(prompt);
            default:
                throw new IllegalArgumentException("Unsupported provider: " + config.getProvider());
        }
    }

    private String buildImagePrompt(String context) {
        if (context != null && !context.isEmpty()) {
            return String.format("""
                    Describe this image in detail, focusing on its content and any text visible.
                    Context: %s

                    Provide a comprehensive description that would be useful for document search and retrieval.""",
                    context);
        }
        return """
                Describe this image in detail, focusing on its content and any text visible.
                Provide a comprehensive description that would be useful for document search and retrieval.""";
    }

    private String callOpenAI(String prompt) throws Exception {
        String requestBody = String.format("""
                {
                    "model": "%s",
                    "messages": [{"role": "user", "content": "%s"}],
                    "max_tokens": %d,
                    "temperature": %.1f
                }""",
                config.getModel(),
                escapeJson(prompt),
                config.getMaxTokens(),
                config.getTemperature());

        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(getOpenAIBaseUrl() + "/chat/completions"))
                .header("Content-Type", "application/json")
                .header("Authorization", "Bearer " + config.getApiKey())
                .POST(HttpRequest.BodyPublishers.ofString(requestBody))
                .build();

        HttpResponse<String> response = httpClient.send(request, HttpResponse.BodyHandlers.ofString());
        return extractOpenAIContent(response.body());
    }

    private String callAzureOpenAI(String prompt) throws Exception {
        String baseUrl = config.getBaseUrl();
        if (baseUrl == null || baseUrl.isEmpty()) {
            throw new IllegalArgumentException("Azure OpenAI requires baseUrl configuration");
        }

        String requestBody = String.format("""
                {
                    "messages": [{"role": "user", "content": "%s"}],
                    "max_tokens": %d,
                    "temperature": %.1f
                }""",
                escapeJson(prompt),
                config.getMaxTokens(),
                config.getTemperature());

        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(baseUrl + "/openai/deployments/" + config.getModel()
                        + "/chat/completions?api-version=2024-02-15-preview"))
                .header("Content-Type", "application/json")
                .header("api-key", config.getApiKey())
                .POST(HttpRequest.BodyPublishers.ofString(requestBody))
                .build();

        HttpResponse<String> response = httpClient.send(request, HttpResponse.BodyHandlers.ofString());
        return extractOpenAIContent(response.body());
    }

    private String callAnthropic(String prompt) throws Exception {
        String requestBody = String.format("""
                {
                    "model": "%s",
                    "max_tokens": %d,
                    "messages": [{"role": "user", "content": "%s"}]
                }""",
                config.getModel(),
                config.getMaxTokens(),
                escapeJson(prompt));

        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create("https://api.anthropic.com/v1/messages"))
                .header("Content-Type", "application/json")
                .header("x-api-key", config.getApiKey())
                .header("anthropic-version", "2023-06-01")
                .POST(HttpRequest.BodyPublishers.ofString(requestBody))
                .build();

        HttpResponse<String> response = httpClient.send(request, HttpResponse.BodyHandlers.ofString());
        return extractAnthropicContent(response.body());
    }

    private String callOllama(String prompt) throws Exception {
        String baseUrl = config.getBaseUrl() != null ? config.getBaseUrl() : "http://localhost:11434";

        String requestBody = String.format("""
                {
                    "model": "%s",
                    "prompt": "%s",
                    "stream": false
                }""",
                config.getModel(),
                escapeJson(prompt));

        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(baseUrl + "/api/generate"))
                .header("Content-Type", "application/json")
                .POST(HttpRequest.BodyPublishers.ofString(requestBody))
                .build();

        HttpResponse<String> response = httpClient.send(request, HttpResponse.BodyHandlers.ofString());
        return extractOllamaContent(response.body());
    }

    private String callOpenAIVision(String base64Image, String prompt) throws Exception {
        String requestBody = String.format("""
                {
                    "model": "%s",
                    "messages": [{
                        "role": "user",
                        "content": [
                            {"type": "text", "text": "%s"},
                            {"type": "image_url", "image_url": {"url": "data:image/png;base64,%s"}}
                        ]
                    }],
                    "max_tokens": %d
                }""",
                config.getModel(),
                escapeJson(prompt),
                base64Image,
                config.getMaxTokens());

        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(getOpenAIBaseUrl() + "/chat/completions"))
                .header("Content-Type", "application/json")
                .header("Authorization", "Bearer " + config.getApiKey())
                .POST(HttpRequest.BodyPublishers.ofString(requestBody))
                .build();

        HttpResponse<String> response = httpClient.send(request, HttpResponse.BodyHandlers.ofString());
        return extractOpenAIContent(response.body());
    }

    private String callAnthropicVision(String base64Image, String prompt) throws Exception {
        String requestBody = String.format("""
                {
                    "model": "%s",
                    "max_tokens": %d,
                    "messages": [{
                        "role": "user",
                        "content": [
                            {"type": "image", "source": {"type": "base64", "media_type": "image/png", "data": "%s"}},
                            {"type": "text", "text": "%s"}
                        ]
                    }]
                }""",
                config.getModel(),
                config.getMaxTokens(),
                base64Image,
                escapeJson(prompt));

        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create("https://api.anthropic.com/v1/messages"))
                .header("Content-Type", "application/json")
                .header("x-api-key", config.getApiKey())
                .header("anthropic-version", "2023-06-01")
                .POST(HttpRequest.BodyPublishers.ofString(requestBody))
                .build();

        HttpResponse<String> response = httpClient.send(request, HttpResponse.BodyHandlers.ofString());
        return extractAnthropicContent(response.body());
    }

    private String callOllamaVision(String base64Image, String prompt) throws Exception {
        String baseUrl = config.getBaseUrl() != null ? config.getBaseUrl() : "http://localhost:11434";

        String requestBody = String.format("""
                {
                    "model": "%s",
                    "prompt": "%s",
                    "images": ["%s"],
                    "stream": false
                }""",
                config.getModel(),
                escapeJson(prompt),
                base64Image);

        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(baseUrl + "/api/generate"))
                .header("Content-Type", "application/json")
                .POST(HttpRequest.BodyPublishers.ofString(requestBody))
                .build();

        HttpResponse<String> response = httpClient.send(request, HttpResponse.BodyHandlers.ofString());
        return extractOllamaContent(response.body());
    }

    private String getOpenAIBaseUrl() {
        if (config.getBaseUrl() != null && !config.getBaseUrl().isEmpty()) {
            return config.getBaseUrl();
        }
        return "https://api.openai.com/v1";
    }

    private String extractOpenAIContent(String response) {
        // Simple JSON parsing - extract content from choices[0].message.content
        int contentStart = response.indexOf("\"content\":");
        if (contentStart == -1)
            return response;

        int valueStart = response.indexOf("\"", contentStart + 10) + 1;
        int valueEnd = response.indexOf("\"", valueStart);

        while (valueEnd > 0 && response.charAt(valueEnd - 1) == '\\') {
            valueEnd = response.indexOf("\"", valueEnd + 1);
        }

        if (valueEnd > valueStart) {
            return response.substring(valueStart, valueEnd).replace("\\n", "\n").replace("\\\"", "\"");
        }
        return response;
    }

    private String extractAnthropicContent(String response) {
        int textStart = response.indexOf("\"text\":");
        if (textStart == -1)
            return response;

        int valueStart = response.indexOf("\"", textStart + 7) + 1;
        int valueEnd = response.indexOf("\"", valueStart);

        while (valueEnd > 0 && response.charAt(valueEnd - 1) == '\\') {
            valueEnd = response.indexOf("\"", valueEnd + 1);
        }

        if (valueEnd > valueStart) {
            return response.substring(valueStart, valueEnd).replace("\\n", "\n").replace("\\\"", "\"");
        }
        return response;
    }

    private String extractOllamaContent(String response) {
        int responseStart = response.indexOf("\"response\":");
        if (responseStart == -1)
            return response;

        int valueStart = response.indexOf("\"", responseStart + 11) + 1;
        int valueEnd = response.indexOf("\"", valueStart);

        while (valueEnd > 0 && response.charAt(valueEnd - 1) == '\\') {
            valueEnd = response.indexOf("\"", valueEnd + 1);
        }

        if (valueEnd > valueStart) {
            return response.substring(valueStart, valueEnd).replace("\\n", "\n").replace("\\\"", "\"");
        }
        return response;
    }

    private String escapeJson(String text) {
        return text
                .replace("\\", "\\\\")
                .replace("\"", "\\\"")
                .replace("\n", "\\n")
                .replace("\r", "\\r")
                .replace("\t", "\\t");
    }
}
