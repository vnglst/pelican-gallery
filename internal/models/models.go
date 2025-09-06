package models

import "time"

// PromptConfig represents the YAML configuration for the LLM prompts
type PromptConfig struct {
	Name               string         `yaml:"name"`
	Description        string         `yaml:"description"`
	SystemPrompts      []SystemPrompt `yaml:"system_prompts"`
	UserPromptTemplate string         `yaml:"user_prompt_template"`
}

// SystemPrompt represents a system prompt with role and content
type SystemPrompt struct {
	Role    string `yaml:"role"`
	Content string `yaml:"content"`
}

// Artwork represents a stored artwork in the database
type Artwork struct {
	ID          int       `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Slug        string    `db:"slug" json:"slug"`
	Category    string    `db:"category" json:"category"`
	Prompt      string    `db:"prompt" json:"prompt"`
	Model       string    `db:"model" json:"model"`
	SVGContent  string    `db:"svg_content" json:"svg_content"`
	Temperature float64   `db:"temperature" json:"temperature"`
	MaxTokens   int       `db:"max_tokens" json:"max_tokens"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

// GenerateRequest represents the request for generating SVG
type GenerateRequest struct {
	Title       string     `json:"title,omitempty"`
	Prompt      string     `json:"prompt"`
	Model       string     `json:"model"`
	Slug        string     `json:"slug,omitempty"`
	Category    string     `json:"category,omitempty"`
	Temperature float64    `json:"temperature,omitempty"`
	MaxTokens   int        `json:"max_tokens,omitempty"`
	Reasoning   *Reasoning `json:"reasoning,omitempty"`
}

// Reasoning represents reasoning configuration for models that support it
type Reasoning struct {
	Enabled bool   `json:"enabled,omitempty"`
	Effort  string `json:"effort,omitempty"`  // "high", "medium", "low" for OpenAI-style
	Exclude bool   `json:"exclude,omitempty"` // Exclude reasoning tokens from response
}

// GenerateResponse represents the response with generated SVG
type GenerateResponse struct {
	SVG   string `json:"svg"`
	Error string `json:"error,omitempty"`
}

// SaveArtworkRequest represents the request for saving an artwork
type SaveArtworkRequest struct {
	Title       string  `json:"title"`
	Slug        string  `json:"slug"`
	Category    string  `json:"category"`
	Prompt      string  `json:"prompt"`
	Model       string  `json:"model"`
	SVGContent  string  `json:"svg_content"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

// SaveArtworkResponse represents the response after saving an artwork
type SaveArtworkResponse struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// ModelInfo represents information about an available model
type ModelInfo struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Checked bool    `json:"checked"`
	Cost    float64 `json:"cost"` // Cost per 1M output tokens in dollars
}

// PromptExample represents an example prompt for users
type PromptExample struct {
	Title    string `json:"title"`
	Category string `json:"category"`
	Prompt   string `json:"prompt"`
}

// TemplateData represents all the data needed to render the index template
type TemplateData struct {
	Models           []ModelInfo     `json:"models"`
	Examples         []PromptExample `json:"examples"`
	DefaultTemp      float64         `json:"default_temp"`
	DefaultMaxTokens int             `json:"default_max_tokens"`
	DefaultModels    []string        `json:"default_models"`
	ReasoningEnabled bool            `json:"reasoning_enabled"`
	ReasoningEffort  string          `json:"reasoning_effort"`
	EditingEnabled   bool            `json:"editing_enabled"`
}

// OpenRouterRequest represents the request to OpenRouter API
type OpenRouterRequest struct {
	Model       string     `json:"model"`
	Messages    []Message  `json:"messages"`
	Temperature float64    `json:"temperature"`
	MaxTokens   int        `json:"max_tokens"`
	Reasoning   *Reasoning `json:"reasoning,omitempty"`
}

// Message represents a message in the OpenRouter request
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenRouterResponse represents the response from OpenRouter API
type OpenRouterResponse struct {
	Choices []Choice         `json:"choices"`
	Error   *OpenRouterError `json:"error,omitempty"`
}

// Choice represents a choice in the OpenRouter response
type Choice struct {
	Message Message `json:"message"`
}

// OpenRouterError represents an error from OpenRouter API
type OpenRouterError struct {
	Message string      `json:"message"`
	Type    string      `json:"type"`
	Code    interface{} `json:"code"` // Can be string or number
}
