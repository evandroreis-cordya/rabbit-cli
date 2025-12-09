package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"rabbitai/internal/config"
)

// ClaudeProvider implements the LLMProvider interface for Anthropic Claude
type ClaudeProvider struct {
	apiKey string
	model  string
}

// NewClaudeProvider creates a new Claude provider
func NewClaudeProvider(cfg *config.LLMProvider) *ClaudeProvider {
	return &ClaudeProvider{
		apiKey: cfg.APIKey,
		model:  cfg.Model,
	}
}

// Name returns the provider name
func (c *ClaudeProvider) Name() string {
	return "claude"
}

// GenerateText generates text using Claude API
func (c *ClaudeProvider) GenerateText(prompt string) (string, error) {
	if c.apiKey == "" || c.apiKey == "YOUR_CLAUDE_API_KEY" {
		return "", fmt.Errorf("chave de API para Claude não definida")
	}

	url := "https://api.anthropic.com/v1/messages"

	reqBody := ClaudeRequest{
		Model:       c.model,
		MaxTokens:   8192,
		Temperature: 0.7,
		TopP:        1.0,
		System:      "You are a helpful assistant.",
		Messages: []ClaudeMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("erro da API Claude: %s - %s", resp.Status, string(body))
	}

	var claudeResp ClaudeResponse
	if err := json.NewDecoder(resp.Body).Decode(&claudeResp); err != nil {
		return "", err
	}

	if len(claudeResp.Content) == 0 {
		return "", fmt.Errorf("nenhum conteúdo gerado")
	}

	return cleanMarkdown(claudeResp.Content[0].Text), nil
}

// ClaudeRequest represents the Claude API request
type ClaudeRequest struct {
	Model       string          `json:"model"`
	MaxTokens   int             `json:"max_tokens"`
	Temperature float64         `json:"temperature"`
	TopP        float64         `json:"top_p"`
	System      string          `json:"system"`
	Messages    []ClaudeMessage `json:"messages"`
}

// ClaudeMessage represents a message in Claude request
type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ClaudeResponse represents the Claude API response
type ClaudeResponse struct {
	Content []ClaudeContent `json:"content"`
}

// ClaudeContent represents content in Claude response
type ClaudeContent struct {
	Text string `json:"text"`
}
