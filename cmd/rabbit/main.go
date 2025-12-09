package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"rabbitai/internal/config"
	"rabbitai/internal/content"
	"rabbitai/internal/imagegen"
	"rabbitai/internal/llm"
	"rabbitai/internal/output"
)

const (
	DefaultLLM     = "gemini"
	DefaultPersona = "neutra"
	DefaultIdiom   = "BR"
)

func main() {
	// 1. Print Banner
	printBanner()

	// 2. Parse Flags
	configPath := flag.String("c", "config.yaml", "Path to configuration file")
	configPathLong := flag.String("config", "config.yaml", "Path to configuration file")
	urlFlag := flag.String("u", "", "URL fonte para o artigo (atalho)")
	urlFlagLong := flag.String("url", "", "URL fonte para o artigo")
	llmFlag := flag.String("l", "", "LLM para usar: gemini, openai, claude, grok")
	llmFlagLong := flag.String("llm", "", "LLM para usar: gemini, openai, claude, grok")
	personaFlag := flag.String("p", "", "Tom da persona: formal, pessoal, neutra")
	personaFlagLong := flag.String("persona", "", "Tom da persona: formal, pessoal, neutra")
	idiomFlag := flag.String("i", "", "Idioma do artigo: BR, EN, ES")
	idiomFlagLong := flag.String("idiom", "", "Idioma do artigo: BR, EN, ES")
	imageProviderFlag := flag.String("ip", "", "Image generator provider: dalle, imagen, stable-diffusion, flux, local")
	imageProviderFlagLong := flag.String("image-provider", "", "Image generator provider: dalle, imagen, stable-diffusion, flux, local")
	assetTypeFlag := flag.String("at", "", "Asset type: image, video, local")
	assetTypeFlagLong := flag.String("asset-type", "", "Asset type: image, video, local")
	imagebankPathFlag := flag.String("ib", "", "Custom imagebank path")
	imagebankPathFlagLong := flag.String("imagebank-path", "", "Custom imagebank path")
	sectionFlag := flag.String("s", "", "Section from editoriais to use (e.g., decretos-governamentais, nomeacoes)")
	sectionFlagLong := flag.String("section", "", "Section from editoriais to use (e.g., decretos-governamentais, nomeacoes)")

	flag.Parse()

	// Resolve flags (long takes precedence)
	targetConfigPath := *configPath
	if *configPathLong != "config.yaml" {
		targetConfigPath = *configPathLong
	}

	targetURL := *urlFlag
	if *urlFlagLong != "" {
		targetURL = *urlFlagLong
	}

	targetLLM := *llmFlag
	if *llmFlagLong != "" {
		targetLLM = *llmFlagLong
	}

	targetPersona := *personaFlag
	if *personaFlagLong != "" {
		targetPersona = *personaFlagLong
	}

	targetIdiom := *idiomFlag
	if *idiomFlagLong != "" {
		targetIdiom = *idiomFlagLong
	}

	targetImageProvider := *imageProviderFlag
	if *imageProviderFlagLong != "" {
		targetImageProvider = *imageProviderFlagLong
	}

	targetAssetType := *assetTypeFlag
	if *assetTypeFlagLong != "" {
		targetAssetType = *assetTypeFlagLong
	}

	targetImagebankPath := *imagebankPathFlag
	if *imagebankPathFlagLong != "" {
		targetImagebankPath = *imagebankPathFlagLong
	}

	targetSection := *sectionFlag
	if *sectionFlagLong != "" {
		targetSection = *sectionFlagLong
	}

	// 3. Load Configuration
	cfg, err := config.LoadConfig(targetConfigPath)
	if err != nil {
		fmt.Printf("Aviso: Falha ao carregar configuração: %v\n", err)
		fmt.Println("Usando configurações padrão...")
		// Create default config if file doesn't exist
		cfg = createDefaultConfig()
	}

	// 4. Resolve defaults from config if flags not provided
	if targetLLM == "" {
		targetLLM = cfg.LLM.Default
	}
	if targetPersona == "" {
		targetPersona = DefaultPersona
	}
	if targetIdiom == "" {
		targetIdiom = DefaultIdiom
	}
	if targetImageProvider == "" {
		targetImageProvider = cfg.ImageGeneration.DefaultProvider
	}
	if targetAssetType == "" {
		targetAssetType = cfg.ImageGeneration.DefaultAssetType
	}
	if targetImagebankPath != "" {
		cfg.ImageGeneration.LocalStorage.Path = targetImagebankPath
	}

	// Resolve section from flag or config
	if targetSection == "" {
		targetSection = cfg.Editorial.DefaultSection
	}

	// 4.5 Validate Section (if provided)
	if targetSection != "" {
		availableSections, err := content.GetAvailableSections(cfg.Editorial.EditoriaisPath)
		if err != nil {
			fmt.Printf("Aviso: Falha ao verificar seções disponíveis: %v\n", err)
		} else {
			sectionValid := false
			for _, section := range availableSections {
				if section == targetSection {
					sectionValid = true
					break
				}
			}
			if !sectionValid {
				fmt.Printf("Aviso: Seção '%s' não encontrada.\n", targetSection)
				if len(availableSections) > 0 {
					fmt.Printf("Seções disponíveis: %s\n", strings.Join(availableSections, ", "))
				}
				fmt.Println("Continuando sem diretrizes editoriais...")
				targetSection = "" // Clear invalid section
			}
		}
	}

	// 5. Validate Parameters
	if targetURL == "" {
		printUsage()
		os.Exit(0)
	}

	fmt.Printf("URL Alvo: %s\n", targetURL)
	fmt.Printf("LLM: %s\n", targetLLM)
	fmt.Printf("Persona: %s\n", targetPersona)
	fmt.Printf("Idioma: %s\n", targetIdiom)
	fmt.Printf("Image Provider: %s\n", targetImageProvider)
	fmt.Printf("Asset Type: %s\n", targetAssetType)
	if targetSection != "" {
		fmt.Printf("Seção: %s\n", targetSection)
	}
	fmt.Println("--------------------------------------------------------")

	// 6. Initialize LLM Provider
	llmProvider, err := llm.NewProvider(targetLLM, cfg)
	if err != nil {
		exitWithError(fmt.Errorf("falha ao inicializar LLM provider: %w", err))
	}

	// 7. Initialize Image Generator
	imageGen, err := imagegen.NewGenerator(targetImageProvider, cfg)
	if err != nil {
		exitWithError(fmt.Errorf("falha ao inicializar image generator: %w", err))
	}

	// Check if generator supports the requested asset type
	assetType := imagegen.AssetType(targetAssetType)
	if !imageGen.SupportsAssetType(assetType) {
		exitWithError(fmt.Errorf("image generator '%s' não suporta asset type '%s'", targetImageProvider, targetAssetType))
	}

	// 8. Core Logic
	fmt.Println("Gerando artigo... por favor aguarde.")

	// 8.1 Fetch Content
	fmt.Printf("Baixando conteúdo de %s...\n", targetURL)
	extractedContent, fetchErr := content.FetchURLContent(targetURL)
	if fetchErr != nil {
		fmt.Printf("Aviso: Falha ao baixar conteúdo da URL: %v\n", fetchErr)
		fmt.Println("Tentando gerar apenas com a URL...")
		extractedContent = ""
	} else {
		fmt.Printf("Conteúdo extraído: %d caracteres.\n", len(extractedContent))
	}

	// 8.1.5 Load Editorial Content
	var editorialContent string
	if targetSection != "" {
		fmt.Printf("Carregando diretrizes editoriais da seção '%s'...\n", targetSection)
		editorial, err := content.LoadEditorialContent(targetSection, cfg.Editorial.EditoriaisPath)
		if err != nil {
			fmt.Printf("Aviso: %v\n", err)
			fmt.Println("Continuando sem diretrizes editoriais...")
			editorialContent = ""
		} else if editorial != "" {
			fmt.Printf("Diretrizes editoriais carregadas: %d caracteres.\n", len(editorial))
			editorialContent = editorial
		} else {
			fmt.Println("Nenhuma diretriz editorial encontrada para esta seção.")
		}
	}

	// 8.2 Construct Prompt
	prompt := content.ConstructPrompt(targetURL, extractedContent, targetPersona, targetIdiom, editorialContent)

	// 8.3 Generate Article
	fmt.Printf("Gerando artigo com %s...\n", llmProvider.Name())
	htmlContent, err := llmProvider.GenerateText(prompt)
	if err != nil {
		exitWithError(fmt.Errorf("falha ao gerar artigo: %w", err))
	}

	// 8.3.5 Extract title and generate article path
	fmt.Println("Extraindo título do artigo...")
	title, err := output.ExtractTitle(htmlContent)
	if err != nil {
		fmt.Printf("Aviso: Falha ao extrair título: %v\n", err)
		fmt.Println("Usando título padrão...")
		title = "article"
	}

	articlePath, err := output.GenerateArticlePath(targetSection, title, cfg.Editorial.EditoriaisPath)
	if err != nil {
		exitWithError(fmt.Errorf("falha ao gerar caminho do artigo: %w", err))
	}

	// 8.4 Image Generation
	fmt.Println("Analisando necessidade de imagem...")
	imagePrompt, hasImagePrompt := output.ExtractImagePrompt(htmlContent)

	if hasImagePrompt {
		fmt.Printf("Gerando imagem com %s: %s\n", imageGen.Name(), imagePrompt)
		asset, err := imageGen.GenerateImage(imagePrompt, imagegen.ImageOptions{})
		if err != nil {
			if cfg.ImageGeneration.LocalStorage.Fallback && targetImageProvider != "local" {
				fmt.Printf("Aviso: Falha ao gerar imagem com %s: %v\n", imageGen.Name(), err)
				fmt.Println("Tentando usar imagebank local como fallback...")
				localGen := imagegen.NewLocalGenerator(&cfg.ImageGeneration.LocalStorage)
				asset, err = localGen.GenerateImage(imagePrompt, imagegen.ImageOptions{})
			}
			if err != nil {
				fmt.Printf("Aviso: Falha ao gerar imagem: %v\n", err)
			} else {
				// Process and save image
				if err := processAndSaveImage(asset, articlePath, &htmlContent, imagePrompt); err != nil {
					fmt.Printf("Aviso: %v\n", err)
				}
			}
		} else {
			// Process and save image
			if err := processAndSaveImage(asset, articlePath, &htmlContent, imagePrompt); err != nil {
				fmt.Printf("Aviso: %v\n", err)
			}
		}
	} else {
		fmt.Println("Nenhuma descrição de imagem encontrada.")
	}

	// 9. Save and Open
	err = output.SaveHTML(articlePath, htmlContent)
	if err != nil {
		exitWithError(fmt.Errorf("falha ao salvar arquivo: %w", err))
	}

	fmt.Printf("Artigo salvo em %s\n", articlePath)
	if err := output.OpenBrowser(articlePath); err != nil {
		fmt.Printf("Aviso: %v\n", err)
	}

	fmt.Println("Mais uma missão cumprida por Humanos e IA. Bem-vindo à era dos AI Citizens.")
}

