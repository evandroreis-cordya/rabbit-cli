package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"rabbitai/internal/config"
)

// OpenAIProvider implements the LLMProvider interface for OpenAI
type OpenAIProvider struct {
	apiKey string
	model  string
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(cfg *config.LLMProvider) *OpenAIProvider {
	return &OpenAIProvider{
		apiKey: cfg.APIKey,
		model:  cfg.Model,
	}
}

// Name returns the provider name
func (o *OpenAIProvider) Name() string {
	return "openai"
}

// GenerateText generates text using OpenAI API
func (o *OpenAIProvider) GenerateText(prompt string) (string, error) {
	if o.apiKey == "" || o.apiKey == "YOUR_OPENAI_API_KEY" {
		return "", fmt.Errorf("chave de API para OpenAI não definida")
	}

	url := "https://api.openai.com/v1/chat/completions"

	reqBody := OpenAIRequest{
		Model: o.model,
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
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("erro da API OpenAI: %s - %s", resp.Status, string(body))
	}

	var openAIResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return "", err
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("nenhum conteúdo gerado")
	}

	return cleanMarkdown(openAIResp.Choices[0].Message.Content), nil
}
