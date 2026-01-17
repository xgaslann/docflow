package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// AzureAISearch client for Azure Cognitive Search.
type AzureAISearch struct {
	endpoint   string
	apiKey     string
	indexName  string
	apiVersion string
	httpClient *http.Client
}

// AzureSearchConfig holds configuration for Azure AI Search.
type AzureSearchConfig struct {
	Endpoint     string
	APIKey       string
	IndexName    string
	APIVersion   string
	VectorFields string
}

// DefaultAzureSearchConfig returns sensible defaults.
func DefaultAzureSearchConfig() AzureSearchConfig {
	return AzureSearchConfig{
		APIVersion:   "2024-05-01-preview",
		VectorFields: "content_vector",
	}
}

// NewAzureAISearch creates a new Azure AI Search client.
func NewAzureAISearch(config AzureSearchConfig) *AzureAISearch {
	return &AzureAISearch{
		endpoint:   config.Endpoint,
		apiKey:     config.APIKey,
		indexName:  config.IndexName,
		apiVersion: config.APIVersion,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// AzureSearchResult represents a search result.
type AzureSearchResult struct {
	ID       string                 `json:"id"`
	Content  string                 `json:"content"`
	Score    float32                `json:"@search.score"`
	Metadata map[string]interface{} `json:"metadata"`
}

// Search performs a keyword search.
func (c *AzureAISearch) Search(query string, top int) ([]AzureSearchResult, error) {
	return c.HybridSearch(query, nil, top)
}

// VectorSearch performs a vector-only search.
func (c *AzureAISearch) VectorSearch(vector []float32, top int) ([]AzureSearchResult, error) {
	return c.HybridSearch("", vector, top)
}

// HybridSearch performs a hybrid keyword + vector search.
func (c *AzureAISearch) HybridSearch(query string, vector []float32, top int) ([]AzureSearchResult, error) {
	requestBody := make(map[string]interface{})

	if query != "" {
		requestBody["search"] = query
	}

	if len(vector) > 0 {
		requestBody["vectorQueries"] = []map[string]interface{}{
			{
				"kind":   "vector",
				"vector": vector,
				"fields": "content_vector",
				"k":      top,
			},
		}
	}

	requestBody["top"] = top
	requestBody["select"] = "id,content,metadata"

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/indexes/%s/docs/search?api-version=%s",
		c.endpoint, c.indexName, c.apiVersion)

	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search failed: %s", string(body))
	}

	var response struct {
		Value []AzureSearchResult `json:"value"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Value, nil
}

// UploadDocuments uploads documents to the index.
func (c *AzureAISearch) UploadDocuments(documents []map[string]interface{}) error {
	for _, doc := range documents {
		doc["@search.action"] = "mergeOrUpload"
	}

	requestBody := map[string]interface{}{
		"value": documents,
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/indexes/%s/docs/index?api-version=%s",
		c.endpoint, c.indexName, c.apiVersion)

	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 207 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed: %s", string(body))
	}

	return nil
}

// DeleteDocuments deletes documents by ID.
func (c *AzureAISearch) DeleteDocuments(ids []string) error {
	documents := make([]map[string]interface{}, len(ids))
	for i, id := range ids {
		documents[i] = map[string]interface{}{
			"@search.action": "delete",
			"id":             id,
		}
	}

	requestBody := map[string]interface{}{
		"value": documents,
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/indexes/%s/docs/index?api-version=%s",
		c.endpoint, c.indexName, c.apiVersion)

	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 207 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete failed: %s", string(body))
	}

	return nil
}

// SearchWithFilter performs a filtered search.
func (c *AzureAISearch) SearchWithFilter(query string, filter string, top int) ([]AzureSearchResult, error) {
	requestBody := map[string]interface{}{
		"search": query,
		"top":    top,
		"select": "id,content,metadata",
	}

	if filter != "" {
		requestBody["filter"] = filter
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/indexes/%s/docs/search?api-version=%s",
		c.endpoint, c.indexName, c.apiVersion)

	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response struct {
		Value []AzureSearchResult `json:"value"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Value, nil
}

// buildODataFilter constructs an OData filter expression.
func buildODataFilter(filters map[string]interface{}) string {
	if len(filters) == 0 {
		return ""
	}

	conditions := make([]string, 0, len(filters))
	for key, value := range filters {
		switch v := value.(type) {
		case string:
			conditions = append(conditions, fmt.Sprintf("%s eq '%s'", key, v))
		case int, int64, float64:
			conditions = append(conditions, fmt.Sprintf("%s eq %v", key, v))
		case bool:
			conditions = append(conditions, fmt.Sprintf("%s eq %t", key, v))
		}
	}

	return strings.Join(conditions, " and ")
}
