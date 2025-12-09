package output

import (
	"fmt"
	"os"
	"path/filepath"
)

// SaveImage saves binary image data to a file
// The directory will be created if it doesn't exist
// File permissions are set to 0644
func SaveImage(filePath string, imageData []byte) error {
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Save image data
	if err := os.WriteFile(filePath, imageData, 0644); err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}

	return nil
}