func createDefaultConfig() *config.Config {
	cfg := &config.Config{}
	cfg.LLM.Default = DefaultLLM
	cfg.LLM.Providers = make(map[string]config.LLMProvider)
	cfg.LLM.Providers["gemini"] = config.LLMProvider{
		APIKey: os.Getenv("GEMINI_API_KEY"),
		Model:  "gemini-2.0-flash-exp",
	}
	cfg.LLM.Providers["openai"] = config.LLMProvider{
		APIKey: os.Getenv("OPENAI_API_KEY"),
		Model:  "gpt-4o",
	}
	cfg.LLM.Providers["claude"] = config.LLMProvider{
		APIKey: os.Getenv("CLAUDE_API_KEY"),
		Model:  "claude-3-5-sonnet-20241022",
	}
	cfg.LLM.Providers["grok"] = config.LLMProvider{
		APIKey: os.Getenv("GROK_API_KEY"),
		Model:  "grok-3",
	}

	cfg.ImageGeneration.DefaultProvider = "dalle"
	cfg.ImageGeneration.DefaultAssetType = "image"
	cfg.ImageGeneration.Providers = make(map[string]config.ImageProvider)
	cfg.ImageGeneration.Providers["dalle"] = config.ImageProvider{
		APIKey: os.Getenv("OPENAI_API_KEY"),
		Model:  "dall-e-3",
		Size:   "1024x1024",
	}
	cfg.ImageGeneration.LocalStorage = config.LocalStorageConfig{
		Enabled:  true,
		Path:     "./imagebank",
		Fallback: true,
	}

	cfg.Editorial.DefaultSection = ""
	cfg.Editorial.EditoriaisPath = "./editoriais"

	return cfg
}

