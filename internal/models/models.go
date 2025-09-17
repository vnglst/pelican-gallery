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

// ArtworkGroup represents a group of artworks with the same prompt
type ArtworkGroup struct {
	ID          int       `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Prompt      string    `db:"prompt" json:"prompt"`
	Category    string    `db:"category" json:"category"`
	OriginalURL string    `db:"original_url" json:"original_url"`
	ArtistName  string    `db:"artist_name" json:"artist_name"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// Artwork represents an individual artwork within a group
type Artwork struct {
	ID        int       `db:"id" json:"id"`
	GroupID   int       `db:"group_id" json:"group_id"`
	Model     string    `db:"model" json:"model"`
	Params    string    `db:"params_json" json:"params"` // JSON string for parameters
	SVG       string    `db:"svg" json:"svg"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Params represents the parameters for an artwork
type Params struct {
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

// GenerateRequest represents the request for generating SVG
type GenerateRequest struct {
	Title       string  `json:"title,omitempty"`
	Prompt      string  `json:"prompt"`
	Model       string  `json:"model"`
	Category    string  `json:"category,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
}

// GenerateResponse represents the response with generated SVG
type GenerateResponse struct {
	SVG   string `json:"svg"`
	Error string `json:"error,omitempty"`
}

// SaveArtworkRequest represents the request for saving an artwork
type SaveArtworkRequest struct {
	Title       string  `json:"title"`
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

// TemplateData represents all the data needed to render template
type TemplateData struct {
	Models         []ModelInfo `json:"models"`
	EditingEnabled bool        `json:"editing_enabled"`
}

// OpenRouterRequest represents the request to OpenRouter API
type OpenRouterRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens"`
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
