package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"rabbitai/internal/config"
)

// GeminiProvider implements the LLMProvider interface for Google Gemini
type GeminiProvider struct {
	apiKey string
	model  string
}

// NewGeminiProvider creates a new Gemini provider
func NewGeminiProvider(cfg *config.LLMProvider) *GeminiProvider {
	return &GeminiProvider{
		apiKey: cfg.APIKey,
		model:  cfg.Model,
	}
}

// Name returns the provider name
func (g *GeminiProvider) Name() string {
	return "gemini"
}

// GenerateText generates text using Gemini API
func (g *GeminiProvider) GenerateText(prompt string) (string, error) {
	if g.apiKey == "" || g.apiKey == "YOUR_GEMINI_API_KEY" {
		return "", fmt.Errorf("chave de API para Gemini não definida")
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", g.model, g.apiKey)

	reqBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Role: "user",
				Parts: []GeminiPart{
					{Text: prompt},
				},
			},
		},
		GenerationConfig: GeminiGenerationConfig{
			Temperature:     0.7,
			TopP:            0.95,
			TopK:            40,
			MaxOutputTokens: 8192,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("erro da API Gemini: %s - %s", resp.Status, string(body))
	}

	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return "", err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("nenhum conteúdo gerado")
	}

	return cleanMarkdown(geminiResp.Candidates[0].Content.Parts[0].Text), nil
}


// GeminiRequest represents the Gemini API request
type GeminiRequest struct {
	Contents         []GeminiContent        `json:"contents"`
	GenerationConfig GeminiGenerationConfig `json:"generationConfig"`
	SafetySettings   []GeminiSafetySetting  `json:"safetySettings"`
}

// GeminiContent represents content in Gemini request
type GeminiContent struct {
	Role  string       `json:"role"`
	Parts []GeminiPart `json:"parts"`
}

// GeminiPart represents a part of Gemini content
type GeminiPart struct {
	Text string `json:"text"`
}

// GeminiGenerationConfig represents generation configuration
type GeminiGenerationConfig struct {
	Temperature     float64 `json:"temperature"`
	TopP            float64 `json:"topP"`
	TopK            int     `json:"topK"`
	MaxOutputTokens int     `json:"maxOutputTokens"`
}

// GeminiSafetySetting represents safety settings
type GeminiSafetySetting struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

// GeminiResponse represents the Gemini API response
type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

// GeminiCandidate represents a candidate in the response
type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}
