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

// StableDiffusionGenerator implements ImageGenerator for Stability AI
type StableDiffusionGenerator struct {
	apiKey string
	model  string
}

// NewStableDiffusionGenerator creates a new Stable Diffusion generator
func NewStableDiffusionGenerator(cfg *config.ImageProvider) *StableDiffusionGenerator {
	model := cfg.Model
	if model == "" {
		model = "stable-diffusion-xl-1024-v1-0"
	}
	return &StableDiffusionGenerator{
		apiKey: cfg.APIKey,
		model:  model,
	}
}

// Name returns the generator name
func (s *StableDiffusionGenerator) Name() string {
	return "stable-diffusion"
}

// SupportsAssetType returns whether this generator supports the asset type
func (s *StableDiffusionGenerator) SupportsAssetType(assetType AssetType) bool {
	return assetType == AssetTypeImage
}

// GenerateImage generates an image using Stability AI API
func (s *StableDiffusionGenerator) GenerateImage(prompt string, options ImageOptions) (*Asset, error) {
	if s.apiKey == "" || s.apiKey == "YOUR_STABILITY_API_KEY" {
		return nil, fmt.Errorf("chave de API para Stability AI não definida")
	}

	model := s.model
	if options.Model != "" {
		model = options.Model
	}

	url := fmt.Sprintf("https://api.stability.ai/v1/generation/%s/text-to-image", model)

	reqBody := StableDiffusionRequest{
		TextPrompts: []StableDiffusionPrompt{
			{
				Text: prompt,
			},
		},
		CfgScale: 7,
		Height:   1024,
		Width:    1024,
		Samples:  1,
		Steps:    30,
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
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro da API Stable Diffusion: %s - %s", resp.Status, string(body))
	}

	var sdResp StableDiffusionResponse
	if err := json.NewDecoder(resp.Body).Decode(&sdResp); err != nil {
		return nil, err
	}

	if len(sdResp.Artifacts) == 0 {
		return nil, fmt.Errorf("nenhuma imagem gerada")
	}

	// Decode base64 image
	data, err := base64.StdEncoding.DecodeString(sdResp.Artifacts[0].Base64)
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

// StableDiffusionRequest represents Stability AI API request
type StableDiffusionRequest struct {
	TextPrompts []StableDiffusionPrompt `json:"text_prompts"`
	CfgScale    int                     `json:"cfg_scale"`
	Height      int                     `json:"height"`
	Width       int                     `json:"width"`
	Samples     int                     `json:"samples"`
	Steps       int                     `json:"steps"`
}

// StableDiffusionPrompt represents a prompt in Stability AI request
type StableDiffusionPrompt struct {
	Text string `json:"text"`
}

// StableDiffusionResponse represents Stability AI API response
type StableDiffusionResponse struct {
	Artifacts []StableDiffusionArtifact `json:"artifacts"`
}

// StableDiffusionArtifact represents an artifact in Stability AI response
type StableDiffusionArtifact struct {
	Base64 string `json:"base64"`
}
