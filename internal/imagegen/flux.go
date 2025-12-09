package imagegen

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"rabbitai/internal/config"
)

// FluxGenerator implements ImageGenerator for Black Forest Labs Flux
type FluxGenerator struct {
	apiKey string
	model  string
}

// NewFluxGenerator creates a new Flux generator
func NewFluxGenerator(cfg *config.ImageProvider) *FluxGenerator {
	model := cfg.Model
	if model == "" {
		model = "flux-pro"
	}
	return &FluxGenerator{
		apiKey: cfg.APIKey,
		model:  model,
	}
}

// Name returns the generator name
func (f *FluxGenerator) Name() string {
	return "flux"
}

// SupportsAssetType returns whether this generator supports the asset type
func (f *FluxGenerator) SupportsAssetType(assetType AssetType) bool {
	return assetType == AssetTypeImage
}

// GenerateImage generates an image using Flux API
func (f *FluxGenerator) GenerateImage(prompt string, options ImageOptions) (*Asset, error) {
	if f.apiKey == "" || f.apiKey == "YOUR_FLUX_API_KEY" {
		return nil, fmt.Errorf("chave de API para Flux não definida")
	}

	model := f.model
	if options.Model != "" {
		model = options.Model
	}

	// Using Black Forest Labs API endpoint (example - adjust based on actual API)
	url := "https://api.blackforestlabs.ai/v1/images/generations"

	reqBody := FluxRequest{
		Model:  model,
		Prompt: prompt,
		Width:  1024,
		Height: 1024,
		Steps:  28,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+f.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro da API Flux: %s - %s", resp.Status, string(body))
	}

	var fluxResp FluxResponse
	if err := json.NewDecoder(resp.Body).Decode(&fluxResp); err != nil {
		return nil, err
	}

	if len(fluxResp.Data) == 0 {
		return nil, fmt.Errorf("nenhuma imagem gerada")
	}

	// If URL is provided, download it
	if fluxResp.Data[0].URL != "" {
		imgResp, err := http.Get(fluxResp.Data[0].URL)
		if err != nil {
			return nil, fmt.Errorf("falha ao baixar imagem: %w", err)
		}
		defer imgResp.Body.Close()

		data, err := io.ReadAll(imgResp.Body)
		if err != nil {
			return nil, fmt.Errorf("falha ao ler imagem: %w", err)
		}

		contentType := http.DetectContentType(data)
		base64Data := base64.StdEncoding.EncodeToString(data)
		dataURI := fmt.Sprintf("data:%s;base64,%s", contentType, base64Data)

		return &Asset{
			Data:        []byte(dataURI),
			ContentType: contentType,
			Type:        AssetTypeImage,
			URL:         fluxResp.Data[0].URL,
		}, nil
	}

	// If base64 is provided
	if fluxResp.Data[0].Base64 != "" {
		data, err := base64.StdEncoding.DecodeString(fluxResp.Data[0].Base64)
		if err != nil {
			return nil, fmt.Errorf("falha ao decodificar imagem: %w", err)
		}

		contentType := http.DetectContentType(data)
		base64Data := base64.StdEncoding.EncodeToString(data)
		dataURI := fmt.Sprintf("data:%s;base64,%s", contentType, base64Data)

		return &Asset{
			Data:        []byte(dataURI),
			ContentType: contentType,
			Type:        AssetTypeImage,
		}, nil
	}

	return nil, fmt.Errorf("nenhuma imagem gerada")
}

// FluxRequest represents Flux API request
type FluxRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Steps  int    `json:"steps"`
}

// FluxResponse represents Flux API response
type FluxResponse struct {
	Data []FluxImage `json:"data"`
}

// FluxImage represents an image in Flux response
type FluxImage struct {
	URL    string `json:"url,omitempty"`
	Base64 string `json:"b64_json,omitempty"`
}
