package llm

import (
	"fmt"

	"rabbitai/internal/config"
)

// NewProvider creates a new LLM provider based on the name
func NewProvider(name string, cfg *config.Config) (LLMProvider, error) {
	providerCfg, err := cfg.GetLLMProvider(name)
	if err != nil {
		return nil, err
	}

	switch name {
	case "gemini":
		return NewGeminiProvider(providerCfg), nil
	case "openai":
		return NewOpenAIProvider(providerCfg), nil
	case "claude":
		return NewClaudeProvider(providerCfg), nil
	case "grok":
		return NewGrokProvider(providerCfg), nil
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", name)
	}
}
