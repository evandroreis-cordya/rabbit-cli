package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"rabbitai/internal/config"
)

// GrokProvider implements the LLMProvider interface for xAI Grok
type GrokProvider struct {
	apiKey string
	model  string
}

// NewGrokProvider creates a new Grok provider
func NewGrokProvider(cfg *config.LLMProvider) *GrokProvider {
	return &GrokProvider{
		apiKey: cfg.APIKey,
		model:  cfg.Model,
	}
}

// Name returns the provider name
func (g *GrokProvider) Name() string {
	return "grok"
}

// GenerateText generates text using Grok API
func (g *GrokProvider) GenerateText(prompt string) (string, error) {
	if g.apiKey == "" || g.apiKey == "YOUR_GROK_API_KEY" {
		return "", fmt.Errorf("chave de API para Grok não definida")
	}

	url := "https://api.x.ai/v1/chat/completions"

	reqBody := OpenAIRequest{
		Model: g.model,
		Messages: []OpenAIMessage{
			{Role: "user", Content: prompt},
		},
		Temperature:      0.7,
		TopP:             1.0,
		MaxTokens:        4096,
		PresencePenalty:  0.0,
		FrequencyPenalty: 0.0,
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
	req.Header.Set("Authorization", "Bearer "+g.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("erro da API Grok: %s - %s", resp.Status, string(body))
	}

	var grokResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&grokResp); err != nil {
		return "", err
	}

	if len(grokResp.Choices) == 0 {
		return "", fmt.Errorf("nenhum conteúdo gerado")
	}

	return cleanMarkdown(grokResp.Choices[0].Message.Content), nil
}
