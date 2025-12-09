package output

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GenerateArticlePath generates the full path for an article file
// Format: {editoriaisPath}/{section}/{sanitized-title}-{date}-{hour}.html
// Date format: YYYY-MM-DD-HHMM (e.g., 2025-01-15-1430)
// If section is empty, defaults to "uncategorized"
func GenerateArticlePath(section, title, editoriaisPath string) (string, error) {
	if editoriaisPath == "" {
		return "", fmt.Errorf("editoriais_path is required")
	}

	// Use default section if not provided
	if section == "" {
		section = "uncategorized"
	}

	// Sanitize title if not already sanitized
	sanitizedTitle := SanitizeFilename(title)

	// Generate date-time string: YYYY-MM-DD-HHMM
	now := time.Now()
	dateStr := now.Format("2006-01-02")
	hourStr := now.Format("1504") // HHMM format
	datetimeStr := fmt.Sprintf("%s-%s", dateStr, hourStr)

	// Construct filename
	filename := fmt.Sprintf("%s-%s.html", sanitizedTitle, datetimeStr)

	// Construct full path
	fullPath := filepath.Join(editoriaisPath, section, filename)

	// Ensure section directory exists
	sectionDir := filepath.Join(editoriaisPath, section)
	if err := os.MkdirAll(sectionDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create section directory: %w", err)
	}

	return fullPath, nil
}

// GenerateImagePath generates the path for an image file based on article path and content type
// Returns: {base-path}.{ext} where ext is determined from content type
func GenerateImagePath(articlePath, contentType string) (string, error) {
	if articlePath == "" {
		return "", fmt.Errorf("article path is required")
	}

	// Get base path without extension
	ext := filepath.Ext(articlePath)
	basePath := strings.TrimSuffix(articlePath, ext)

	// Determine image extension from content type
	imageExt := getExtensionFromContentType(contentType)
	if imageExt == "" {
		imageExt = "jpeg" // default fallback
	}

	// Construct image path
	imagePath := fmt.Sprintf("%s.%s", basePath, imageExt)

	return imagePath, nil
}

// getExtensionFromContentType maps MIME content types to file extensions
func getExtensionFromContentType(contentType string) string {
	contentType = strings.ToLower(strings.TrimSpace(contentType))

	// Remove any parameters (e.g., "image/jpeg; charset=utf-8" -> "image/jpeg")
	if idx := strings.Index(contentType, ";"); idx != -1 {
		contentType = contentType[:idx]
	}

	contentType = strings.TrimSpace(contentType)

	// Map common image content types to extensions
	switch contentType {
	case "image/jpeg", "image/jpg":
		return "jpeg"
	case "image/png":
		return "png"
	case "image/gif":
		return "gif"
	case "image/webp":
		return "webp"
	case "image/svg+xml":
		return "svg"
	case "image/bmp":
		return "bmp"
	case "image/tiff":
		return "tiff"
	default:
		// Try to extract extension from content type if it's in format "image/xyz"
		if strings.HasPrefix(contentType, "image/") {
			ext := strings.TrimPrefix(contentType, "image/")
			if ext != "" && !strings.Contains(ext, "/") {
				return ext
			}
		}
		return "jpeg" // default fallback
	}
}

