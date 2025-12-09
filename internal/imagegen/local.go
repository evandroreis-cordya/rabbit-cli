package imagegen

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"rabbitai/internal/config"
)

// LocalGenerator implements ImageGenerator for local imagebank
type LocalGenerator struct {
	path string
}

// NewLocalGenerator creates a new local storage generator
func NewLocalGenerator(cfg *config.LocalStorageConfig) *LocalGenerator {
	path := cfg.Path
	if path == "" {
		path = "./imagebank"
	}
	return &LocalGenerator{
		path: path,
	}
}

// Name returns the generator name
func (l *LocalGenerator) Name() string {
	return "local"
}

// SupportsAssetType returns whether this generator supports the asset type
func (l *LocalGenerator) SupportsAssetType(assetType AssetType) bool {
	return assetType == AssetTypeLocal || assetType == AssetTypeImage
}

// GenerateImage searches for a matching image in the local imagebank
func (l *LocalGenerator) GenerateImage(prompt string, options ImageOptions) (*Asset, error) {
	// Ensure directory exists
	if err := os.MkdirAll(l.path, 0755); err != nil {
		return nil, fmt.Errorf("falha ao criar diretório imagebank: %w", err)
	}

	// Search for matching images
	imagePath, err := l.findMatchingImage(prompt)
	if err != nil {
		return nil, fmt.Errorf("imagem não encontrada no imagebank: %w", err)
	}

	// Read image file
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler arquivo de imagem: %w", err)
	}

	// Detect content type
	ext := strings.ToLower(filepath.Ext(imagePath))
	var contentType string
	switch ext {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".webp":
		contentType = "image/webp"
	default:
		contentType = "image/jpeg" // Default
	}

	// Encode as base64 data URI
	base64Data := base64.StdEncoding.EncodeToString(data)
	dataURI := fmt.Sprintf("data:%s;base64,%s", contentType, base64Data)

	return &Asset{
		Data:        []byte(dataURI),
		ContentType: contentType,
		Type:        AssetTypeLocal,
		Path:        imagePath,
	}, nil
}

// findMatchingImage searches for an image matching the prompt
func (l *LocalGenerator) findMatchingImage(prompt string) (string, error) {
	// Simple fuzzy matching: look for keywords in filename
	promptLower := strings.ToLower(prompt)
	keywords := strings.Fields(promptLower)

	var bestMatch string
	var bestScore int

	err := filepath.Walk(l.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Check if it's an image file
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
			return nil
		}

		// Score based on keyword matches in filename
		filename := strings.ToLower(filepath.Base(path))
		score := 0
		for _, keyword := range keywords {
			if len(keyword) > 3 && strings.Contains(filename, keyword) {
				score++
			}
		}

		if score > bestScore {
			bestScore = score
			bestMatch = path
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if bestMatch == "" {
		// If no match found, return first image found
		err = filepath.Walk(l.path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				ext := strings.ToLower(filepath.Ext(path))
				if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp" {
					bestMatch = path
					return io.EOF // Stop walking
				}
			}
			return nil
		})
		if err != nil && err != io.EOF {
			return "", err
		}
	}

	if bestMatch == "" {
		return "", fmt.Errorf("nenhuma imagem encontrada no imagebank")
	}

	return bestMatch, nil
}
