package com.docflow.search;

import com.docflow.config.AISearchConfig;

import java.io.*;
import java.net.URI;
import java.net.http.*;
import java.util.*;

/**
 * Azure AI Search client for hybrid search.
 */
public class AzureAISearch {

    private final AISearchConfig config;
    private final HttpClient httpClient;

    public AzureAISearch(AISearchConfig config) {
        this.config = config;
        this.httpClient = HttpClient.newBuilder()
                .connectTimeout(java.time.Duration.ofSeconds(30))
                .build();
    }

    public AzureAISearch(String endpoint, String apiKey, String indexName) {
        this.config = new AISearchConfig();
        this.config.setEndpoint(endpoint);
        this.config.setApiKey(apiKey);
        this.config.setIndexName(indexName);
        this.httpClient = HttpClient.newBuilder()
                .connectTimeout(java.time.Duration.ofSeconds(30))
                .build();
    }

    /**
     * Perform a keyword search.
     */
    public List<SearchResult> search(String query, int top) throws Exception {
        return search(query, null, top, false);
    }

    /**
     * Perform a vector search.
     */
    public List<SearchResult> vectorSearch(float[] vector, int top) throws Exception {
        return search(null, vector, top, false);
    }

    /**
     * Perform a hybrid search (keyword + vector).
     */
    public List<SearchResult> hybridSearch(String query, float[] vector, int top) throws Exception {
        return search(query, vector, top, true);
    }

    /**
     * Search with all options.
     */
    public List<SearchResult> search(String query, float[] vector, int top, boolean hybrid) throws Exception {
        StringBuilder requestBody = new StringBuilder("{");

        if (query != null) {
            requestBody.append(String.format("\"search\": \"%s\",", escapeJson(query)));
        }

        if (vector != null) {
            requestBody.append("\"vectorQueries\": [{");
            requestBody.append(String.format("\"kind\": \"vector\","));
            requestBody.append(String.format("\"vector\": %s,", arrayToJson(vector)));
            requestBody.append(String.format("\"fields\": \"%s\",", config.getVectorFields()));
            requestBody.append(String.format("\"k\": %d", top));
            requestBody.append("}],");
        }

        if (hybrid && query != null && vector != null) {
            requestBody.append("\"searchMode\": \"all\",");
        }

        requestBody.append(String.format("\"top\": %d,", top));
        requestBody.append("\"select\": \"id,content,metadata\"");
        requestBody.append("}");

        String url = String.format("%s/indexes/%s/docs/search?api-version=%s",
                config.getEndpoint(), config.getIndexName(), config.getApiVersion());

        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(url))
                .header("Content-Type", "application/json")
                .header("api-key", config.getApiKey())
                .POST(HttpRequest.BodyPublishers.ofString(requestBody.toString()))
                .build();

        HttpResponse<String> response = httpClient.send(request, HttpResponse.BodyHandlers.ofString());

        if (response.statusCode() != 200) {
            throw new RuntimeException("Search failed: " + response.body());
        }

