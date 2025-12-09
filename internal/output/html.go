package output

import (
	"encoding/base64"
	"fmt"
	"html"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// InjectImageIntoHTML injects an image into HTML content
// imagePath should be a relative path to the image file (e.g., "article-title-2025-01-15-1430.jpeg")
func InjectImageIntoHTML(htmlContent string, imagePath string, imagePrompt string) string {
	re := regexp.MustCompile(`<!-- IMAGE_PROMPT: (.*?) -->`)
	matches := re.FindStringSubmatch(htmlContent)

	if len(matches) > 1 {
		// Escape HTML special characters in imagePrompt to prevent XSS
		escapedPrompt := html.EscapeString(imagePrompt)
		// Use just the filename (not full path) for relative reference
		filename := filepath.Base(imagePath)
		imgTag := fmt.Sprintf(`<img src="%s" alt="%s" style="width:100%%; height:auto; border-radius:8px; margin-bottom:20px;">`, filename, escapedPrompt)
		return strings.Replace(htmlContent, matches[0], imgTag, 1)
	}

	return htmlContent
}

// ExtractImagePrompt extracts the image prompt from HTML content
func ExtractImagePrompt(htmlContent string) (string, bool) {
	re := regexp.MustCompile(`<!-- IMAGE_PROMPT: (.*?) -->`)
	matches := re.FindStringSubmatch(htmlContent)

	if len(matches) > 1 {
		return matches[1], true
	}

	return "", false
}

// ExtractTitle extracts the article title from HTML content
// Tries to find <title> tag first, then first <h1> within <article>
func ExtractTitle(htmlContent string) (string, error) {
	// Try to extract from <title> tag
	titleRe := regexp.MustCompile(`(?i)<title[^>]*>(.*?)</title>`)
	matches := titleRe.FindStringSubmatch(htmlContent)
	if len(matches) > 1 {
		title := strings.TrimSpace(matches[1])
		if title != "" {
			return SanitizeFilename(title), nil
		}
	}

	// Try to extract from first <h1> within <article>
	articleRe := regexp.MustCompile(`(?is)<article[^>]*>.*?<h1[^>]*>(.*?)</h1>`)
	matches = articleRe.FindStringSubmatch(htmlContent)
	if len(matches) > 1 {
		title := strings.TrimSpace(matches[1])
		// Remove HTML tags from title
		tagRe := regexp.MustCompile(`<[^>]*>`)
		title = tagRe.ReplaceAllString(title, "")
		if title != "" {
			return SanitizeFilename(title), nil
		}
	}

	// Try to extract from any <h1> tag
	h1Re := regexp.MustCompile(`(?i)<h1[^>]*>(.*?)</h1>`)
	matches = h1Re.FindStringSubmatch(htmlContent)
	if len(matches) > 1 {
		title := strings.TrimSpace(matches[1])
		tagRe := regexp.MustCompile(`<[^>]*>`)
		title = tagRe.ReplaceAllString(title, "")
		if title != "" {
			return SanitizeFilename(title), nil
		}
	}

	return "", fmt.Errorf("no title found in HTML content")
}

// SanitizeFilename sanitizes a title string to be suitable for use in filenames
// Removes/replaces special characters, spaces, and invalid filename characters
// Converts to lowercase, replaces spaces with hyphens, limits length
func SanitizeFilename(title string) string {
	// Convert to lowercase
	title = strings.ToLower(title)

	// Remove HTML entities and decode if needed
	title = html.UnescapeString(title)

	// Replace spaces and multiple spaces with single hyphen
	spaceRe := regexp.MustCompile(`\s+`)
	title = spaceRe.ReplaceAllString(title, "-")

	// Remove or replace invalid filename characters
	var result strings.Builder
	for _, r := range title {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' {
			result.WriteRune(r)
		} else if r == ' ' {
			result.WriteRune('-')
		}
		// Skip other characters
	}

	title = result.String()

	// Remove leading/trailing hyphens and multiple consecutive hyphens
	title = strings.Trim(title, "-")
	hyphenRe := regexp.MustCompile(`-+`)
	title = hyphenRe.ReplaceAllString(title, "-")

	// Limit length to 100 characters
	if len(title) > 100 {
		title = title[:100]
		// Remove trailing hyphen if truncated
		title = strings.TrimSuffix(title, "-")
	}

	// If empty after sanitization, return default
	if title == "" {
		title = "article"
	}

	return title
}

// ExtractImageFromDataURI extracts binary image data and content type from a data URI
func ExtractImageFromDataURI(dataURI string) ([]byte, string, error) {
	// Data URI format: data:image/jpeg;base64,<base64data>
	// or: data:image/png;base64,<base64data>
	parts := strings.SplitN(dataURI, ",", 2)
	if len(parts) != 2 {
		return nil, "", fmt.Errorf("invalid data URI format")
	}

	header := parts[0]
	data := parts[1]

	// Extract content type
	contentType := "image/jpeg" // default
	if strings.HasPrefix(header, "data:") {
		header = strings.TrimPrefix(header, "data:")
		parts := strings.Split(header, ";")
		if len(parts) > 0 {
			contentType = strings.TrimSpace(parts[0])
		}
	}

	// Decode base64
	imageData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, "", fmt.Errorf("failed to decode base64 data: %w", err)
	}

	return imageData, contentType, nil
}
