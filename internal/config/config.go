package config

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"pelican-gallery/internal/models"

	"gopkg.in/yaml.v3"
)

var (
	modelsCache []models.ModelInfo
	cacheExpiry time.Time
	modelsMu    sync.RWMutex
)

type openRouterResponse struct {
	Data []openRouterModel `json:"data"`
}

type openRouterModel struct {
	ID      string                 `json:"id"`
	Name    string                 `json:"name"`
	Pricing map[string]interface{} `json:"pricing"`
}

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

	// Try to fetch live models from OpenRouter when an API key is present.
	var allModels []models.ModelInfo
	if openModels, err := fetchOpenRouterModels(); err == nil && len(openModels) > 0 {
		allModels = openModels
	} else {
		allModels = getAllModels()
	}

	// Sort models by cost (cheapest first)
	sort.Slice(allModels, func(i, j int) bool {
		return allModels[i].Cost < allModels[j].Cost
	})

	// Filter out the "openrouter/auto" model
	var filteredModels []models.ModelInfo
	for _, model := range allModels {
		if model.ID != "openrouter/auto" {
			filteredModels = append(filteredModels, model)
		}
	}

	// Set the Checked field based on whether the model is in defaults
	for i := range filteredModels {
		filteredModels[i].Checked = defaultSet[filteredModels[i].ID]
	}

	return filteredModels
}

// fetchOpenRouterModels fetches models from the OpenRouter API
func fetchOpenRouterModels() ([]models.ModelInfo, error) {
	// Return cached value if valid
	modelsMu.RLock()
	if time.Now().Before(cacheExpiry) && len(modelsCache) > 0 {
		models := make([]models.ModelInfo, len(modelsCache))
		copy(models, modelsCache)
		modelsMu.RUnlock()
		return models, nil
	}
	modelsMu.RUnlock()

	// Fetch from API
	modelsMu.Lock()
	defer modelsMu.Unlock()

	resp, err := http.Get("https://openrouter.ai/api/v1/models")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var apiResp openRouterResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	var modelInfos []models.ModelInfo
	for _, model := range apiResp.Data {
		cost := 0.0
		if completion, ok := model.Pricing["completion"].(string); ok {
			if f, err := parseFloat(completion); err == nil {
				// Convert from per-token to per-million-tokens cost
				cost = f * 1000000
			}
		}
		modelInfos = append(modelInfos, models.ModelInfo{
			ID:   model.ID,
			Name: model.Name,
			Cost: cost,
		})
	}

	// Update cache
	modelsCache = make([]models.ModelInfo, len(modelInfos))
	copy(modelsCache, modelInfos)
	cacheExpiry = time.Now().Add(5 * time.Minute)

	log.Printf("Fetched %d models from OpenRouter", len(modelInfos))
	return modelInfos, nil
}

// parseFloat parses a string to float64
func parseFloat(s string) (float64, error) {
	if s == "" {
		return 0, fmt.Errorf("empty string")
	}
	return strconv.ParseFloat(s, 64)
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
	// Return live models from OpenRouter only. If the API call fails return an
	// empty list.
	if live, err := fetchOpenRouterModels(); err == nil && len(live) > 0 {
		return live
	}

	return []models.ModelInfo{}
}
