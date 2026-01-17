package rag

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/xgaslan/docflow/sdks/go/docflow/config"
)

// LLMProcessor handles unified LLM processing for images, tables, and text.
type LLMProcessor struct {
	Config config.LLMConfig
	client *http.Client
}

// NewLLMProcessor creates a new LLM processor.
func NewLLMProcessor(cfg config.LLMConfig) *LLMProcessor {
	timeout := time.Duration(cfg.Timeout) * time.Second
	if timeout == 0 {
		timeout = 60 * time.Second
	}

	return &LLMProcessor{
		Config: cfg,
		client: &http.Client{Timeout: timeout},
	}
}

// ============== Image Processing ==============

// DescribeImage generates a description for an image.
func (p *LLMProcessor) DescribeImage(image ExtractedImage, context string) (string, error) {
	prompt := p.buildImagePrompt(image, context)
	return p.callVisionAPI(image.Data, image.Format, prompt)
}

// AnalyzeImageForRAG performs full RAG analysis on an image.
func (p *LLMProcessor) AnalyzeImageForRAG(image ExtractedImage) (map[string]interface{}, error) {
	prompt := `Analyze this image for RAG (Retrieval-Augmented Generation):

1. **Description**: Detailed description of the image content
2. **Key Information**: Important facts, numbers, or data shown
3. **Entities**: People, organizations, locations, products mentioned
4. **Context**: How this image relates to document content
5. **Data Extraction**: Any text, charts, or tables visible

Respond in JSON format.`

	response, err := p.callVisionAPI(image.Data, image.Format, prompt)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(extractJSON(response)), &result); err != nil {
		// Fallback if JSON parsing fails
		return map[string]interface{}{
			"description":     response,
			"key_information": []string{},
			"entities":        []string{},
		}, nil
	}
	return result, nil
}

// ============== Table Processing ==============

// AnalyzeTable generates analysis/summary for a table.
func (p *LLMProcessor) AnalyzeTable(table ExtractedTable, context string) (string, error) {
	tableMD := p.tableToMarkdown(table)
	prompt := fmt.Sprintf(`Analyze this table and provide:
1. Brief summary of what the table contains
2. Key insights or patterns
3. Important data points

Table:
%s

Context: %s`, tableMD, context)

	return p.callTextAPI(prompt)
}

// ExtractTableData extracts structured data from a table.
func (p *LLMProcessor) ExtractTableData(table ExtractedTable) (map[string]interface{}, error) {
	tableMD := p.tableToMarkdown(table)
	prompt := fmt.Sprintf(`Extract structured information from this table:

%s

Respond with JSON containing:
{
    "summary": "Brief table summary",
    "columns": ["column descriptions"],
    "key_values": {"important": "values"},
    "statistics": {"if applicable": "stats"},
    "trends": ["observed patterns"],
    "entities": ["mentioned entities"]
}`, tableMD)

	response, err := p.callTextAPI(prompt)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(extractJSON(response)), &result); err != nil {
		return map[string]interface{}{"summary": response}, nil
	}
	return result, nil
}

// ============== Text Processing ==============

// GenerateSummary generates a summary of the content.
func (p *LLMProcessor) GenerateSummary(content string, maxLength int) (string, error) {
	if len(content) > 8000 {
		content = content[:8000]
	}
	prompt := fmt.Sprintf(`Summarize the following content in %d characters or less.
Focus on the main points and key information.

Content:
%s`, maxLength, content)

	return p.callTextAPI(prompt)
}

// ExtractKeyPoints extracts key points from content.
func (p *LLMProcessor) ExtractKeyPoints(content string, maxPoints int) ([]string, error) {
	if len(content) > 8000 {
		content = content[:8000]
	}
	prompt := fmt.Sprintf(`Extract the %d most important key points from this content.
Return as a JSON array of strings.

Content:
%s`, maxPoints, content)

	response, err := p.callTextAPI(prompt)
	if err != nil {
		return []string{}, err
	}

	var result []string
	if err := json.Unmarshal([]byte(extractJSON(response)), &result); err != nil {
		// Fallback line splitting
		lines := strings.Split(response, "\n")
		var points []string
		for _, line := range lines {
			line = strings.TrimSpace(strings.TrimLeft(line, "- •0123456789."))
			if line != "" {
				points = append(points, line)
			}
		}
		if len(points) > maxPoints {
			points = points[:maxPoints]
		}
		return points, nil
	}
	return result, nil
}

