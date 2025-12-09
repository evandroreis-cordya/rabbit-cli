package llm

// LLMProvider defines the interface for LLM providers
type LLMProvider interface {
	// GenerateText generates text from a prompt
	GenerateText(prompt string) (string, error)

	// Name returns the name of the provider
	Name() string
}
