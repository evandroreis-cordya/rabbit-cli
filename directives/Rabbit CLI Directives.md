# Rabbit CLI

Create a Go lang, command-line application, named rabbit, that generates a compelling and engaging article (as a self-contained HTML file) using as the source the contents from the URL passed and calling a LLM (also passed as a parameter, with the default being "gemini"). The article must contain. It MUST use the LLM model defined generate: title, one-liner, article content, image to illustrate the article, captions, and references. The application must implement the following features described below:

## Accept command-line parameters

- **-u (or --url)**: a well-formed URL that will serve as the source for the article to be created.
- **-l (or --llm)**: the name of the LLM to be called (defined below). If none, it must use the default LLM (gemini). If the parameter is "gemini", it must call Gemini API (according to the appropriate instructions below); If the parameter is "openai", it must call OpenAI API (according to the appropriate instructions below); If the parameter is "claude", it must call Claude API (according to the appropriate instructions below); If the parameter is "grok", it must call Grok (xAI) API (according to the appropriate instructions below);
- **-p (or --persona)**: defines the tone of the article: "formal", a formal, press release-type article, written with most strict journalistic guidelines; "personal", a more relax, personal, almost 'what it means to me'-kind of article, and "neutral" (the default), a article that mimics the style of the articles published in the portal O Tempo (<https://www.otempo.com.br/>)

### Prompt for Code Generation: Gemini API Call

"Generate a self-contained Go program that calls the Google Gemini API over raw HTTP using only the standard library (net/http, encoding/json, etc.). Use variables for every configurable API parameter instead of hardcoding them.

Requirements:

- Declare variables for:
  - apiKey (string) – the Gemini API key.
  - modelName (string) – for example 'gemini-2.5-flash' or any caller-supplied value.
  - endpointBaseURL (string) – default '<https://generativelanguage.googleapis.com>'.
  - apiVersion (string) – for example 'v1beta' or 'v1'.
  - temperature (float64).
  - topP (float64).
  - topK (int).
  - maxOutputTokens (int).
  - requestPrompt (string) – the user text to send.
  - safetyCategory (string) – for example 'HARM_CATEGORY_HATE_SPEECH'.
  - safetyThreshold (string) – for example 'BLOCK_MEDIUM_AND_ABOVE'.
- Build the full URL as:
  endpointBaseURL + '/' + apiVersion + '/models/' + modelName + ':generateContent'
  and pass the apiKey as the query parameter 'key' (e.g. '?key=' + apiKey).
- Construct the JSON request body to match the Gemini generateContent REST schema:
  {
    'contents': [
      {
        'role': 'user',
        'parts': [
          {'text': requestPrompt}
        ]
      }
    ],
    'generationConfig': {
      'temperature': temperature,
      'topP': topP,
      'topK': topK,
      'maxOutputTokens': maxOutputTokens
    },
    'safetySettings': [
      {
        'category': safetyCategory,
        'threshold': safetyThreshold
      }
    ]
  }
- Define Go structs that map to this JSON shape and marshal them with encoding/json.
- Use http.NewRequest with method POST, set:
  - URL to the constructed endpoint URL.
  - Header 'Content-Type' to 'application/json'.
- Execute the request with http.Client.Do, handle any errors, and:
  - If status code is not 2xx, print the status and response body as an error.
  - If status is 2xx, decode the JSON response into Go structs that cover at least:
    - candidates
    - content.parts.text
  - Print the first candidate's first text part to stdout.
- Organize the code so all variable configuration (apiKey, modelName, endpointBaseURL, apiVersion, temperature, topP, topK, maxOutputTokens, requestPrompt, safetyCategory, safetyThreshold) is grouped near the top of main() or in a separate config function, making it easy to change without editing the rest of the code.
- Include basic context usage (context.Background) with http.Client if you structure it that way, but do not introduce external dependencies.

Produce idiomatic, formatted Go code ready to paste called as a function within main.go and run with 'go run .'

Sources
[1] Gemini API quickstart - Google AI for Developers <https://ai.google.dev/gemini-api/docs/quickstart>
[2] Gemini API reference | Google AI for Developers <https://ai.google.dev/api>
[3] Using Gemini API keys | Google AI for Developers <https://ai.google.dev/gemini-api/docs/api-key>
[4] Get started with Gemini using the REST API - Colab - Google <https://colab.research.google.com/github/google/generative-ai-docs/blob/main/site/en/gemini-api/docs/get-started/rest.ipynb>
[5] google-gemini/api-examples <https://github.com/google-gemini/api-examples>
[6] Building an AI-Powered CLI with Golang and Google Gemini <https://dev.to/pradumnasaraf/building-an-ai-powered-cli-with-golang-and-google-gemini-45a1>
[7] Gemini Example with Go <https://www.pixelstech.net/article/1734225104-gemini-example-with-go>
[8] go-genai/example_test.go at main <https://github.com/googleapis/go-genai/blob/main/example_test.go>
[9] Gemini API: cURL command successfully calls API, but my ... <https://stackoverflow.com/questions/78564261/gemini-api-curl-command-successfully-calls-api-but-my-program-does-not>
[10] Generate content with the Gemini API in Vertex AI <https://docs.cloud.google.com/vertex-ai/generative-ai/docs/model-reference/inference>
[11] Google Generative AI provider <https://ai-sdk.dev/providers/ai-sdk-providers/google-generative-ai>
[12] Use Gemini AI API Without Coding <https://www.youtube.com/watch?v=sScZmKaLYq8>
[13] Gemini API: Getting started with Gemini models - Colab - Google <https://colab.research.google.com/github/google-gemini/cookbook/blob/main/quickstarts/Get_started.ipynb>
[14] Google Gen AI SDK documentation <https://googleapis.github.io/python-genai/>
[15] working with gemini api in go : r/golang <https://www.reddit.com/r/golang/comments/1bj970w/working_with_gemini_api_in_go/>
[16] Generative Language API v1beta2 <https://docs.cloud.google.com/go/docs/reference/cloud.google.com/go/ai/latest/generativelanguage/apiv1beta2>
[17] Gemini API in Vertex AI quickstart <https://docs.cloud.google.com/vertex-ai/generative-ai/docs/start/quickstart>
[18] All methods | Gemini API - Google AI for Developers <https://ai.google.dev/api/all-methods>
[19] Gerar conteúdo com a API Gemini na Vertex AI <https://docs.cloud.google.com/vertex-ai/generative-ai/docs/model-reference/inference?hl=pt-br>
[20] Gemini API: Prompting Quickstart with REST <https://github.com/google-gemini/cookbook/blob/main/quickstarts/rest/Prompting_REST.ipynb>

### Prompt for Code Generation: OpenAI API Call

"Generate a self-contained Go program that calls the OpenAI API over raw HTTP using only the standard library (net/http, encoding/json, etc.). Use variables for every configurable API parameter instead of hardcoding them.

Requirements:

- Declare variables for:
  - apiKey (string) – the OpenAI API key.
  - baseURL (string) – default '<https://api.openai.com/v1>'.
  - endpointPath (string) – for example '/chat/completions'.
  - modelName (string) – for example 'gpt-4.1' or any caller-supplied value.
  - temperature (float64).
  - topP (float64).
  - maxTokens (int).
  - presencePenalty (float64).
  - frequencyPenalty (float64).
  - userPrompt (string) – the user content to send.
  - systemPrompt (string) – an optional system instruction string.
- Build the full URL as:
  baseURL + endpointPath
- Use HTTP Bearer authentication for the API key by setting the header:
  'Authorization: Bearer ' + apiKey
- Construct the JSON request body to match the OpenAI Chat Completions schema:
  {
    'model': modelName,
    'messages': [
      {
        'role': 'system',
        'content': systemPrompt
      },
      {
        'role': 'user',
        'content': userPrompt
      }
    ],
    'temperature': temperature,
    'top_p': topP,
    'max_tokens': maxTokens,
    'presence_penalty': presencePenalty,
    'frequency_penalty': frequencyPenalty
  }
- Define Go structs that map to this JSON shape and marshal them with encoding/json.
- Use http.NewRequest with method POST, set:
  - URL to the constructed endpoint URL.
  - Headers:
    - 'Content-Type': 'application/json'
    - 'Authorization': 'Bearer ' + apiKey
- Execute the request with http.Client.Do, handle any errors, and:
  - If status code is not 2xx, print the status and response body as an error.
  - If status is 2xx, decode the JSON response into Go structs that cover at least:
    - choices
    - choices.message.content
  - Print the assistant reply from choices.message.content to stdout.
- Organize the code so all variable configuration (apiKey, baseURL, endpointPath, modelName, temperature, topP, maxTokens, presencePenalty, frequencyPenalty, userPrompt, systemPrompt) is grouped near the top of main() or in a separate config function, making it easy to change without editing the rest of the code.
- Include basic context usage (context.Background) with http.Client if you structure it that way, but do not introduce external dependencies.

Produce idiomatic, formatted Go code ready to paste called as a function within main.go and run with 'go run .'

Sources
[1] API Reference <https://platform.openai.com/docs/api-reference/chat>
[2] Open AI authentication API keys - OpenAI Platform <https://platform.openai.com/docs/api-reference/authentication>
[3] Completions API <https://platform.openai.com/docs/guides/completions>
[4] Create chat completion <https://platform.openai.com/docs/api-reference/chat/create>
[5] API Reference - OpenAI Platform <https://platform.openai.com/docs/api-reference/introduction>
[6] API Reference - OpenAI API <https://platform.openai.com/docs/api-reference>
[7] api-reference/completions/create <https://platform.openai.com/docs/api-reference/completions>
[8] OpenAI Chat Completion Object documentation <https://platform.openai.com/docs/api-reference/chat/object>
[9] OpenAI Chat Completions API - Sensedia Docs <https://docs.sensedia.com/en/api-management-guide/Latest/other-info/openai-chat-completions.html>
[10] Help to interact between openai with curl and with python <https://community.openai.com/t/help-to-interact-between-openai-with-curl-and-with-python/471482>
[11] How do I authenticate API requests with OpenAI? - Milvus <https://milvus.io/ai-quick-reference/how-do-i-authenticate-api-requests-with-openai>
[12] Azure OpenAI in Microsoft Foundry Models REST API ... <https://learn.microsoft.com/en-us/azure/ai-foundry/openai/reference?view=foundry-classic>
[13] Authenticate to Azure OpenAI API - Azure API Management <https://learn.microsoft.com/en-us/azure/api-management/api-management-authenticate-authorize-azure-openai>
[14] OpenAI API Chat Completions curl command Detailed ... <https://www.youtube.com/watch?v=Ya25_m-Jn5o>
[15] OpenAI Chat :: Spring AI Reference <https://docs.spring.io/spring-ai/reference/api/chat/openai-chat.html>
[16] cURL OpenAI API: Step-by-Step Tutorial for Beginners <https://muneebdev.com/curl-openai-api-tutorial/>
[17] OpenAI API chat/completions endpoint <https://docs.openvino.ai/2024/openvino-workflow/model-server/ovms_docs_rest_api_chat.html>
[18] You didn't provide an API key. You need to provide your API key in ... <https://community.openai.com/t/you-didnt-provide-an-api-key-you-need-to-provide-your-api-key-in-an-authorization-header-using-bearer-auth/561756>
[19] How can I execute a cURL command to interact with ... <https://community.openai.com/t/how-can-i-execute-a-curl-command-to-interact-with-the-openai-api-in-a-colab-notebook/553081>
[20] API Reference <https://platform.openai.com/docs/api-reference/completions?lang=python>

### Prompt for Code Generation: Claude API Call

"Generate a self-contained Go program that calls the Anthropic Claude Messages API over raw HTTP using only the standard library (net/http, encoding/json, etc.). Use variables for every configurable API parameter instead of hardcoding them.

Requirements:

- Declare variables for:
  - apiKey (string) – the Anthropic API key.
  - baseURL (string) – default '<https://api.anthropic.com>'.
  - endpointPath (string) – default '/v1/messages'.
  - apiVersion (string) – for example '2023-06-01'.
  - modelName (string) – for example 'claude-3-5-sonnet-20241022' or any caller-supplied value.
  - maxTokens (int).
  - temperature (float64).
  - topP (float64).
  - userPrompt (string) – the user content to send.
  - systemPrompt (string) – an optional system instruction string.
- Build the full URL as:
  baseURL + endpointPath
- Use API key authentication by setting the header:
  'x-api-key: ' + apiKey
- Set the required version header:
  'anthropic-version: ' + apiVersion
- Construct the JSON request body to match the Claude Messages API schema:
  {
    'model': modelName,
    'max_tokens': maxTokens,
    'temperature': temperature,
    'top_p': topP,
    'system': systemPrompt,
    'messages': [
      {
        'role': 'user',
        'content': userPrompt
      }
    ]
  }
- Define Go structs that map to this JSON shape and marshal them with encoding/json.
- Use http.NewRequest with method POST, set:
  - URL to the constructed endpoint URL.
  - Headers:
    - 'Content-Type': 'application/json'
    - 'x-api-key': apiKey
    - 'anthropic-version': apiVersion
- Execute the request with http.Client.Do, handle any errors, and:
  - If status code is not 2xx, print the status and response body as an error.
  - If status is 2xx, decode the JSON response into Go structs that cover at least:
    - content (array)
    - the first text block in content (e.g., content.text) representing the assistant reply
  - Print the assistant reply text to stdout.
- Organize the code so all variable configuration (apiKey, baseURL, endpointPath, apiVersion, modelName, maxTokens, temperature, topP, userPrompt, systemPrompt) is grouped near the top of main() or in a separate config function, making it easy to change without editing the rest of the code.
- Include basic context usage (context.Background) with http.Client if you structure it that way, but do not introduce external dependencies.

Produce idiomatic, formatted Go code ready to paste called as a function within main.go and run with 'go run .'

Sources
[1] Messages API reference - Documentation - Claude Docs <https://docs.claude.com/en/api/messages>
[2] Using the Messages API - Claude Docs <https://platform.claude.com/docs/en/build-with-claude/working-with-messages>
[3] Anthropic Claude API Key: The Essential Guide - Nightfall AI <https://www.nightfall.ai/ai-security-101/anthropic-claude-api-key>
[4] Anthropic Claude Messages API - Amazon Bedrock <https://docs.aws.amazon.com/bedrock/latest/userguide/model-parameters-anthropic-claude-messages.html>
[5] Anthropic Academy: Claude API Development Guide <https://www.anthropic.com/learn/build-with-claude>
[6] List Message Batches | Claude API Reference <https://console.anthropic.com/docs/en/api/python/messages/batches/list>
[7] Anthropic.HttpClient.Utils — anthropic_community v0.4.3 - Hexdocs <https://hexdocs.pm/anthropic_community/Anthropic.HttpClient.Utils.html>
[8] Anthropic - LiteLLM <https://docs.litellm.ai/docs/providers/anthropic>
[9] Anthropic Chat :: Spring AI Reference <https://docs.spring.io/spring-ai/reference/api/chat/anthropic-chat.html>
[10] Anthropic Chat Messages - New API <https://docs.newapi.pro/en/api/anthropic-chat/>
[11] Anthropic API | DeepSeek API Docs <https://api-docs.deepseek.com/guides/anthropic_api>
[12] Anthropic – agentgateway | Agent Connectivity Solved <https://agentgateway.dev/docs/llm/providers/anthropic/>
[13] Mensagens - Claude Docs - Home - Anthropic <https://anthropic.mintlify.app/pt/api/messages>
[14] API Messages do Claude da Anthropic - Amazon Bedrock <https://docs.aws.amazon.com/pt_br/bedrock/latest/userguide/model-parameters-anthropic-claude-messages.html>
[15] Dapatkan Kunci API - Claude Docs - Home - Anthropic <https://anthropic.mintlify.app/id/api/admin-api/apikeys/get-api-key>
[16] Melhorar um prompt - Claude Docs - Anthropic <https://docs.anthropic.com/pt/api/prompt-tools-improve>
[17] List Message Batches - Claude Docs - Home - Anthropic <https://anthropic.mintlify.app/en/api/listing-message-batches>
[18] [BUG] Fix the docs. You use --header to add an authorization header ... <https://github.com/anthropics/claude-code/issues/2324>
[19] Anthropic Claude Code Execution Tool via API - AI Engineer Guide <https://aiengineerguide.com/blog/anthropic-claude-code-execution-tool/>
[20] Anthropic API | Documentation | Postman API Network <https://www.postman.com/ai-engineer/generative-ai-apis/documentation/lqv1fm6/anthropic-api>

## Prompt for Code Generation: Grok (xAI) API Call

"Generate a self-contained Go program that calls the xAI Grok Chat Completions API over raw HTTP using only the standard library (net/http, encoding/json, etc.). Use variables for every configurable API parameter instead of hardcoding them.

Requirements:

- Declare variables for:
  - apiKey (string) – the xAI API key.
  - baseURL (string) – default '<https://api.x.ai/v1>'.
  - endpointPath (string) – default '/chat/completions'.
  - modelName (string) – for example 'grok-4' or any caller-supplied Grok model name.
  - temperature (float64).
  - topP (float64).
  - maxTokens (int).
  - presencePenalty (float64).
  - frequencyPenalty (float64).
  - userPrompt (string) – the user content to send.
  - systemPrompt (string) – an optional system instruction string.
- Build the full URL as:
  baseURL + endpointPath
- Use HTTP Bearer authentication for the API key by setting the header:
  'Authorization: Bearer ' + apiKey
- Construct the JSON request body to match the Grok chat completions schema (OpenAI-compatible style):
  {
    'model': modelName,
    'messages': [
      {
        'role': 'system',
        'content': systemPrompt
      },
      {
        'role': 'user',
        'content': userPrompt
      }
    ],
    'temperature': temperature,
    'top_p': topP,
    'max_tokens': maxTokens,
    'presence_penalty': presencePenalty,
    'frequency_penalty': frequencyPenalty
  }
- Define Go structs that map to this JSON shape and marshal them with encoding/json.
- Use http.NewRequest with method POST, set:
  - URL to the constructed endpoint URL.
  - Headers:
    - 'Content-Type': 'application/json'
    - 'Authorization': 'Bearer ' + apiKey
- Execute the request with http.Client.Do, handle any errors, and:
  - If status code is not 2xx, print the status and response body as an error.
  - If status is 2xx, decode the JSON response into Go structs that cover at least:
    - choices
    - choices.message.content
  - Print the assistant reply from choices.message.content to stdout.
- Organize the code so all variable configuration (apiKey, baseURL, endpointPath, modelName, temperature, topP, maxTokens, presencePenalty, frequencyPenalty, userPrompt, systemPrompt) is grouped near the top of main() or in a separate config function, making it easy to change without editing the rest of the code.
- Include basic context usage (context.Background) with http.Client if you structure it that way, but do not introduce external dependencies.

Produce idiomatic, formatted Go code ready to paste called as a function within main.go and run with 'go run .'

Sources
[1] REST API Reference <https://docs.x.ai/docs/api-reference>
[2] API <https://x.ai/api>
[3] Migration from Other Providers - xAI API <https://docs.x.ai/docs/guides/migration>
[4] How to Get Your Grok (XAI) API Key <https://www.apideck.com/blog/how-to-get-your-grok-xai-api-key>
[5] Getting Started with xAI's Grok API: Your First AI Integration <https://lablab.ai/t/xai-beginner-tutorial>
[6] The Hitchhiker's Guide to Grok <https://docs.x.ai/docs/tutorial>
[7] xAI Grok Component | Prismatic Docs <https://prismatic.io/docs/components/xai-grok/>
[8] xAI Grok Provider <https://ai-sdk.dev/providers/ai-sdk-providers/xai>
[9] xAI · Cloudflare AI Gateway docs <https://developers.cloudflare.com/ai-gateway/usage/providers/grok/>
[10] Chat - Grok API <https://grok-api.apidog.io/chat-15796842e0>
[11] Complete Guide to xAI's Grok: API Documentation and ... <https://latenode.com/blog/ai-technology-language-models/xai-grok-grok-2-grok-3/complete-guide-to-xais-grok-api-documentation-and-implementation>
[12] Grok | API References <https://docs.console.zenlayer.com/api/aigw/dialogue-generation/xai-chat-completion>
[13] How to Use xAI Grok API A Simple Guide to Setup and Integration <https://www.aionlinecourse.com/blog/how-to-use-xai-grok-api-a-simple-guide-to-setup-and-integration>
[14] Getting Started with xAI (Grok) - TensorZero Docs <https://www.tensorzero.com/docs/integrations/model-providers/xai>
[15] Grok 4 API: A Step-by-Step Guide With Examples <https://www.datacamp.com/tutorial/grok-4-api>
[16] Como acessar a API do Grok 4 - Todos os modelos de IA em uma API <https://www.cometapi.com/pt/how-to-access-grok-4-api/>
[17] xAI API <https://docs.x.ai/docs/overview>
[18] xAI Grok Connector - What is Flow Builder <https://flowbuilder-lansweeper.document360.io/docs/xai-grok-connector>
[19] xAI REST API Review | Zuplo Learning Center <https://zuplo.com/learning-center/xai-rest-api-review>
[20] grok-3-beta <https://docs.aimlapi.com/api-references/text-models-llm/xai/grok-3-beta>

## Article Prompt

"You are an expert journalist and front-end designer. Write a complete, self-contained HTML file for a journalistic article that looks polished, modern, and minimalist.  

Requirements:

1. The HTML, CSS, and JavaScript must all be included in a single file (no external dependencies or links).  
2. Use semantic HTML5 structure (header, main, article, section, footer).  
3. The design should follow minimalist principles:
   - Clean typography (prefer sans-serif font like “Inter”, “Lato”, or “Roboto”).  
   - Generous white space and consistent line spacing.  
   - Subtle color palette (light background, dark text, limited accent colors).  
   - Responsive layout that works on desktop and mobile.  
4. Include an attention-grabbing article title, byline (author name, date), and compelling lead paragraph.  
5. The body of the article should read like a professional news or feature story covering a recent topic (technology, science, culture, or world news).  
6. Use short paragraphs, subheadings, and pull quotes to improve readability.  
7. Include subtle interactive or animated elements—such as a fade-in effect for paragraphs or smooth scroll for anchor links—using only lightweight, embedded JavaScript.  
8. Add internal CSS styling (using `<style>` tags) and internal JS (using `<script>` tags). Do not rely on external sources (CDNs, frameworks, or fonts).  
9. End the HTML file cleanly with no incomplete tags.  

Output only the complete HTML document.
"

## Output

Once the HTML document is create the application must launch a web broswer and display the HTML document

## Addition features

- The application must print the following reader when launched:

"
RabbitAI, a plataforma jornalistica do Portal O Tempo.
Versão 0.20-beta-0981271.
Criada por AI Workers da plataforma Cordya AI MaxwellAI
--------------------------------------------------------

"

- The application must print to the friendly messages of all its steps it takes.
- The application must print a friendly message explaining how it should be called in the case of no parameter is provided.
- The application must exit gracefully when it finishes, printing the message "Another mission accomplished by Humans and AI. Welcome to the era of AI Citizens."
- The application must exit gracefully when an error occurs, printing a friendly message of the error, followed by a new line and the message "Errar é Humano. I am sorry!"
- The application must be FULLY implemented in main.go
