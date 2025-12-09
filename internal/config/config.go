package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	LLM struct {
		Default   string                 `yaml:"default"`
		Providers map[string]LLMProvider `yaml:"providers"`
	} `yaml:"llm"`
	ImageGeneration struct {
		DefaultProvider  string                   `yaml:"default_provider"`
		DefaultAssetType string                   `yaml:"default_asset_type"`
		Providers        map[string]ImageProvider `yaml:"providers"`
		LocalStorage     LocalStorageConfig       `yaml:"local_storage"`
	} `yaml:"image_generation"`
	Editorial EditorialConfig `yaml:"editorial"`
}

// LLMProvider represents LLM provider configuration
type LLMProvider struct {
	APIKey string `yaml:"api_key"`
	Model  string `yaml:"model"`
}

// ImageProvider represents image generation provider configuration
type ImageProvider struct {
	APIKey    string `yaml:"api_key"`
	Model     string `yaml:"model"`
	Size      string `yaml:"size,omitempty"`
	ProjectID string `yaml:"project_id,omitempty"` // For Google Cloud services like Imagen
}

// LocalStorageConfig represents local storage configuration
type LocalStorageConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Path     string `yaml:"path"`
	Fallback bool   `yaml:"fallback"`
}

// EditorialConfig represents editorial section configuration
type EditorialConfig struct {
	DefaultSection  string `yaml:"default_section"`
	EditoriaisPath  string `yaml:"editoriais_path"`
}

var envVarRegex = regexp.MustCompile(`\$\{([^}]+)\}`)

// LoadEnvFile loads environment variables from a .env file.
// If the file doesn't exist, it returns nil (not an error).
// Returns an error only if the file exists but cannot be parsed.
func LoadEnvFile(path string) error {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// File doesn't exist - this is not an error
		return nil
	}

	// File exists, try to load it
	if err := godotenv.Load(path); err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	return nil
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Resolve environment variables
	resolveEnvVars(&config)

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// resolveEnvVars replaces ${VAR_NAME} placeholders with environment variable values
func resolveEnvVars(config *Config) {
	// Resolve LLM provider API keys
	for name, provider := range config.LLM.Providers {
		provider.APIKey = resolveEnvVar(provider.APIKey)
		config.LLM.Providers[name] = provider
	}

	// Resolve image provider API keys
	for name, provider := range config.ImageGeneration.Providers {
		provider.APIKey = resolveEnvVar(provider.APIKey)
		config.ImageGeneration.Providers[name] = provider
	}
}

// resolveEnvVar replaces ${VAR_NAME} with actual environment variable value
func resolveEnvVar(value string) string {
	return envVarRegex.ReplaceAllStringFunc(value, func(match string) string {
		varName := strings.TrimPrefix(strings.TrimSuffix(match, "}"), "${")
		if envValue := os.Getenv(varName); envValue != "" {
			return envValue
		}
		return match // Return original if env var not found
	})
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	if config.LLM.Default == "" {
		return fmt.Errorf("llm.default is required")
	}

	if _, exists := config.LLM.Providers[config.LLM.Default]; !exists {
		return fmt.Errorf("llm.default provider '%s' not found in providers", config.LLM.Default)
	}

	if config.ImageGeneration.DefaultProvider == "" {
		return fmt.Errorf("image_generation.default_provider is required")
	}

	if _, exists := config.ImageGeneration.Providers[config.ImageGeneration.DefaultProvider]; !exists {
		return fmt.Errorf("image_generation.default_provider '%s' not found in providers", config.ImageGeneration.DefaultProvider)
	}

	validAssetTypes := map[string]bool{
		"image": true,
		"video": true,
		"local": true,
	}
	if !validAssetTypes[config.ImageGeneration.DefaultAssetType] {
		return fmt.Errorf("image_generation.default_asset_type must be one of: image, video, local")
	}

	// Validate editorial configuration
	if config.Editorial.EditoriaisPath != "" {
		// Check if path exists (if provided)
		if _, err := os.Stat(config.Editorial.EditoriaisPath); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("editorial.editoriais_path '%s' does not exist", config.Editorial.EditoriaisPath)
			}
		}
	}

	return nil
}

// GetLLMProvider returns the LLM provider configuration
func (c *Config) GetLLMProvider(name string) (*LLMProvider, error) {
	provider, exists := c.LLM.Providers[name]
	if !exists {
		return nil, fmt.Errorf("LLM provider '%s' not found", name)
	}
	return &provider, nil
}

// GetImageProvider returns the image provider configuration
func (c *Config) GetImageProvider(name string) (*ImageProvider, error) {
	provider, exists := c.ImageGeneration.Providers[name]
	if !exists {
		return nil, fmt.Errorf("image provider '%s' not found", name)
	}
	return &provider, nil
}

// GetEditorialPath returns the full path to a section's EDITORIA.md file
func (c *Config) GetEditorialPath(section string) string {
	if section == "" || c.Editorial.EditoriaisPath == "" {
		return ""
	}
	return filepath.Join(c.Editorial.EditoriaisPath, section, "EDITORIA.md")
}
