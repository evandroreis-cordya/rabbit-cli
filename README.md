# Rabbit CLI

```text
██████╗  █████╗ ██████╗ ██████╗ ██╗████████╗     ██████╗██╗     ██╗
██╔══██╗██╔══██╗██╔══██╗██╔══██╗██║╚══██╔══╝    ██╔════╝██║     ██║
██████╔╝███████║██████╔╝██████╔╝██║   ██║       ██║     ██║     ██║
██╔══██╗██╔══██║██╔══██╗██╔══██╗██║   ██║       ██║     ██║     ██║
██║  ██║██║  ██║██████╔╝██████╔╝██║   ██║       ╚██████╗███████╗██║
╚═╝  ╚═╝╚═╝  ╚═╝╚═════╝ ╚═════╝ ╚═╝   ╚═╝        ╚═════╝╚══════╝╚═╝
                                                                    
```

Rabbit CLI is a modular, configurable command-line tool for the Rabbit Platform. The Rabbit Platform generates compelling, self-contained journalistic articles in based on the contents from DOU (Diário Oficial da União) using Large Language Models (LLMs). It extracts content from DOU and uses AI to create polished, modern, and minimalist articles with configurable image generation from multiple providers.

**Version:** 1.0.0  
**Created by:** AI Workers from MaxwellAI  
**Copyright (C) 2025 Cordya AI**

## Features

- 📰 **Article Generation**: Creates professional journalistic articles from any URL
- 🤖 **Multiple LLM Support**: Works with Gemini, OpenAI, Claude, and Grok
- 🎨 **Self-Contained HTML**: Generates complete HTML files with embedded CSS and JavaScript
- 🖼️ **Multiple Image Generators**: Supports DALL-E 3, Google Imagen 3, Stable Diffusion, Flux, and local imagebank
- 🎯 **Configurable**: YAML-based configuration with environment variable support
- 🌍 **Multi-Language**: Supports Brazilian Portuguese (BR), English (EN), and Spanish (ES)
- 🎭 **Persona Options**: Choose between formal, personal, or neutral writing styles
- 📁 **Local Image Storage**: Use local imagebank directory for image assets
- 🌐 **Auto-Browser**: Automatically opens the generated article in your default browser
- 📂 **Section-Based Organization**: Organize articles by editorial sections with automatic directory structure
- 📝 **Editorial Guidelines**: Load custom editorial guidelines from `EDITORIA.md` files for consistent article style

## Architecture

Rabbit CLI uses a modular architecture with separate packages for:

- **LLM Providers** (`internal/llm/`): Modular LLM implementations (Gemini, OpenAI, Claude, Grok)
- **Image Generators** (`internal/imagegen/`): Multiple image generation backends
- **Configuration** (`internal/config/`): YAML configuration with environment variable resolution
- **Content Processing** (`internal/content/`): URL fetching and content extraction
- **Output** (`internal/output/`): HTML manipulation and browser utilities

## Installation

### Prerequisites

- Go 1.25.4 or later
- API keys for at least one LLM provider (configured via `config.yaml` or environment variables)

### Build from Source

```bash
git clone <repository-url>
cd rabbitai
go build -o rabbit ./cmd/rabbit
```

Or run directly:

```bash
go run ./cmd/rabbit -u <URL>
```

## Configuration

Rabbit CLI uses a YAML configuration file (`config.yaml`) for settings. The configuration supports environment variable placeholders using `${VAR_NAME}` syntax.

### Configuration File Structure

```yaml
llm:
  default: gemini
  providers:
    gemini:
      api_key: ${GEMINI_API_KEY}
      model: gemini-2.0-flash-exp
    openai:
      api_key: ${OPENAI_API_KEY}
      model: gpt-4o
    # ... other providers

image_generation:
  default_provider: dalle
  default_asset_type: image
  providers:
    dalle:
      api_key: ${OPENAI_API_KEY}
      model: dall-e-3
      size: 1024x1024
    imagen:
      api_key: ${GOOGLE_API_KEY}
      model: imagen-3
      size: 1024x1024
      project_id: ${GOOGLE_CLOUD_PROJECT}
    stable_diffusion:
      api_key: ${STABILITY_API_KEY}
      model: stable-diffusion-xl-1024-v1-0
    flux:
      api_key: ${FLUX_API_KEY}
      model: flux-pro
  local_storage:
    enabled: true
    path: ./imagebank
    fallback: true
```

### Environment Variables

Set the following environment variables (or use `${VAR_NAME}` in config.yaml):

