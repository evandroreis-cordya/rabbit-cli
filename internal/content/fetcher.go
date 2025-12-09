package content

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

// FetchURLContent fetches and extracts text content from a URL
func FetchURLContent(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,pt-BR;q=0.8,pt;q=0.7")
	req.Header.Set("Referer", "https://www.google.com/")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	content := string(body)

	// Remove scripts and styles
	reScript := regexp.MustCompile(`(?s)<script.*?>.*?</script>`)
	content = reScript.ReplaceAllString(content, "")

	reStyle := regexp.MustCompile(`(?s)<style.*?>.*?</style>`)
	content = reStyle.ReplaceAllString(content, "")

	// Remove HTML tags
	reTags := regexp.MustCompile(`<[^>]*>`)
	content = reTags.ReplaceAllString(content, " ")

	// Collapse whitespace
	reSpace := regexp.MustCompile(`\s+`)
	content = reSpace.ReplaceAllString(content, " ")

	return strings.TrimSpace(content), nil
}