// ExtractEntities extracts named entities.
func (p *LLMProcessor) ExtractEntities(content string) ([]string, error) {
	if len(content) > 6000 {
		content = content[:6000]
	}
	prompt := `Extract all named entities from this content.
Include: people, organizations, locations, products, dates, numbers.
Return as a JSON array of strings.

Content:
` + content

	response, err := p.callTextAPI(prompt)
	if err != nil {
		return []string{}, err
	}

	var result []string
	if err := json.Unmarshal([]byte(extractJSON(response)), &result); err != nil {
		return []string{}, nil
	}
	return result, nil
}

// EnhanceMetadata enhances document metadata using LLM.
func (p *LLMProcessor) EnhanceMetadata(meta DocumentMetadata, content string) (DocumentMetadata, error) {
	summary, _ := p.GenerateSummary(content, 500)
	meta.Summary = summary

	keyPoints, _ := p.ExtractKeyPoints(content, 5)
	meta.KeyPoints = keyPoints

	entities, _ := p.ExtractEntities(content)
	meta.Entities = entities

	return meta, nil
}

// ============== Internal Helpers ==============

func (p *LLMProcessor) buildImagePrompt(image ExtractedImage, context string) string {
	prompt := "Describe this image in detail for use in a document retrieval system."
	if image.Caption != "" {
		prompt += fmt.Sprintf("\n\nOriginal caption: %s", image.Caption)
	}
	if context != "" {
		prompt += fmt.Sprintf("\n\nSurrounding context: %s", context)
	}
	prompt += "\n\nFocus on: key information, text visible, data shown, and relevance to the document."
	return prompt
}

func (p *LLMProcessor) tableToMarkdown(table ExtractedTable) string {
	var sb strings.Builder
	if len(table.Header) > 0 {
		sb.WriteString("| " + strings.Join(table.Header, " | ") + " |\n")
		separators := make([]string, len(table.Header))
		for i := range separators {
			separators[i] = "---"
		}
		sb.WriteString("| " + strings.Join(separators, " | ") + " |\n")
	}
	for _, row := range table.Rows {
		sb.WriteString("| " + strings.Join(row, " | ") + " |\n")
	}
	return sb.String()
}

func extractJSON(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start == -1 || end == -1 {
		start = strings.Index(s, "[")
		end = strings.LastIndex(s, "]")
	}
	if start != -1 && end != -1 && end > start {
		return s[start : end+1]
	}
	return s
}

// ============== API Calls ==============

func (p *LLMProcessor) callVisionAPI(data []byte, format, prompt string) (string, error) {
	switch p.Config.Provider {
	case "openai":
		return p.callOpenAIVision(data, format, prompt)
	case "anthropic":
		return p.callAnthropicVision(data, format, prompt)
	case "ollama":
		return p.callOllamaVision(data, prompt)
	default:
		return "", fmt.Errorf("unsupported provider: %s", p.Config.Provider)
	}
}

func (p *LLMProcessor) callTextAPI(prompt string) (string, error) {
	switch p.Config.Provider {
	case "openai":
		return p.callOpenAIText(prompt)
	case "anthropic":
		return p.callAnthropicText(prompt)
	case "ollama":
		return p.callOllamaText(prompt)
	default:
		return "", fmt.Errorf("unsupported provider: %s", p.Config.Provider)
	}
}

// OpenAI Implementation
func (p *LLMProcessor) callOpenAIVision(data []byte, format, prompt string) (string, error) {
	b64Image := base64.StdEncoding.EncodeToString(data)
	mediaType := fmt.Sprintf("image/%s", format)
	detail := "auto"
	if p.Config.Detail != "" {
		detail = p.Config.Detail
	}

	payload := map[string]interface{}{
		"model": p.Config.Model,
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{"type": "text", "text": prompt},
					{
						"type": "image_url",
						"image_url": map[string]interface{}{
							"url":    fmt.Sprintf("data:%s;base64,%s", mediaType, b64Image),
							"detail": detail,
						},
					},
				},
			},
		},
		"max_tokens":  p.Config.MaxTokens,
		"temperature": p.Config.Temperature,
	}

	return p.makeRequest("POST", "https://api.openai.com/v1/chat/completions", payload, p.Config.APIKey)
}

