package imagegen

// AssetType represents the type of asset
type AssetType string

const (
	AssetTypeImage AssetType = "image"
	AssetTypeVideo AssetType = "video"
	AssetTypeLocal AssetType = "local"
)

// Asset represents a generated or retrieved asset
type Asset struct {
	Data        []byte
	ContentType string
	Type        AssetType
	Path        string // For local assets
	URL         string // For remote assets
}

// ImageOptions represents options for image generation
type ImageOptions struct {
	Size   string
	Model  string
	Format string
}

// ImageGenerator defines the interface for image generation providers
type ImageGenerator interface {
	// GenerateImage generates an image from a prompt
	GenerateImage(prompt string, options ImageOptions) (*Asset, error)

	// SupportsAssetType returns whether this generator supports the given asset type
	SupportsAssetType(assetType AssetType) bool

	// Name returns the name of the generator
	Name() string
}
