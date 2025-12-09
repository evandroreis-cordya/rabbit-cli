package llm

import "strings"

// OpenAIRequest represents the OpenAI-compatible API request (used by OpenAI and Grok)
type OpenAIRequest struct {
	Model            string          `json:"model"`
	Messages         []OpenAIMessage `json:"messages"`
	Temperature      float64         `json:"temperature"`
	TopP             float64         `json:"top_p"`
	MaxTokens        int             `json:"max_tokens"`
	PresencePenalty  float64         `json:"presence_penalty"`
	FrequencyPenalty float64         `json:"frequency_penalty"`
}

// OpenAIMessage represents a message in OpenAI-compatible request
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIResponse represents the OpenAI-compatible API response
type OpenAIResponse struct {
	Choices []OpenAIChoice `json:"choices"`
}

// OpenAIChoice represents a choice in the response
type OpenAIChoice struct {
	Message OpenAIMessage `json:"message"`
}

// cleanMarkdown removes markdown code block markers from text
func cleanMarkdown(text string) string {
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```html")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	return text
}

