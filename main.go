package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// --- Configuration ---

// API Keys - Loaded from environment variables or config
var (
	GeminiAPIKey = os.Getenv("GEMINI_API_KEY")
	OpenAIAPIKey = os.Getenv("OPENAI_API_KEY")
	ClaudeAPIKey = os.Getenv("CLAUDE_API_KEY")
	GrokAPIKey   = os.Getenv("GROK_API_KEY")
)

// Default Configuration
const (
	DefaultLLM     = "gemini"
	DefaultPersona = "neutra"
	DefaultIdiom   = "BR"
)

// --- Structs ---

// Gemini Request/Response Structs
type GeminiRequest struct {
	Contents         []GeminiContent        `json:"contents"`
	GenerationConfig GeminiGenerationConfig `json:"generationConfig"`
	SafetySettings   []GeminiSafetySetting  `json:"safetySettings"`
}

type GeminiContent struct {
	Role  string       `json:"role"`
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiGenerationConfig struct {
	Temperature     float64 `json:"temperature"`
	TopP            float64 `json:"topP"`
	TopK            int     `json:"topK"`
	MaxOutputTokens int     `json:"maxOutputTokens"`
}

type GeminiSafetySetting struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}

// OpenAI/Grok Request/Response Structs (Compatible)
type OpenAIRequest struct {
	Model            string          `json:"model"`
	Messages         []OpenAIMessage `json:"messages"`
	Temperature      float64         `json:"temperature"`
	TopP             float64         `json:"top_p"`
	MaxTokens        int             `json:"max_tokens"`
	PresencePenalty  float64         `json:"presence_penalty"`
	FrequencyPenalty float64         `json:"frequency_penalty"`
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []OpenAIChoice `json:"choices"`
}

type OpenAIChoice struct {
	Message OpenAIMessage `json:"message"`
}

// Claude Request/Response Structs
type ClaudeRequest struct {
	Model       string          `json:"model"`
	MaxTokens   int             `json:"max_tokens"`
	Temperature float64         `json:"temperature"`
	TopP        float64         `json:"top_p"`
	System      string          `json:"system"`
	Messages    []ClaudeMessage `json:"messages"`
}

type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClaudeResponse struct {
	Content []ClaudeContent `json:"content"`
}

type ClaudeContent struct {
	Text string `json:"text"`
}

// DALL-E 3 Request/Response Structs
type DalleRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Size   string `json:"size"`
	N      int    `json:"n"`
}

type DalleResponse struct {
	Data []DalleImage `json:"data"`
}

type DalleImage struct {
	Url string `json:"url"`
}

// --- Main ---

