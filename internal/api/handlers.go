package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"pelican-gallery/internal/config"
	"pelican-gallery/internal/database"
	"pelican-gallery/internal/models"
)

// Handler contains the API handlers
// Handler contains the API handlers
type Handler struct {
	promptConfig *models.PromptConfig
	db           *database.DB
	tmpl         *template.Template
}

// NewHandler creates a new API handler
// NewHandler creates a new API handler
func NewHandler(promptConfig *models.PromptConfig, db *database.DB, tmpl *template.Template) *Handler {
	return &Handler{
		promptConfig: promptConfig,
		db:           db,
		tmpl:         tmpl,
	}
}

// generateSlugFromTitle creates a URL-friendly slug from a title
func generateSlugFromTitle(title string) string {
	slug := strings.ToLower(title)

	re := regexp.MustCompile(`[^a-z0-9\s-]`)
	slug = re.ReplaceAllString(slug, "")

	re = regexp.MustCompile(`[\s-]+`)
	slug = re.ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")

	if len(slug) > 50 {
		slug = slug[:50]
		slug = strings.Trim(slug, "-")
	}

	return slug
}

// isEditingEnabled checks if artwork editing/creating is enabled
// isEditingEnabled checks if artwork editing/creating is enabled
func isEditingEnabled() bool {
	return config.IsEditingEnabled()
}

