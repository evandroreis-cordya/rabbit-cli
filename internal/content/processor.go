package content

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ConstructPrompt constructs the prompt for article generation
func ConstructPrompt(url, content, persona, idiom, editorialContent string) string {
	basePrompt := `You are an expert journalist and front-end designer. Write a complete, self-contained HTML file for a journalistic article that looks polished, modern, and minimalist.
The source material for the article is the content found at this URL: ` + url + `

Extracted Content from the URL:
"""
` + content + `
"""
`

	// Add editorial guidelines if provided
	if editorialContent != "" {
		basePrompt += `
Editorial Guidelines:
"""
` + editorialContent + `
"""

`
	}

	basePrompt += `Requirements:
0. The article MUST be written in `

	switch strings.ToUpper(idiom) {
	case "EN":
		basePrompt += "English."
	case "ES":
		basePrompt += "Spanish."
	default:
		basePrompt += "Brazilian Portuguese."
	}

	basePrompt += `
1. The HTML, CSS, and JavaScript must all be included in a single file (no external dependencies or links).
2. Use semantic HTML5 structure (header, main, article, section, footer).
3. The design should follow minimalist principles:
   - Clean typography (prefer sans-serif font like "Inter", "Lato", or "Roboto").
   - Generous white space and consistent line spacing.
   - Subtle color palette (light background, dark text, limited accent colors).
   - Responsive layout that works on desktop and mobile.
4. Include an attention-grabbing article title, byline (author name, date), and compelling lead paragraph.
5. The body of the article should read like a professional news or feature story covering the topic.
6. Use short paragraphs, subheadings, and pull quotes to improve readability.
7. Include subtle interactive or animated elements—such as a fade-in effect for paragraphs or smooth scroll for anchor links—using only lightweight, embedded JavaScript.
8. Add internal CSS styling (using <style> tags) and internal JS (using <script> tags). Do not rely on external sources (CDNs, frameworks, or fonts).
9. End the HTML file cleanly with no incomplete tags.
10. Generate: title, one-liner, article content, captions, and references.
11. IMPORTANT: You MUST generate a detailed English description for a header image that illustrates the article's content. Place this description inside an HTML comment in the following format EXACTLY: <!-- IMAGE_PROMPT: Your detailed image description here -->. Place this comment inside the <article> tag, just before the main content starts. Do NOT include any other <img> tags or placeholders.


Tone/Persona: `

	switch persona {
	case "formal":
		basePrompt += "Formal, press release-type article, written with most strict journalistic guidelines."
	case "personal":
		basePrompt += "Relaxed, pessoal, almost 'what it means to me'-kind of article."
	case "neutra":
		basePrompt += "Neutral style, mimicking articles published in the portal O Tempo (https://www.otempo.com.br/)."
	default:
		basePrompt += "Neutral journalistic style."
	}

	basePrompt += "\n\nOutput only the complete HTML document."
	return basePrompt
}

// GetAvailableSections scans the editoriais directory and returns a list of available section names
func GetAvailableSections(editoriaisPath string) ([]string, error) {
	if editoriaisPath == "" {
		return []string{}, nil
	}

	entries, err := os.ReadDir(editoriaisPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read editoriais directory: %w", err)
	}

	var sections []string
	for _, entry := range entries {
		if entry.IsDir() {
			sections = append(sections, entry.Name())
		}
	}

	return sections, nil
}

// LoadEditorialContent reads the EDITORIA.md file from the specified section folder
// Returns empty string if section is empty, folder doesn't exist, or file doesn't exist (non-fatal)
func LoadEditorialContent(section string, editoriaisPath string) (string, error) {
	if section == "" || editoriaisPath == "" {
		return "", nil
	}

	editorialPath := filepath.Join(editoriaisPath, section, "EDITORIA.md")

	// Check if file exists
	if _, err := os.Stat(editorialPath); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("EDITORIA.md not found in section '%s'", section)
		}
		return "", fmt.Errorf("failed to access EDITORIA.md: %w", err)
	}

	content, err := os.ReadFile(editorialPath)
	if err != nil {
		return "", fmt.Errorf("failed to read EDITORIA.md: %w", err)
	}

	return strings.TrimSpace(string(content)), nil
}