func main() {
	// 1. Print Banner
	printBanner()

	// 2. Parse Flags
	urlFlag := flag.String("u", "", "URL fonte para o artigo (atalho)")
	urlFlagLong := flag.String("url", "", "URL fonte para o artigo")
	llmFlag := flag.String("l", DefaultLLM, "LLM para usar: gemini (default), openai, claude, grok")
	llmFlagLong := flag.String("llm", DefaultLLM, "LLM para usar: gemini, openai, claude, grok")
	personaFlag := flag.String("p", DefaultPersona, "Tom da persona: formal, pessoal, neutra (atalho)")
	personaFlagLong := flag.String("persona", DefaultPersona, "Tom da persona: formal, pessoal, neutra")
	idiomFlag := flag.String("i", DefaultIdiom, "Idioma do artigo: BR, EN, ES (atalho)")
	idiomFlagLong := flag.String("idiom", DefaultIdiom, "Idioma do artigo: BR, EN, ES")

	flag.Parse()

	// Resolve flags (long takes precedence if set, otherwise short, otherwise default is handled by flag default)
	// Actually, standard flag behavior is separate. Let's just check which one is set or use the value.
	// A common pattern for aliasing is checking if the non-default value is present.
	// Since we can't easily distinguish "not set" from "set to default" without more logic,
	// let's simplify: check if long is set, if not use short.

	targetURL := *urlFlag
	if *urlFlagLong != "" {
		targetURL = *urlFlagLong
	}

	targetLLM := *llmFlag
	if *llmFlagLong != DefaultLLM { // If long flag is explicitly set (assuming user doesn't set it to default manually, which is fine)
		targetLLM = *llmFlagLong
	}

	targetPersona := *personaFlag
	if *personaFlagLong != DefaultPersona {
		targetPersona = *personaFlagLong
	}

	targetIdiom := *idiomFlag
	if *idiomFlagLong != DefaultIdiom {
		targetIdiom = *idiomFlagLong
	}

	// 3. Validate Parameters
	if targetURL == "" {
		printUsage()
		os.Exit(0)
	}

	fmt.Printf("URL Alvo: %s\n", targetURL)
	fmt.Printf("LLM: %s\n", targetLLM)
	fmt.Printf("Persona: %s\n", targetPersona)
	fmt.Printf("Idioma: %s\n", targetIdiom)
	fmt.Println("--------------------------------------------------------")

	// 4. Core Logic
	fmt.Println("Gerando artigo... por favor aguarde.")

	// 4.1 Fetch Content
	fmt.Printf("Baixando conteúdo de %s...\n", targetURL)
	extractedContent, fetchErr := fetchURLContent(targetURL)
	if fetchErr != nil {
		fmt.Printf("Aviso: Falha ao baixar conteúdo da URL: %v\n", fetchErr)
		fmt.Println("Tentando gerar apenas com a URL...")
		extractedContent = ""
	} else {
		fmt.Printf("Conteúdo extraído: %d caracteres.\n", len(extractedContent))
	}

	// Construct Prompt
	prompt := constructPrompt(targetURL, extractedContent, targetPersona, targetIdiom)

	var htmlContent string
	var err error

	switch strings.ToLower(targetLLM) {
	case "gemini":
		htmlContent, err = callGemini(prompt)
	case "openai":
		htmlContent, err = callOpenAI(prompt)
	case "claude":
		htmlContent, err = callClaude(prompt)
	case "grok":
		htmlContent, err = callGrok(prompt)
	default:
		exitWithError(fmt.Errorf("LLM não suportado: %s", targetLLM))
	}

	if err != nil {
		exitWithError(err)
	}

	if err != nil {
		exitWithError(err)
	}

	// 4.5 Image Generation
	fmt.Println("Analisando necessidade de imagem...")
	re := regexp.MustCompile(`<!-- IMAGE_PROMPT: (.*?) -->`)
	matches := re.FindStringSubmatch(htmlContent)

	if len(matches) > 1 {
		imagePrompt := matches[1]
		fmt.Printf("Gerando imagem com DALL-E 3: %s\n", imagePrompt)
		imageURL, err := generateImage(imagePrompt)
		if err != nil {
			fmt.Printf("Aviso: Falha ao gerar imagem: %v\n", err)
		} else {
			fmt.Println("Imagem gerada. Baixando e codificando...")
			base64Image, err := downloadAndEncodeImage(imageURL)
			if err != nil {
				fmt.Printf("Aviso: Falha ao processar imagem: %v\n", err)
			} else {
				// Inject image into HTML
				imgTag := fmt.Sprintf(`<img src="%s" alt="%s" style="width:100%%; height:auto; border-radius:8px; margin-bottom:20px;">`, base64Image, imagePrompt)
				// Replace a placeholder or insert at the beginning of the article
				// We'll look for the first <article> tag and insert after it, or replace a specific placeholder if we asked for one.
				// Let's replace the comment with the image tag.
				htmlContent = strings.Replace(htmlContent, matches[0], imgTag, 1)
			}
		}
	} else {
		fmt.Println("Nenhuma descrição de imagem encontrada.")
	}

	// 5. Save and Open
	filename := "article.html"
	err = os.WriteFile(filename, []byte(htmlContent), 0644)
	if err != nil {
		exitWithError(fmt.Errorf("falha ao salvar arquivo: %v", err))
	}

	fmt.Printf("Artigo salvo em %s\n", filename)
	openBrowser(filename)

	fmt.Println("Mais uma missão cumprida por Humanos e IA. Bem-vindo à era dos AI Citizens.")
}

// --- Helper Functions ---

