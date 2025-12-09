package imagegen

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"rabbitai/internal/config"
)

// ImagenGenerator implements ImageGenerator for Google Imagen 3
type ImagenGenerator struct {
	apiKey    string
	model     string
	size      string
	projectID string
}

// NewImagenGenerator creates a new Imagen generator
func NewImagenGenerator(cfg *config.ImageProvider) *ImagenGenerator {
	size := cfg.Size
	if size == "" {
		size = "1024x1024"
	}
	projectID := cfg.ProjectID
	if projectID == "" {
		// Try to get from environment variable as fallback
		projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
		if projectID == "" {
			projectID = os.Getenv("GCP_PROJECT_ID")
		}
	}
	return &ImagenGenerator{
		apiKey:    cfg.APIKey,
		model:     cfg.Model,
		size:      size,
		projectID: projectID,
	}
}

// Name returns the generator name
func (i *ImagenGenerator) Name() string {
	return "imagen"
}

// SupportsAssetType returns whether this generator supports the asset type
func (i *ImagenGenerator) SupportsAssetType(assetType AssetType) bool {
	return assetType == AssetTypeImage
}

// GenerateImage generates an image using Imagen 3 API
func (i *ImagenGenerator) GenerateImage(prompt string, options ImageOptions) (*Asset, error) {
	if i.apiKey == "" || i.apiKey == "YOUR_GOOGLE_API_KEY" {
		return nil, fmt.Errorf("chave de API para Google não definida")
	}

	if i.projectID == "" {
		return nil, fmt.Errorf("project ID não definido. Configure 'project_id' no config.yaml ou defina GOOGLE_CLOUD_PROJECT")
	}

	// Using Vertex AI Imagen API endpoint
	url := fmt.Sprintf("https://us-central1-aiplatform.googleapis.com/v1/projects/%s/locations/us-central1/publishers/google/models/%s:predict", i.projectID, i.model)

	size := i.size
	if options.Size != "" {
		size = options.Size
	}

	// Imagen API request structure
	reqBody := ImagenRequest{
		Instances: []ImagenInstance{
			{
				Prompt: prompt,
			},
		},
		Parameters: ImagenParameters{
			SampleCount:    1,
			Size:           size,
			GuidanceScale:  7.5,
			NegativePrompt: "",
		},
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
	req.Header.Set("Authorization", "Bearer "+i.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro da API Imagen: %s - %s", resp.Status, string(body))
	}

	var imagenResp ImagenResponse
	if err := json.NewDecoder(resp.Body).Decode(&imagenResp); err != nil {
		return nil, err
	}

	if len(imagenResp.Predictions) == 0 || imagenResp.Predictions[0].BytesBase64Encoded == "" {
		return nil, fmt.Errorf("nenhuma imagem gerada")
	}

	// Decode base64 image
	data, err := base64.StdEncoding.DecodeString(imagenResp.Predictions[0].BytesBase64Encoded)
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

// ImagenRequest represents Imagen API request
type ImagenRequest struct {
	Instances  []ImagenInstance `json:"instances"`
	Parameters ImagenParameters `json:"parameters"`
}

// ImagenInstance represents an instance in Imagen request
type ImagenInstance struct {
	Prompt string `json:"prompt"`
}

// ImagenParameters represents parameters for Imagen
type ImagenParameters struct {
	SampleCount    int     `json:"sampleCount"`
	Size           string  `json:"size"`
	GuidanceScale  float64 `json:"guidanceScale"`
	NegativePrompt string  `json:"negativePrompt"`
}

// ImagenResponse represents Imagen API response
type ImagenResponse struct {
	Predictions []ImagenPrediction `json:"predictions"`
}

// ImagenPrediction represents a prediction in Imagen response
type ImagenPrediction struct {
	BytesBase64Encoded string `json:"bytesBase64Encoded"`
}