func (p *LLMProcessor) callOpenAIText(prompt string) (string, error) {
	model := p.Config.Model
	if strings.Contains(model, "vision") {
		model = "gpt-4"
	}

	payload := map[string]interface{}{
		"model":       model,
		"messages":    []map[string]interface{}{{"role": "user", "content": prompt}},
		"max_tokens":  p.Config.MaxTokens,
		"temperature": p.Config.Temperature,
	}

	return p.makeRequest("POST", "https://api.openai.com/v1/chat/completions", payload, p.Config.APIKey)
}

// Anthropic Implementation
func (p *LLMProcessor) callAnthropicVision(data []byte, format, prompt string) (string, error) {
	b64Image := base64.StdEncoding.EncodeToString(data)
	payload := map[string]interface{}{
		"model":      p.Config.Model,
		"max_tokens": p.Config.MaxTokens,
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type": "image",
						"source": map[string]interface{}{
							"type":       "base64",
							"media_type": "image/" + format,
							"data":       b64Image,
						},
					},
					{"type": "text", "text": prompt},
				},
			},
		},
	}
	return p.makeRequestUnwrapContent("POST", "https://api.anthropic.com/v1/messages", payload, p.Config.APIKey, "anthropic")
}

func (p *LLMProcessor) callAnthropicText(prompt string) (string, error) {
	model := strings.Replace(p.Config.Model, "-vision", "", 1)
	payload := map[string]interface{}{
		"model":      model,
		"max_tokens": p.Config.MaxTokens,
		"messages":   []map[string]interface{}{{"role": "user", "content": prompt}},
	}
	return p.makeRequestUnwrapContent("POST", "https://api.anthropic.com/v1/messages", payload, p.Config.APIKey, "anthropic")
}

// Ollama Implementation
func (p *LLMProcessor) callOllamaVision(data []byte, prompt string) (string, error) {
	b64Image := base64.StdEncoding.EncodeToString(data)
	baseURL := "http://localhost:11434"
	if p.Config.BaseURL != "" {
		baseURL = p.Config.BaseURL
	}

	payload := map[string]interface{}{
		"model":  p.Config.Model,
		"prompt": prompt,
		"images": []string{b64Image},
		"stream": false,
	}
	return p.makeRequestOllama(baseURL+"/api/generate", payload)
}

func (p *LLMProcessor) callOllamaText(prompt string) (string, error) {
	baseURL := "http://localhost:11434"
	if p.Config.BaseURL != "" {
		baseURL = p.Config.BaseURL
	}

	payload := map[string]interface{}{
		"model":  p.Config.Model,
		"prompt": prompt,
		"stream": false,
	}
	return p.makeRequestOllama(baseURL+"/api/generate", payload)
}

// Generic Request Helpers
func (p *LLMProcessor) makeRequest(method, url string, payload interface{}, apiKey string) (string, error) {
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(method, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API error: %s", string(respBody))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &result); err == nil && len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("failed to parse response: %s", string(respBody))
}

func (p *LLMProcessor) makeRequestUnwrapContent(method, url string, payload interface{}, apiKey, typeStr string) (string, error) {
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(method, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if typeStr == "anthropic" {
		req.Header.Set("x-api-key", apiKey)
		req.Header.Set("anthropic-version", "2023-06-01")
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API error: %s", string(respBody))
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(respBody, &result); err == nil && len(result.Content) > 0 {
		return result.Content[0].Text, nil
	}
	return "", fmt.Errorf("failed to parse response: %s", string(respBody))
}

func (p *LLMProcessor) makeRequestOllama(url string, payload interface{}) (string, error) {
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Ollama error: %s", string(respBody))
	}

	var result struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(respBody, &result); err == nil {
		return result.Response, nil
	}
	return "", fmt.Errorf("failed to parse response: %s", string(respBody))
}