func constructPrompt(url, content, persona, idiom string) string {
	basePrompt := `You are an expert journalist and front-end designer. Write a complete, self-contained HTML file for a journalistic article that looks polished, modern, and minimalist.
The source material for the article is the content found at this URL: ` + url + `

Extracted Content from the URL:
"""
` + content + `
"""

Requirements:
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
   - Clean typography (prefer sans-serif font like “Inter”, “Lato”, or “Roboto”).
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

func callGemini(prompt string) (string, error) {
	if GeminiAPIKey == "" || GeminiAPIKey == "YOUR_GEMINI_API_KEY" {
		return "", fmt.Errorf("chave de API para Gemini não definida")
	}

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash-exp:generateContent?key=" + GeminiAPIKey

	reqBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Role: "user",
				Parts: []GeminiPart{
					{Text: prompt},
				},
			},
		},
		GenerationConfig: GeminiGenerationConfig{
			Temperature:     0.7,
			TopP:            0.95,
			TopK:            40,
			MaxOutputTokens: 8192,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("erro da API Gemini: %s - %s", resp.Status, string(body))
	}

	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return "", err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("nenhum conteúdo gerado")
	}

	return cleanMarkdown(geminiResp.Candidates[0].Content.Parts[0].Text), nil
}

func callOpenAI(prompt string) (string, error) {
	if OpenAIAPIKey == "" || OpenAIAPIKey == "YOUR_OPENAI_API_KEY" {
		return "", fmt.Errorf("chave de API para OpenAI não definida")
	}

	url := "https://api.openai.com/v1/chat/completions"

	reqBody := OpenAIRequest{
		Model: "gpt-4o", // Using a modern model as default
		Messages: []OpenAIMessage{
			{Role: "user", Content: prompt},
		},
		Temperature:      0.7,
		TopP:             1.0,
		MaxTokens:        4096,
		PresencePenalty:  0.0,
		FrequencyPenalty: 0.0,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+OpenAIAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("erro da API OpenAI: %s - %s", resp.Status, string(body))
	}

	var openAIResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return "", err
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("nenhum conteúdo gerado")
	}

	return cleanMarkdown(openAIResp.Choices[0].Message.Content), nil
}

func callClaude(prompt string) (string, error) {
	if ClaudeAPIKey == "" || ClaudeAPIKey == "YOUR_CLAUDE_API_KEY" {
		return "", fmt.Errorf("chave de API para Claude não definida")
	}

	url := "https://api.anthropic.com/v1/messages"

	reqBody := ClaudeRequest{
		Model:       "claude-3-5-sonnet-20241022",
		MaxTokens:   8192,
		Temperature: 0.7,
		TopP:        1.0,
		System:      "You are a helpful assistant.",
		Messages: []ClaudeMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", ClaudeAPIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("erro da API Claude: %s - %s", resp.Status, string(body))
	}

	var claudeResp ClaudeResponse
	if err := json.NewDecoder(resp.Body).Decode(&claudeResp); err != nil {
		return "", err
	}

	if len(claudeResp.Content) == 0 {
		return "", fmt.Errorf("nenhum conteúdo gerado")
	}

	return cleanMarkdown(claudeResp.Content[0].Text), nil
}

func callGrok(prompt string) (string, error) {
	if GrokAPIKey == "" || GrokAPIKey == "YOUR_GROK_API_KEY" {
		return "", fmt.Errorf("chave de API para Grok não definida")
	}

	url := "https://api.x.ai/v1/chat/completions"

	// Grok uses OpenAI-compatible API
	reqBody := OpenAIRequest{
		Model: "grok-3", // Or grok-1, depending on availability
		Messages: []OpenAIMessage{
			{Role: "user", Content: prompt},
		},
		Temperature:      0.7,
		TopP:             1.0,
		MaxTokens:        4096,
		PresencePenalty:  0.0,
		FrequencyPenalty: 0.0,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+GrokAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("erro da API Grok: %s - %s", resp.Status, string(body))
	}

	var grokResp OpenAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&grokResp); err != nil {
		return "", err
	}

	if len(grokResp.Choices) == 0 {
		return "", fmt.Errorf("nenhum conteúdo gerado")
	}

	return cleanMarkdown(grokResp.Choices[0].Message.Content), nil
}

func generateImage(prompt string) (string, error) {
	if OpenAIAPIKey == "" || OpenAIAPIKey == "YOUR_OPENAI_API_KEY" {
		return "", fmt.Errorf("chave de API para OpenAI não definida")
	}

	url := "https://api.openai.com/v1/images/generations"

	reqBody := DalleRequest{
		Model:  "dall-e-3",
		Prompt: prompt,
		Size:   "1024x1024",
		N:      1,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+OpenAIAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("erro da API DALL-E: %s - %s", resp.Status, string(body))
	}

	var dalleResp DalleResponse
	if err := json.NewDecoder(resp.Body).Decode(&dalleResp); err != nil {
		return "", err
	}

	if len(dalleResp.Data) == 0 {
		return "", fmt.Errorf("nenhuma imagem gerada")
	}

	return dalleResp.Data[0].Url, nil
}

func downloadAndEncodeImage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("falha ao baixar imagem: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Detect content type
	contentType := http.DetectContentType(data)
	base64Data := base64.StdEncoding.EncodeToString(data)

	return fmt.Sprintf("data:%s;base64,%s", contentType, base64Data), nil
}

func fetchURLContent(url string) (string, error) {
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

func openBrowser(filename string) {
	var err error
	absPath, _ := filepath.Abs(filename)
	url := "file://" + absPath

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("plataforma não suportada")
	}
	if err != nil {
		fmt.Printf("Erro ao abrir navegador: %v\n", err)
	}
}

func cleanMarkdown(text string) string {
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```html")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	return text
}

func printBanner() {
	fmt.Println(`
RabbitAI, a plataforma jornalística do Portal O Tempo.
Versão 0.20-beta-0981271. 
Criada por AI Workers da MaxwellAI. Copyright (C) 2025 Cordya AI.
-----------------------------------------------------------------`)
}

func printUsage() {
	fmt.Println("Uso: rabbitai -u <url> [-l <llm>] [-p <persona>]")
	fmt.Println("  -u, --url     URL fonte para o artigo (OBRIGATÓRIO)")
	fmt.Println("  -l, --llm     LLM para usar: gemini, openai, claude, grok (padrão: gemini)")
	fmt.Println("  -p, --persona Tom da persona: formal, pessoal, neutra (padrão: neutra)")
	fmt.Println("  -i, --idiom   Idioma do artigo: BR, EN, ES (padrão: BR)")
	fmt.Println("\nPor favor forneça o parâmetro -u para iniciar.")
}

func exitWithError(err error) {
	fmt.Printf("Erro: %v\n", err)
	fmt.Println("Errar é Humano. Sinto muito!")
	os.Exit(1)
}
