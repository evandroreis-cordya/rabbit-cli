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

// DalleGenerator implements ImageGenerator for DALL-E
type DalleGenerator struct {
	apiKey string
	model  string
	size   string
}

// NewDalleGenerator creates a new DALL-E generator
func NewDalleGenerator(cfg *config.ImageProvider) *DalleGenerator {
	size := cfg.Size
	if size == "" {
		size = "1024x1024"
	}
	return &DalleGenerator{
		apiKey: cfg.APIKey,
		model:  cfg.Model,
		size:   size,
	}
}

// Name returns the generator name
func (d *DalleGenerator) Name() string {
	return "dalle"
}

// SupportsAssetType returns whether this generator supports the asset type
func (d *DalleGenerator) SupportsAssetType(assetType AssetType) bool {
	return assetType == AssetTypeImage
}

// GenerateImage generates an image using DALL-E API
func (d *DalleGenerator) GenerateImage(prompt string, options ImageOptions) (*Asset, error) {
	if d.apiKey == "" || d.apiKey == "YOUR_OPENAI_API_KEY" {
		return nil, fmt.Errorf("chave de API para OpenAI não definida")
	}

	url := "https://api.openai.com/v1/images/generations"

	size := d.size
	if options.Size != "" {
		size = options.Size
	}

	model := d.model
	if options.Model != "" {
		model = options.Model
	}

	reqBody := DalleRequest{
		Model:  model,
		Prompt: prompt,
		Size:   size,
		N:      1,
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
	req.Header.Set("Authorization", "Bearer "+d.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro da API DALL-E: %s - %s", resp.Status, string(body))
	}

	var dalleResp DalleResponse
	if err := json.NewDecoder(resp.Body).Decode(&dalleResp); err != nil {
		return nil, err
	}

	if len(dalleResp.Data) == 0 {
		return nil, fmt.Errorf("nenhuma imagem gerada")
	}

	imageURL := dalleResp.Data[0].Url

	// Download and encode image
	imgResp, err := http.Get(imageURL)
	if err != nil {
		return nil, fmt.Errorf("falha ao baixar imagem: %w", err)
	}
	defer imgResp.Body.Close()

	if imgResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("falha ao baixar imagem: %s", imgResp.Status)
	}

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
		URL:         imageURL,
	}, nil
}

// DalleRequest represents DALL-E API request
type DalleRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Size   string `json:"size"`
	N      int    `json:"n"`
}

// DalleResponse represents DALL-E API response
type DalleResponse struct {
	Data []DalleImage `json:"data"`
}

// DalleImage represents an image in DALL-E response
type DalleImage struct {
	Url string `json:"url"`
}
