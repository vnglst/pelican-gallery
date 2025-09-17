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
		{ID: "openai/gpt-5-mini", Name: "OpenAI: GPT-5 Mini", Cost: 2.00},
		{ID: "anthropic/claude-3.5-haiku", Name: "Anthropic: Claude 3.5 Haiku", Cost: 0.80},
		{ID: "deepseek/deepseek-chat-v3-0324", Name: "DeepSeek: DeepSeek V3 0324", Cost: 0.80},
		{ID: "deepseek/deepseek-chat-v3.1", Name: "DeepSeek: DeepSeek V3.1", Cost: 0.80},
		{ID: "qwen/qwen3-coder", Name: "Qwen: Qwen3 Coder 480B A35B", Cost: 0.80},
		{ID: "deepseek/deepseek-r1-0528", Name: "DeepSeek: R1 0528", Cost: 0.80},
		{ID: "openai/gpt-4.1", Name: "OpenAI: GPT-4.1", Cost: 8.00},
		{ID: "openai/gpt-5", Name: "OpenAI: GPT-5", Cost: 1.25},
		{ID: "openai/gpt-4o-2024-11-20", Name: "OpenAI: GPT-4o 2024-11-20", Cost: 10.00},
		{ID: "google/gemini-2.5-pro", Name: "Google: Gemini 2.5 Pro", Cost: 10.00},
		{ID: "z-ai/glm-4.5", Name: "Z.AI: GLM 4.5", Cost: 1.32},
		{ID: "moonshotai/kimi-k2", Name: "MoonshotAI: Kimi K2", Cost: 2.49},
		{ID: "google/gemini-2.5-flash", Name: "Google: Gemini 2.5 Flash", Cost: 2.50},
		{ID: "anthropic/claude-3.7-sonnet", Name: "Anthropic: Claude 3.7 Sonnet", Cost: 3.00},
		{ID: "x-ai/grok-4", Name: "xAI: Grok 4", Cost: 3.00},
		{ID: "anthropic/claude-sonnet-4", Name: "Anthropic: Claude Sonnet 4", Cost: 3.00},
		{ID: "openai/gpt-4-turbo", Name: "OpenAI: GPT-4 Turbo", Cost: 30.00},
		{ID: "anthropic/claude-opus-4", Name: "Anthropic: Claude Opus 4", Cost: 15.00},
		{ID: "anthropic/claude-opus-4.1", Name: "Anthropic: Claude Opus 4.1", Cost: 15.00},
		{ID: "openai/gpt-3.5-turbo-0613", Name: "OpenAI: GPT-3.5 Turbo (older v0613)", Cost: 1.50},
	}
}