- `GEMINI_API_KEY` - Google Gemini API key
- `OPENAI_API_KEY` - OpenAI API key (for GPT-4 and DALL-E 3)
- `CLAUDE_API_KEY` - Anthropic Claude API key
- `GROK_API_KEY` - xAI Grok API key
- `GOOGLE_API_KEY` - Google API key (for Imagen)
- `GOOGLE_CLOUD_PROJECT` - Google Cloud Project ID (for Imagen)
- `STABILITY_API_KEY` - Stability AI API key
- `FLUX_API_KEY` - Black Forest Labs Flux API key

### Editorial Guidelines Configuration

The `editorial` section in `config.yaml` configures section-based organization:

```yaml
editorial:
  default_section: ""  # Optional default section
  editoriais_path: "./editoriais"  # Path to editoriais directory
```

**Editorial Guidelines**: Each section in the `editoriais/` directory can contain an `EDITORIA.md` file with custom writing instructions that are automatically loaded when generating articles for that section.

**Example `editoriais/portarias-anvisa/EDITORIA.md`:**
```markdown
# Editoria Portarias ANVISA

Preciso que encontre, na página do Diário Oficial da União (DOU), 
quais resoluções da ANVISA entraram em vigor e qual foi a respectiva motivação.

A partir disso, produza um conteúdo informativo, no formato jornalístico, 
por meio da técnica da pirâmide invertida (o que, quem, como, quando, onde, por quê)...
```

When you use `-s portarias-anvisa`, the tool automatically:
1. Loads the editorial guidelines from `editoriais/portarias-anvisa/EDITORIA.md`
2. Includes them in the prompt sent to the LLM
3. Saves the article in `editoriais/portarias-anvisa/` directory

## Usage

### Basic Usage

```bash
./rabbit -u <URL>
```

### Command-Line Options

| Flag | Long Form | Description | Default |
|------|-----------|-------------|---------|
| `-u` | `--url` | Source URL for the article (required) | - |
| `-c` | `--config` | Path to configuration file | `config.yaml` |
| `-l` | `--llm` | LLM provider: `gemini`, `openai`, `claude`, `grok` | From config |
| `-p` | `--persona` | Writing tone: `formal`, `pessoal`, `neutra` | `neutra` |
| `-i` | `--idiom` | Article language: `BR`, `EN`, `ES` | `BR` |
| `-ip` | `--image-provider` | Image generator: `dalle`, `imagen`, `stable-diffusion`, `flux`, `local` | From config |
| `-at` | `--asset-type` | Asset type: `image`, `video`, `local` | `image` |
| `-ib` | `--imagebank-path` | Custom imagebank directory path | From config |
| `-s` | `--section` | Editorial section from editoriais (e.g., `portarias-anvisa`, `nomeacoes`) | From config |

### Examples

#### Basic Examples

Generate an article using default settings (Gemini, neutral tone, Portuguese, DALL-E):

```bash
./rabbit -u https://example.com/article
```

Generate an article with a specific URL:

```bash
./rabbit --url https://www.in.gov.br/web/dou/-/portaria-n-123-de-15-de-janeiro-de-2025
```

#### LLM Provider Examples

Use OpenAI GPT-4o for article generation:

```bash
./rabbit -u https://example.com/article -l openai
```

Use Claude 3.5 Sonnet for article generation:

```bash
./rabbit -u https://example.com/article --llm claude
```

Use Grok for article generation:

```bash
./rabbit -u https://example.com/article -l grok
```

#### Persona and Language Examples

Generate a formal press release-style article in Portuguese:

```bash
./rabbit -u https://example.com/article -p formal
```

Generate a personal, relaxed article in English:

```bash
./rabbit -u https://example.com/article -p pessoal -i EN
```

Generate a neutral journalistic article in Spanish:

```bash
./rabbit -u https://example.com/article -p neutra -i ES
```

#### Image Generation Examples

Generate an article with DALL-E 3 images:

```bash
./rabbit -u https://example.com/article -ip dalle
```

Generate an article with Google Imagen 3:

```bash
./rabbit -u https://example.com/article --image-provider imagen
```

Generate an article with Stable Diffusion:

```bash
./rabbit -u https://example.com/article -ip stable-diffusion
```

Generate an article with Flux Pro:

```bash
./rabbit -u https://example.com/article -ip flux
```

Use local imagebank for images:

```bash
./rabbit -u https://example.com/article -ip local --asset-type local
```

#### Combined Examples

Generate a formal English article using Claude with DALL-E images:

```bash
./rabbit -u https://example.com/article -l claude -p formal -i EN -ip dalle
```

Generate a personal Portuguese article using Gemini with local imagebank:

```bash
./rabbit -u https://example.com/article -l gemini -p pessoal -i BR -ip local -at local
```

Generate a neutral Spanish article using OpenAI with Flux images:

```bash
./rabbit -u https://example.com/article --llm openai --persona neutra --idiom ES --image-provider flux
```

#### Custom Configuration Examples

Use a custom configuration file:

```bash
./rabbit -u https://example.com/article -c /path/to/custom-config.yaml
```

Use a custom imagebank directory:

```bash
./rabbit -u https://example.com/article -ip local --imagebank-path /path/to/my/images
```

#### Section-Based Organization Examples

Generate an article for a specific editorial section (e.g., ANVISA ordinances):

```bash
./rabbit -u https://www.in.gov.br/web/dou/-/portaria-anvisa-123 -s portarias-anvisa
```

Generate an article for government decrees section:

```bash
./rabbit -u https://www.in.gov.br/web/dou/-/decreto-123 -s decretos-governamentais -p formal
```

Generate an article for nominations section with editorial guidelines:

```bash
./rabbit -u https://www.in.gov.br/web/dou/-/nomeacao-123 -s nomeacoes -l claude -i BR
```

Available sections include:
- `portarias-anvisa` - ANVISA ordinances
- `portarias-ans` - ANS ordinances
- `decretos-governamentais` - Government decrees
- `nomeacoes` - Appointments/nominations
- `exoneracoes` - Dismissals
- `medidas-provisorias` - Provisional measures
- `destinacao-de-verbas` - Budget allocations
- `balancos-governamentais` - Government balances
- `reformas-ministeriais` - Ministerial reforms

#### Real-World Use Cases

Generate an article from a DOU (Diário Oficial da União) URL with section:

```bash
./rabbit -u https://www.in.gov.br/web/dou/-/portaria-n-123-de-15-de-janeiro-de-2025 -s portarias-anvisa -p formal -i BR
```

Generate an article from DOU with editorial guidelines automatically loaded:

```bash
./rabbit -u https://www.in.gov.br/web/dou/-/portaria-ans-456 -s portarias-ans -l gemini -p neutra
# Automatically loads EDITORIA.md from editoriais/portarias-ans/ if it exists
```

Generate an article from a news article URL in English:

```bash
./rabbit -u https://news.example.com/article -l openai -i EN -p neutra -ip dalle
```

Generate an article with fallback to local images (if API fails):

```bash
./rabbit -u https://example.com/article -ip dalle
# If DALL-E fails, automatically falls back to local imagebank (if enabled in config)
```

#### Quick Reference Examples

All flags combined example:

```bash
./rabbit \
  --url https://example.com/article \
  --config config.yaml \
  --llm gemini \
  --persona neutra \
  --idiom BR \
  --image-provider dalle \
  --asset-type image \
  --imagebank-path ./imagebank \
  --section portarias-anvisa
```

Short flags version:

```bash
./rabbit -u https://example.com/article -c config.yaml -l gemini -p neutra -i BR -ip dalle -at image -ib ./imagebank -s portarias-anvisa
```

## How It Works

1. **Configuration Loading**: Loads settings from `config.yaml` (or custom path)
2. **Content Extraction**: Fetches and extracts text content from the provided URL
3. **Editorial Guidelines Loading**: If a section is specified, loads `EDITORIA.md` from `editoriais/{section}/` directory
4. **Prompt Construction**: Builds a detailed prompt based on the selected persona, language, extracted content, and editorial guidelines
5. **LLM Generation**: Calls the selected LLM API to generate the article HTML
6. **Image Generation**: If the LLM includes an image prompt comment, generates an image using the configured provider
7. **Fallback Handling**: Falls back to local imagebank if API generation fails (if enabled)
8. **File Output**: Saves the complete HTML file to `editoriais/{section}/{title}-{date}-{time}.html` (or `uncategorized/` if no section)
9. **Image Saving**: Saves generated images alongside the article with matching filename
10. **Browser Launch**: Automatically opens the article in your default browser

## Persona Styles

- **Formal** (`formal`): Press release-style article written with strict journalistic guidelines
- **Personal** (`pessoal`): Relaxed, personal article with a "what it means to me" approach
- **Neutral** (`neutra`): Neutral journalistic style mimicking articles from Portal O Tempo

## Supported LLMs