// GenerateHandler handles SVG generation requests
func (h *Handler) GenerateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !isEditingEnabled() {
		log.Printf("Generate API access denied: editing is disabled")
		http.Error(w, "Artwork creation is currently disabled", http.StatusForbidden)
		return
	}

	var req models.GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Generate SVG request: model=%s, prompt length=%d", req.Model, len(req.Prompt))

	// Basic validation - ensure required fields are not empty
	if req.Prompt == "" {
		log.Printf("Error: Empty prompt provided")
		http.Error(w, "Prompt is required", http.StatusBadRequest)
		return
	}

	if req.Model == "" {
		log.Printf("Error: Model is required")
		http.Error(w, "Model is required", http.StatusBadRequest)
		return
	}

	// Set defaults if not provided
	if req.Temperature == 0 {
		req.Temperature = config.GetDefaultTemperature()
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = config.GetDefaultMaxTokens()
	}

	svg, err := h.generateSVG(req.Prompt, req.Model, req.Temperature, req.MaxTokens, req.Reasoning)
	if err != nil {
		log.Printf("Error generating SVG: %v", err)
		resp := models.GenerateResponse{
			Error: err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	log.Printf("Successfully generated SVG with length: %d characters", len(svg))

	resp := models.GenerateResponse{
		SVG: svg,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// generateSVG calls the OpenRouter API to generate SVG
func (h *Handler) generateSVG(prompt, model string, temperature float64, maxTokens int, reasoning *models.Reasoning) (string, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENROUTER_API_KEY environment variable is not set")
	}

	log.Printf("Calling OpenRouter API with model: %s", model)

	var messages []models.Message

	for _, sysPrompt := range h.promptConfig.SystemPrompts {
		messages = append(messages, models.Message{
			Role:    sysPrompt.Role,
			Content: sysPrompt.Content,
		})
	}

	userPrompt := config.FormatUserPrompt(h.promptConfig.UserPromptTemplate, prompt)
	messages = append(messages, models.Message{
		Role:    "user",
		Content: userPrompt,
	})

	log.Printf("Sending %d messages to OpenRouter", len(messages))

	openRouterReq := models.OpenRouterRequest{
		Model:       model,
		Messages:    messages,
		Temperature: temperature,
		MaxTokens:   maxTokens,
		Reasoning:   reasoning,
	}

	jsonData, err := json.Marshal(openRouterReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Printf("Request payload size: %d bytes", len(jsonData))

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("HTTP-Referer", "http://localhost:8080")
	req.Header.Set("X-Title", "Pelican Art Gallery")

	client := &http.Client{}
	log.Printf("Making request to OpenRouter API...")
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("OpenRouter API responded with status: %d", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	log.Printf("Response body size: %d bytes", len(body))

	if resp.StatusCode != http.StatusOK {
		log.Printf("OpenRouter API error response: %s", string(body))
		return "", fmt.Errorf("OpenRouter API returned status %d: %s", resp.StatusCode, string(body))
	}

	var openRouterResp models.OpenRouterResponse
	if err := json.Unmarshal(body, &openRouterResp); err != nil {
		log.Printf("Failed to parse OpenRouter response: %s", string(body))
		return "", fmt.Errorf("failed to parse response: %w (body: %s)", err, string(body))
	}

	if openRouterResp.Error != nil {
		log.Printf("OpenRouter API error: %+v", openRouterResp.Error)
		return "", fmt.Errorf("OpenRouter API error: %s", openRouterResp.Error.Message)
	}

	if len(openRouterResp.Choices) == 0 {
		log.Printf("No choices in OpenRouter response")
		return "", fmt.Errorf("no response from OpenRouter API")
	}

	log.Printf("Received %d choices from OpenRouter", len(openRouterResp.Choices))

	svgContent := strings.TrimSpace(openRouterResp.Choices[0].Message.Content)

	svgContent = strings.ReplaceAll(svgContent, "```svg", "")
	svgContent = strings.ReplaceAll(svgContent, "```", "")
	svgContent = strings.TrimSpace(svgContent)

	log.Printf("Cleaned SVG content length: %d characters", len(svgContent))

	return svgContent, nil
}

// SaveArtworkHandler handles requests to save artwork to the database
func (h *Handler) SaveArtworkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !isEditingEnabled() {
		log.Printf("Save artwork API access denied: editing is disabled")
		http.Error(w, "Artwork creation is currently disabled", http.StatusForbidden)
		return
	}

	var req models.SaveArtworkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding save artwork request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Save artwork request: title=%s, model=%s", req.Title, req.Model)

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	if req.SVGContent == "" {
		http.Error(w, "SVG content is required", http.StatusBadRequest)
		return
	}

	slug := generateSlugFromTitle(req.Title)
	if req.Slug != "" {
		slug = req.Slug
	}

	artwork := models.Artwork{
		Title:       req.Title,
		Slug:        slug,
		Category:    req.Category,
		Prompt:      req.Prompt,
		Model:       req.Model,
		SVGContent:  req.SVGContent,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		CreatedAt:   time.Now(),
	}

	log.Printf("Attempting to save artwork: slug=%s, model=%s, title=%s", artwork.Slug, artwork.Model, artwork.Title)
	if err := h.db.SaveArtwork(artwork); err != nil {
		log.Printf("Failed to save artwork to database: %v", err)
		resp := models.SaveArtworkResponse{
			Error: "Failed to save artwork to database",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	log.Printf("Artwork saved successfully for slug: %s, model: %s", artwork.Slug, artwork.Model)

	resp := models.SaveArtworkResponse{
		ID:      artwork.ID,
		Message: "Artwork saved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// RegenerateArtworkHandler handles requests to regenerate an existing artwork
func (h *Handler) RegenerateArtworkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !isEditingEnabled() {
		log.Printf("Regenerate artwork API access denied: editing is disabled")
		http.Error(w, "Artwork creation is currently disabled", http.StatusForbidden)
		return
	}

	var req struct {
		Slug  string `json:"slug"`
		Model string `json:"model"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding regenerate request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Slug == "" {
		http.Error(w, "Slug is required", http.StatusBadRequest)
		return
	}

	if req.Model == "" {
		http.Error(w, "Model is required", http.StatusBadRequest)
		return
	}

	artworks, err := h.db.GetArtworksBySlug(req.Slug)
	if err != nil {
		log.Printf("Error fetching artworks for slug %s: %v", req.Slug, err)
		http.Error(w, "Failed to fetch artwork", http.StatusInternalServerError)
		return
	}

	var existingArtwork *models.Artwork
	for _, artwork := range artworks {
		if artwork.Model == req.Model {
			existingArtwork = &artwork
			break
		}
	}

	if existingArtwork == nil {
		http.Error(w, "Artwork not found", http.StatusNotFound)
		return
	}

	svg, err := h.generateSVG(existingArtwork.Prompt, existingArtwork.Model, existingArtwork.Temperature, existingArtwork.MaxTokens, nil)
	if err != nil {
		log.Printf("Error generating artwork: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedArtwork := models.Artwork{
		Slug:        existingArtwork.Slug,
		Category:    existingArtwork.Category,
		Prompt:      existingArtwork.Prompt,
		Model:       existingArtwork.Model,
		SVGContent:  svg,
		Temperature: existingArtwork.Temperature,
		MaxTokens:   existingArtwork.MaxTokens,
		CreatedAt:   time.Now(),
	}

	if err := h.db.SaveArtwork(updatedArtwork); err != nil {
		log.Printf("Failed to save regenerated artwork: %v", err)
		http.Error(w, "Failed to save regenerated artwork", http.StatusInternalServerError)
		return
	}

	log.Printf("Artwork regenerated successfully for slug: %s, model: %s", req.Slug, req.Model)

	resp := struct {
		SVGContent string `json:"svgContent"`
		Message    string `json:"message"`
	}{
		SVGContent: svg,
		Message:    "Artwork regenerated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DeleteArtworkHandler handles artwork deletion requests
func (h *Handler) DeleteArtworkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !isEditingEnabled() {
		log.Printf("Delete artwork API access denied: editing is disabled")
		http.Error(w, "Artwork editing is currently disabled", http.StatusForbidden)
		return
	}

	var req struct {
		ID int `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Delete artwork request: ID=%d", req.ID)

	if err := h.db.DeleteArtwork(req.ID); err != nil {
		log.Printf("Error deleting artwork: %v", err)
		http.Error(w, "Failed to delete artwork", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully deleted artwork with ID: %d", req.ID)

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success": true,
		"message": "Artwork deleted successfully",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