func printBanner() {
	fmt.Println(`
Rabbit CLI, a plataforma jornalística do Portal O Tempo.
Versão 1.0.0
Criada por AI Workers da MaxwellAI. Copyright (C) 2025 Cordya AI.
-----------------------------------------------------------------`)
}

func printUsage() {
	fmt.Println("Uso: rabbit -u <url> [opções]")
	fmt.Println("\nOpções:")
	fmt.Println("  -u, --url            URL fonte para o artigo (OBRIGATÓRIO)")
	fmt.Println("  -c, --config          Caminho para arquivo de configuração (padrão: config.yaml)")
	fmt.Println("  -l, --llm            LLM para usar: gemini, openai, claude, grok")
	fmt.Println("  -p, --persona         Tom da persona: formal, pessoal, neutra")
	fmt.Println("  -i, --idiom          Idioma do artigo: BR, EN, ES")
	fmt.Println("  -ip, --image-provider Image generator: dalle, imagen, stable-diffusion, flux, local")
	fmt.Println("  -at, --asset-type    Tipo de asset: image, video, local")
	fmt.Println("  -ib, --imagebank-path Caminho customizado para imagebank")
	fmt.Println("  -s, --section       Seção de editoriais para usar (ex: decretos-governamentais, nomeacoes)")
	fmt.Println("\nPor favor forneça o parâmetro -u para iniciar.")
}

