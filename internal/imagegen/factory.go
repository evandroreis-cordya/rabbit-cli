package imagegen

import (
	"fmt"

	"rabbitai/internal/config"
)

// NewGenerator creates a new image generator based on the name
func NewGenerator(name string, cfg *config.Config) (ImageGenerator, error) {
	switch name {
	case "dalle":
		providerCfg, err := cfg.GetImageProvider("dalle")
		if err != nil {
			return nil, err
		}
		return NewDalleGenerator(providerCfg), nil
	case "imagen":
		providerCfg, err := cfg.GetImageProvider("imagen")
		if err != nil {
			return nil, err
		}
		return NewImagenGenerator(providerCfg), nil
	case "stable-diffusion":
		providerCfg, err := cfg.GetImageProvider("stable_diffusion")
		if err != nil {
			return nil, err
		}
		return NewStableDiffusionGenerator(providerCfg), nil
	case "flux":
		providerCfg, err := cfg.GetImageProvider("flux")
		if err != nil {
			return nil, err
		}
		return NewFluxGenerator(providerCfg), nil
	case "local":
		return NewLocalGenerator(&cfg.ImageGeneration.LocalStorage), nil
	default:
		return nil, fmt.Errorf("unsupported image generator: %s", name)
	}
}
