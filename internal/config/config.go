package config

import (
	"fmt"
	"os"
	"strings"

	"pelican-gallery/internal/models"

	"gopkg.in/yaml.v3"
)

// LoadPromptConfig loads the prompt configuration from the YAML file
func LoadPromptConfig(filename string) (*models.PromptConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config models.PromptConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// FormatUserPrompt formats the user prompt template with the provided description
func FormatUserPrompt(template, description string) string {
	return strings.ReplaceAll(template, "{art_work_description}", description)
}

// GetAvailableModels returns a list of available models for the dropdown
func GetAvailableModels() []models.ModelInfo {
	defaultModels := GetDefaultModels()
	defaultSet := make(map[string]bool)
	for _, id := range defaultModels {
		defaultSet[id] = true
	}

	allModels := getAllModels()

	// Set the Checked field based on whether the model is in defaults
	for i := range allModels {
		allModels[i].Checked = defaultSet[allModels[i].ID]
	}

	return allModels
}

// GetDefaultTemperature returns the default temperature for generation
func GetDefaultTemperature() float64 {
	return 0.5
}

// GetDefaultMaxTokens returns the default max tokens for generation
func GetDefaultMaxTokens() int {
	return 30_000
}

// GetDefaultReasoningEnabled returns the default reasoning enabled state
func GetDefaultReasoningEnabled() bool {
	return false
}

// GetDefaultReasoningEffort returns the default reasoning effort level
func GetDefaultReasoningEffort() string {
	return "medium"
}

// IsEditingEnabled checks if artwork editing/creating is enabled
func IsEditingEnabled() bool {
	// Check if editing is explicitly enabled (defaults to false if not set)
	enableEditing := os.Getenv("ENABLE_EDITING")
	if enableEditing == "" {
		return false // Default to disabled
	}
	return enableEditing == "true" || enableEditing == "1"
}

// GetDefaultModels returns the default model IDs
func GetDefaultModels() []string {
	// Get all available models and filter for free ones or those under $0.40/1M tokens
	allModels := getAllModels() // Helper function to get the raw model data
	var defaultModelIDs []string

	for _, model := range allModels {
		if model.Cost == 0.00 || model.Cost < 0.20 {
			defaultModelIDs = append(defaultModelIDs, model.ID)
		}
	}

	return defaultModelIDs
}

// getAllModels returns the raw model data (helper function)
func getAllModels() []models.ModelInfo {
	return []models.ModelInfo{
		{ID: "deepseek/deepseek-chat-v3.1:free", Name: "DeepSeek: DeepSeek V3.1 (free)", Cost: 0.00},
		{ID: "openai/gpt-5-nano", Name: "OpenAI: GPT-5 Nano", Cost: 0.05},
		{ID: "meta-llama/llama-3.3-70b-instruct", Name: "Meta: Llama 3.3 70B Instruct", Cost: 0.12},
		{ID: "google/gemma-3-12b-it", Name: "Google: Gemma 3 12B", Cost: 0.193},
		{ID: "x-ai/grok-code-fast-1", Name: "xAI: Grok Code Fast 1", Cost: 0.20},
		{ID: "google/gemma-3-27b-it", Name: "Google: Gemma 3 27B", Cost: 0.267},
		{ID: "openai/gpt-oss-120b", Name: "OpenAI: gpt-oss-120b", Cost: 0.28},
		{ID: "google/gemini-2.0-flash-001", Name: "Google: Gemini 2.0 Flash", Cost: 0.40},
		{ID: "google/gemini-2.5-flash-lite", Name: "Google: Gemini 2.5 Flash Lite", Cost: 0.40},
		{ID: "google/gemini-2.5-flash-lite-preview-06-17", Name: "Google: Gemini 2.5 Flash Lite Preview 06-17", Cost: 0.40},
		{ID: "anthropic/claude-3.5-haiku", Name: "Anthropic: Claude 3.5 Haiku", Cost: 0.80},
		{ID: "deepseek/deepseek-chat-v3-0324", Name: "DeepSeek: DeepSeek V3 0324", Cost: 0.80},
		{ID: "deepseek/deepseek-chat-v3.1", Name: "DeepSeek: DeepSeek V3.1", Cost: 0.80},
		{ID: "qwen/qwen3-coder", Name: "Qwen: Qwen3 Coder 480B A35B", Cost: 0.80},
		{ID: "deepseek/deepseek-r1-0528", Name: "DeepSeek: R1 0528", Cost: 0.80},
		{ID: "openai/gpt-5", Name: "OpenAI: GPT-5", Cost: 1.25},
		{ID: "z-ai/glm-4.5", Name: "Z.AI: GLM 4.5", Cost: 1.32},
		{ID: "moonshotai/kimi-k2", Name: "MoonshotAI: Kimi K2", Cost: 2.49},
		{ID: "google/gemini-2.5-flash", Name: "Google: Gemini 2.5 Flash", Cost: 2.50},
		{ID: "anthropic/claude-3.7-sonnet", Name: "Anthropic: Claude 3.7 Sonnet", Cost: 3.00},
		{ID: "x-ai/grok-4", Name: "xAI: Grok 4", Cost: 3.00},
		{ID: "anthropic/claude-sonnet-4", Name: "Anthropic: Claude Sonnet 4", Cost: 3.00},
		{ID: "anthropic/claude-opus-4", Name: "Anthropic: Claude Opus 4", Cost: 15.00},
		{ID: "anthropic/claude-opus-4.1", Name: "Anthropic: Claude Opus 4.1", Cost: 15.00},
	}
}

// GetExamplePrompts returns a list of example prompts for users
func GetExamplePrompts() []models.PromptExample {
	return []models.PromptExample{
		{
			Title:    "The Night Watch",
			Category: "classic",
			Prompt:   "The Night Watch. Rembrandt. A group portrait of Amsterdam militiamen led by Captain Frans Banning Cocq and Lieutenant Willem van Ruytenburch, stepping out from a gateway. The scene includes officers, a drummer, a girl with a chicken, and a barking dog, set in dramatic lighting.",
		},
		{
			Title:    "Het Melkmeisje",
			Category: "classic",
			Prompt:   "Het Melkmeisje. Johannes Vermeer. A kitchen maid stands behind a table, pouring milk from a jug into a cooking pot. On the table are a basket of bread and a stone jug; to the left, a window with a wicker basket and a copper pot. At the bottom right, a row of tiles and a small stove on the wall.",
		},
		{
			Title:    "Dutch Maritime Scene",
			Category: "classic",
			Prompt:   "A fleet of ships at anchor, greeting a government barge with salutes. Calm water reflects the boats and a dramatic cloudscape fills the sky.",
		},
		{
			Title:    "Geometric Landscape",
			Category: "modern",
			Prompt:   "A mountain landscape with angular peaks, geometric patterns, and a stylized sun. Bold colors and sharp shapes.",
		},
		{
			Title:    "Art Deco Architecture",
			Category: "modern",
			Prompt:   "An Art Deco skyscraper with symmetrical facade, sunburst patterns, zigzags, and metallic colors.",
		},
		{
			Title:    "Botanical Illustration",
			Category: "nature",
			Prompt:   "Detailed illustration of exotic flowering plants with precise linework, leaf structures, and natural colors.",
		},
		{
			Title:    "Abstract Composition",
			Category: "abstract",
			Prompt:   "Flowing organic shapes and curves, overlapping forms, limited color palette, and transparent elements.",
		},
		{
			Title:  "Minimalist Icon",
			Prompt: "A minimalist icon of growth, using simple geometric shapes like an upward arrow or sprouting plant.",
		},
	}
}