func exitWithError(err error) {
	fmt.Printf("Erro: %v\n", err)
	fmt.Println("Errar é Humano. Sinto muito!")
	os.Exit(1)
}

// processAndSaveImage processes an image asset and saves it to disk, updating HTML
func processAndSaveImage(asset *imagegen.Asset, articlePath string, htmlContent *string, imagePrompt string) error {
	fmt.Println("Imagem gerada. Processando e salvando...")
	
	// Extract image data from data URI
	imageDataURI := string(asset.Data)
	imageBytes, extractedContentType, err := output.ExtractImageFromDataURI(imageDataURI)
	if err != nil {
		return fmt.Errorf("falha ao extrair dados da imagem: %w", err)
	}

	// Use asset.ContentType if available, otherwise use extracted content type
	contentType := asset.ContentType
	if contentType == "" {
		contentType = extractedContentType
	}

	// Generate image path
	imagePath, err := output.GenerateImagePath(articlePath, contentType)
	if err != nil {
		return fmt.Errorf("falha ao gerar caminho da imagem: %w", err)
	}

	// Save image to file
	if err := output.SaveImage(imagePath, imageBytes); err != nil {
		return fmt.Errorf("falha ao salvar imagem: %w", err)
	}

	fmt.Printf("Imagem salva em %s\n", imagePath)
	
	// Update HTML to reference local image file
	*htmlContent = output.InjectImageIntoHTML(*htmlContent, imagePath, imagePrompt)
	
	return nil
}