        return parseSearchResults(response.body());
    }

    /**
     * Upload documents to the index.
     */
    public void uploadDocuments(List<Map<String, Object>> documents) throws Exception {
        StringBuilder requestBody = new StringBuilder("{\"value\": [");

        for (int i = 0; i < documents.size(); i++) {
            if (i > 0)
                requestBody.append(",");
            requestBody.append(mapToJson(documents.get(i)));
        }

        requestBody.append("]}");

        String url = String.format("%s/indexes/%s/docs/index?api-version=%s",
                config.getEndpoint(), config.getIndexName(), config.getApiVersion());

        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(url))
                .header("Content-Type", "application/json")
                .header("api-key", config.getApiKey())
                .POST(HttpRequest.BodyPublishers.ofString(requestBody.toString()))
                .build();

        HttpResponse<String> response = httpClient.send(request, HttpResponse.BodyHandlers.ofString());

        if (response.statusCode() != 200 && response.statusCode() != 207) {
            throw new RuntimeException("Upload failed: " + response.body());
        }
    }

    /**
     * Delete documents by ID.
     */
    public void deleteDocuments(List<String> ids) throws Exception {
        StringBuilder requestBody = new StringBuilder("{\"value\": [");

        for (int i = 0; i < ids.size(); i++) {
            if (i > 0)
                requestBody.append(",");
            requestBody.append(String.format("{\"@search.action\": \"delete\", \"id\": \"%s\"}", ids.get(i)));
        }

        requestBody.append("]}");

        String url = String.format("%s/indexes/%s/docs/index?api-version=%s",
                config.getEndpoint(), config.getIndexName(), config.getApiVersion());

        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(url))
                .header("Content-Type", "application/json")
                .header("api-key", config.getApiKey())
                .POST(HttpRequest.BodyPublishers.ofString(requestBody.toString()))
                .build();

        HttpResponse<String> response = httpClient.send(request, HttpResponse.BodyHandlers.ofString());

        if (response.statusCode() != 200 && response.statusCode() != 207) {
            throw new RuntimeException("Delete failed: " + response.body());
        }
    }

    private List<SearchResult> parseSearchResults(String json) {
        List<SearchResult> results = new ArrayList<>();

        // Simple JSON parsing for search results
        int valueStart = json.indexOf("\"value\":");
        if (valueStart == -1)
            return results;

        int arrayStart = json.indexOf("[", valueStart);
        int arrayEnd = findMatchingBracket(json, arrayStart);

        if (arrayStart == -1 || arrayEnd == -1)
            return results;

        String arrayContent = json.substring(arrayStart + 1, arrayEnd);

        // Parse individual objects
        int objectStart = 0;
        while ((objectStart = arrayContent.indexOf("{", objectStart)) != -1) {
            int objectEnd = findMatchingBrace(arrayContent, objectStart);
            if (objectEnd == -1)
                break;

            String objectContent = arrayContent.substring(objectStart, objectEnd + 1);

            SearchResult result = new SearchResult();
            result.setId(extractJsonValue(objectContent, "id"));
            result.setContent(extractJsonValue(objectContent, "content"));

            String scoreStr = extractJsonValue(objectContent, "@search.score");
            if (scoreStr != null && !scoreStr.isEmpty()) {
                result.setScore(Float.parseFloat(scoreStr));
            }

            results.add(result);
            objectStart = objectEnd + 1;
        }

        return results;
    }

    private String extractJsonValue(String json, String key) {
        String searchKey = "\"" + key + "\":";
        int keyStart = json.indexOf(searchKey);
        if (keyStart == -1)
            return null;

        int valueStart = keyStart + searchKey.length();
        while (valueStart < json.length() && Character.isWhitespace(json.charAt(valueStart))) {
            valueStart++;
        }

        if (valueStart >= json.length())
            return null;

        char firstChar = json.charAt(valueStart);
        if (firstChar == '"') {
            int valueEnd = json.indexOf("\"", valueStart + 1);
            while (valueEnd > 0 && json.charAt(valueEnd - 1) == '\\') {
                valueEnd = json.indexOf("\"", valueEnd + 1);
            }
            if (valueEnd > valueStart) {
                return json.substring(valueStart + 1, valueEnd);
            }
        } else {
            int valueEnd = valueStart;
            while (valueEnd < json.length() && !",}]".contains(String.valueOf(json.charAt(valueEnd)))) {
                valueEnd++;
            }
            return json.substring(valueStart, valueEnd).trim();
        }

        return null;
    }

    private int findMatchingBracket(String s, int start) {
        int count = 0;
        for (int i = start; i < s.length(); i++) {
            if (s.charAt(i) == '[')
                count++;
            else if (s.charAt(i) == ']') {
                count--;
                if (count == 0)
                    return i;
            }
        }
        return -1;
    }

    private int findMatchingBrace(String s, int start) {
        int count = 0;
        for (int i = start; i < s.length(); i++) {
            if (s.charAt(i) == '{')
                count++;
            else if (s.charAt(i) == '}') {
                count--;
                if (count == 0)
                    return i;
            }
        }
        return -1;
    }

    private String arrayToJson(float[] array) {
        StringBuilder sb = new StringBuilder("[");
        for (int i = 0; i < array.length; i++) {
            if (i > 0)
                sb.append(",");
            sb.append(array[i]);
        }
        sb.append("]");
        return sb.toString();
    }

    private String mapToJson(Map<String, Object> map) {
        StringBuilder sb = new StringBuilder("{");
        sb.append("\"@search.action\": \"mergeOrUpload\",");

        int i = 0;
        for (Map.Entry<String, Object> entry : map.entrySet()) {
            if (i > 0)
                sb.append(",");
            sb.append("\"").append(entry.getKey()).append("\": ");

            Object value = entry.getValue();
            if (value == null) {
                sb.append("null");
            } else if (value instanceof String) {
                sb.append("\"").append(escapeJson((String) value)).append("\"");
            } else if (value instanceof Number) {
                sb.append(value);
            } else if (value instanceof float[]) {
                sb.append(arrayToJson((float[]) value));
            } else {
                sb.append("\"").append(escapeJson(value.toString())).append("\"");
            }
            i++;
        }

        sb.append("}");
        return sb.toString();
    }

    private String escapeJson(String text) {
        return text
                .replace("\\", "\\\\")
                .replace("\"", "\\\"")
                .replace("\n", "\\n")
                .replace("\r", "\\r")
                .replace("\t", "\\t");
    }

    /**
     * Search result.
     */
    public static class SearchResult {
        private String id;
        private String content;
        private float score;
        private Map<String, Object> metadata;

        public String getId() {
            return id;
        }

        public void setId(String id) {
            this.id = id;
        }

        public String getContent() {
            return content;
        }

        public void setContent(String content) {
            this.content = content;
        }

        public float getScore() {
            return score;
        }

        public void setScore(float score) {
            this.score = score;
        }

        public Map<String, Object> getMetadata() {
            return metadata;
        }

        public void setMetadata(Map<String, Object> metadata) {
            this.metadata = metadata;
        }
    }
}