- **Gemini**: Uses `gemini-2.0-flash-exp` model
- **OpenAI**: Uses `gpt-4o` model for text generation
- **Claude**: Uses `claude-3-5-sonnet-20241022` model
- **Grok**: Uses `grok-3` model

## Supported Image Generators

- **DALL-E 3**: OpenAI's image generation model
- **Imagen 3**: Google's image generation model (requires Google Cloud Project ID)
- **Stable Diffusion**: Stability AI's image generation model
- **Flux**: Black Forest Labs' Flux Pro model
- **Local**: Uses images from the `./imagebank` directory (fuzzy matching by filename)

## Section-Based Organization

Rabbit CLI organizes articles by editorial sections, making it easy to manage content for different categories. Each section can have its own editorial guidelines and directory structure.

### Setting Up Sections

1. Create section directories in `editoriais/`:
   ```bash
   mkdir -p editoriais/portarias-anvisa
   mkdir -p editoriais/nomeacoes
   mkdir -p editoriais/decretos-governamentais
   ```

2. Optionally add `EDITORIA.md` files to each section with editorial guidelines:
   ```bash
   # Example: editoriais/portarias-anvisa/EDITORIA.md
   # Contains writing style guidelines for ANVISA ordinances
   ```

3. Use the `-s` or `--section` flag when generating articles:
   ```bash
   ./rabbit -u <URL> -s portarias-anvisa
   ```

### Article Path Structure

Articles are automatically saved with descriptive filenames:
- Format: `{sanitized-title}-{YYYY-MM-DD-HHMM}.html`
- Location: `editoriais/{section}/` or `editoriais/uncategorized/`
- Images: Saved alongside with matching base filename

## Local Imagebank

The local imagebank feature allows you to store images in a directory (`./imagebank` by default) and have Rabbit CLI automatically select matching images based on the prompt. The system uses fuzzy matching on filenames to find relevant images.

To use the imagebank:

1. Create an `imagebank` directory in your project root
2. Add image files (JPG, PNG, GIF, WebP) to the directory
3. Use descriptive filenames that match your article topics
4. Use `--image-provider local` or set `default_provider: local` in config

## Output

The application generates organized, self-contained HTML files with the following structure:

### File Organization

Articles are saved in a structured directory layout:
- **With section**: `editoriais/{section}/{title}-{YYYY-MM-DD-HHMM}.html`
- **Without section**: `editoriais/uncategorized/{title}-{YYYY-MM-DD-HHMM}.html`
- **Images**: Saved alongside articles as `{title}-{YYYY-MM-DD-HHMM}.{ext}` (jpeg, png, etc.)

Example:
```
editoriais/
  ├── portarias-anvisa/
  │   ├── portaria-anvisa-123-2025-01-15-1430.html
  │   └── portaria-anvisa-123-2025-01-15-1430.jpeg
  ├── nomeacoes/
  │   └── nomeacao-ministro-2025-01-15-1500.html
  └── uncategorized/
      └── article-2025-01-15-1200.html
```

### HTML File Contents

Each generated HTML file includes:

- Semantic HTML5 structure
- Embedded CSS styling (minimalist design)
- Embedded JavaScript (subtle animations and interactions)
- Responsive layout for desktop and mobile
- Professional typography and spacing
- Header image (local file reference or base64 encoded)
- No external dependencies (fully self-contained)

## Migration from RabbitAI

If you're migrating from the old RabbitAI version:

1. **Create `config.yaml`**: Copy the template from the repository and set your API keys
2. **Update command**: Use `rabbit` instead of `rabbitai`
3. **Build command**: Build with `go build -o rabbit ./cmd/rabbit`
4. **API Keys**: Move from hardcoded values to `config.yaml` or environment variables

## Security Notes

- API keys should never be committed to version control
- Use environment variables or secure configuration management
- The application escapes HTML content in image alt attributes to prevent XSS attacks
- Local imagebank paths are validated to prevent directory traversal

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Credits

- **Created by**: AI Workers from MaxwellAI
- **Copyright**: © 2025 Cordya AI
- **Platform**: Rabbit CLI, a journalistic platform from Portal O Tempo

## Notes

- The application requires internet connectivity to fetch URLs and call LLM APIs
- Image generation requires appropriate API keys for the selected provider
- The generated HTML files are self-contained and can be shared or hosted independently
- Google Imagen requires a Google Cloud Project ID configured in `config.yaml` or `GOOGLE_CLOUD_PROJECT` environment variable

---

### Quote

> "Mais uma missão cumprida por Humanos e IA. Bem-vindo à era dos AI Citizens."
