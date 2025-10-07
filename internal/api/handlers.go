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
	"strconv"
	"strings"
	"time"

	"pelican-gallery/internal/config"
	"pelican-gallery/internal/database"
	"pelican-gallery/internal/models"
)

// Handler contains the API handlers
type Handler struct {
	promptConfig *models.PromptConfig
	db           *database.DB
	tmpl         *template.Template
}

// NewHandler creates a new API handler
func NewHandler(promptConfig *models.PromptConfig, db *database.DB, tmpl *template.Template) *Handler {
	return &Handler{
		promptConfig: promptConfig,
		db:           db,
		tmpl:         tmpl,
	}
}

// jsonError is a simple structured error returned to clients
type jsonError struct {
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

func writeJSONError(w http.ResponseWriter, status int, message string, details ...interface{}) {
	var det interface{}
	if len(details) > 0 {
		det = details[0]
	}
	writeJSON(w, status, jsonError{Message: message, Details: det})
}

// isEditingEnabled checks if artwork editing/creating is enabled
func isEditingEnabled() bool {
	return config.IsEditingEnabled()
}

// GenerateHandler handles SVG generation requests
func (h *Handler) GenerateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	if !isEditingEnabled() {
		log.Printf("Generate API access denied: editing is disabled")
		writeJSONError(w, http.StatusForbidden, "Artwork creation is currently disabled")
		return
	}

	var req models.GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding generate request body: %v", err)
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Prompt == "" {
		writeJSONError(w, http.StatusBadRequest, "Prompt is required")
		return
	}

	if req.Model == "" {
		writeJSONError(w, http.StatusBadRequest, "Model is required")
		return
	}

	if req.Temperature < 0 || req.Temperature > 1 {
		writeJSONError(w, http.StatusBadRequest, "Temperature must be between 0 and 1")
		return
	}

	if req.MaxTokens <= 0 {
		writeJSONError(w, http.StatusBadRequest, "MaxTokens must be positive")
		return
	}

	log.Printf("Generate SVG request: model=%s, prompt length=%d", req.Model, len(req.Prompt))

	svg, err := h.generateSVG(req.Prompt, req.Model, req.Temperature, req.MaxTokens)
	if err != nil {
		log.Printf("Error generating SVG: %v", err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Successfully generated SVG with length: %d characters", len(svg))

	resp := models.GenerateResponse{
		SVG: svg,
	}

	writeJSON(w, http.StatusOK, resp)
}

// generateSVG calls the OpenRouter API to generate SVG
func (h *Handler) generateSVG(prompt, model string, temperature float64, maxTokens int) (string, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENROUTER_API_KEY environment variable is not set")
	}

	log.Printf("Calling OpenRouter API with model: %s", model)

	var messages []models.Message

	for _, sysPrompt := range h.promptConfig.SystemPrompts {
		messages = append(messages, models.Message(sysPrompt))
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
		Reasoning: &models.Reasoning{
			Effort:  "medium",
			Enabled: true,
			Exclude: true,
		},
	}

	// Note: reasoning is enabled for supported models at medium effort.
	// We exclude reasoning from the response (exclude=true) and do not log reasoning content.
	log.Printf("Request will use reasoning: effort=%s, exclude=%t", "medium", true)

	jsonData, err := json.Marshal(openRouterReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))

	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("X-Title", "Pelican Art Gallery")

	client := &http.Client{
		Timeout: 300 * time.Second, // 5 minutes
	}
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

	if resp.StatusCode != http.StatusOK {
		log.Printf("OpenRouter API error (status %d): %s", resp.StatusCode, string(body))
		return "", fmt.Errorf("OpenRouter API returned status %d: %s", resp.StatusCode, string(body))
	}

	var openRouterResp models.OpenRouterResponse
	if err := json.Unmarshal(body, &openRouterResp); err != nil {
		log.Printf("Failed to parse OpenRouter response")
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if openRouterResp.Error != nil {
		log.Printf("OpenRouter API error: %s", openRouterResp.Error.Message)
		return "", fmt.Errorf("OpenRouter API error: %s", openRouterResp.Error.Message)
	}

	if len(openRouterResp.Choices) == 0 {
		log.Printf("No choices in OpenRouter response")
		return "", fmt.Errorf("no response from OpenRouter API")
	}

	log.Printf("Received %d choices from OpenRouter", len(openRouterResp.Choices))

	svgContent := strings.TrimSpace(openRouterResp.Choices[0].Message.Content)
	log.Printf("Raw OpenRouter response content length: %d", len(svgContent))

	return svgContent, nil
}

// DeleteArtworkHandler handles artwork deletion requests
func (h *Handler) DeleteArtworkHandler(w http.ResponseWriter, r *http.Request, artworkIDStr string) {
	if !isEditingEnabled() {
		writeJSONError(w, http.StatusForbidden, "Artwork editing is currently disabled")
		return
	}

	artworkID, err := strconv.Atoi(artworkIDStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid artwork ID")
		return
	}

	log.Printf("Delete artwork request: ID=%d", artworkID)

	if err := h.db.DeleteArtwork(artworkID); err != nil {
		log.Printf("Error deleting artwork (id=%d): %v", artworkID, err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to delete artwork")
		return
	}

	log.Printf("Successfully deleted artwork with ID: %d", artworkID)

	response := map[string]interface{}{
		"success": true,
		"message": "Artwork deleted successfully",
	}
	writeJSON(w, http.StatusOK, response)
}

// ListGroupsHandler handles GET /api/groups
func (h *Handler) ListGroupsHandler(w http.ResponseWriter, r *http.Request) {
	groups, err := h.db.ListGroups()
	if err != nil {
		log.Printf("Error listing groups: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to list groups")
		return
	}
	writeJSON(w, http.StatusOK, groups)
}

// CreateGroupHandler handles POST /api/groups
func (h *Handler) CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	if !isEditingEnabled() {
		writeJSONError(w, http.StatusForbidden, "Artwork creation is currently disabled")
		return
	}

	var req struct {
		Title       string `json:"title"`
		Prompt      string `json:"prompt"`
		Category    string `json:"category"`
		OriginalURL string `json:"original_url"`
		ArtistName  string `json:"artist_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("CreateGroup invalid body: %v", err)
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" || req.Prompt == "" {
		writeJSONError(w, http.StatusBadRequest, "Title and prompt are required")
		return
	}

	group := models.ArtworkGroup{
		Title:       req.Title,
		Prompt:      req.Prompt,
		Category:    req.Category,
		OriginalURL: req.OriginalURL,
		ArtistName:  req.ArtistName,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	id, err := h.db.CreateGroup(group)
	if err != nil {
		log.Printf("Error creating group: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to create group")
		return
	}

	group.ID = id
	writeJSON(w, http.StatusCreated, group)
}

// UpdateGroupHandler handles PUT /api/groups/{id}
func (h *Handler) UpdateGroupHandler(w http.ResponseWriter, r *http.Request, groupIDStr string) {
	if !isEditingEnabled() {
		writeJSONError(w, http.StatusForbidden, "Artwork editing is currently disabled")
		return
	}

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	var req struct {
		Title       string `json:"title"`
		Prompt      string `json:"prompt"`
		Category    string `json:"category"`
		OriginalURL string `json:"original_url"`
		ArtistName  string `json:"artist_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("UpdateGroup invalid body: %v", err)
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" || req.Prompt == "" {
		writeJSONError(w, http.StatusBadRequest, "Title and prompt are required")
		return
	}

	group := models.ArtworkGroup{
		ID:          groupID,
		Title:       req.Title,
		Prompt:      req.Prompt,
		Category:    req.Category,
		OriginalURL: req.OriginalURL,
		ArtistName:  req.ArtistName,
		UpdatedAt:   time.Now(),
	}

	if err := h.db.UpdateGroup(group); err != nil {
		log.Printf("Error updating group (id=%d): %v", groupID, err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to update group")
		return
	}

	writeJSON(w, http.StatusOK, group)
}

// DeleteGroupHandler handles DELETE /api/groups/{id}
func (h *Handler) DeleteGroupHandler(w http.ResponseWriter, r *http.Request, groupIDStr string) {
	if !isEditingEnabled() {
		writeJSONError(w, http.StatusForbidden, "Artwork editing is currently disabled")
		return
	}

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	log.Printf("Delete group request: ID=%d", groupID)

	if err := h.db.DeleteGroup(groupID); err != nil {
		log.Printf("Error deleting group (id=%d): %v", groupID, err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to delete group")
		return
	}

	log.Printf("Successfully deleted group with ID: %d (cascaded to all artworks)", groupID)

	response := map[string]interface{}{
		"success": true,
		"message": "Group and all associated artworks deleted successfully",
	}
	writeJSON(w, http.StatusOK, response)
}

// GetGroupHandler handles GET /api/groups/{id}
func (h *Handler) GetGroupHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/groups/")
	idStr := strings.TrimSuffix(path, "/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	group, err := h.db.GetGroup(id)
	if err != nil {
		log.Printf("Error getting group: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to get group")
		return
	}

	artworks, err := h.db.ListArtworksByGroup(id)
	if err != nil {
		log.Printf("Error listing artworks: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to list artworks")
		return
	}

	response := struct {
		Group    *models.ArtworkGroup `json:"group"`
		Artworks []models.Artwork     `json:"artworks"`
	}{
		Group:    group,
		Artworks: artworks,
	}

	writeJSON(w, http.StatusOK, response)
}

// CreateArtworkHandler handles POST /api/artworks
func (h *Handler) CreateArtworkHandler(w http.ResponseWriter, r *http.Request) {
	if !isEditingEnabled() {
		writeJSONError(w, http.StatusForbidden, "Artwork creation is currently disabled")
		return
	}

	var req struct {
		GroupID     int     `json:"group_id"`
		Model       string  `json:"model"`
		Temperature float64 `json:"temperature"`
		MaxTokens   int     `json:"max_tokens"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("CreateArtwork invalid body: %v", err)
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.GroupID == 0 || req.Model == "" {
		writeJSONError(w, http.StatusBadRequest, "Group ID and model are required")
		return
	}

	artwork := models.Artwork{
		GroupID:     req.GroupID,
		Model:       req.Model,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	id, err := h.db.CreateArtwork(artwork)
	if err != nil {
		log.Printf("Error creating artwork (group_id=%d, model=%s): %v", req.GroupID, req.Model, err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to create artwork")
		return
	}

	artwork.ID = id
	writeJSON(w, http.StatusCreated, artwork)
}

// UpdateArtworkHandler handles PATCH /api/artworks/{id}
func (h *Handler) UpdateArtworkHandler(w http.ResponseWriter, r *http.Request, artworkIDStr string) {
	if !isEditingEnabled() {
		writeJSONError(w, http.StatusForbidden, "Artwork editing is currently disabled")
		return
	}

	artworkID, err := strconv.Atoi(artworkIDStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid artwork ID")
		return
	}

	var req struct {
		Temperature float64 `json:"temperature"`
		MaxTokens   int     `json:"max_tokens"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("UpdateArtwork invalid body: %v", err)
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.db.UpdateArtwork(artworkID, req.Temperature, req.MaxTokens); err != nil {
		log.Printf("Error updating artwork (id=%d): %v", artworkID, err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to update artwork")
		return
	}

	artwork, err := h.db.GetArtwork(artworkID)
	if err != nil {
		log.Printf("Error getting updated artwork (id=%d): %v", artworkID, err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to get updated artwork")
		return
	}

	writeJSON(w, http.StatusOK, artwork)
}

// GenerateArtworkHandler handles POST /api/generate
func (h *Handler) GenerateArtworkHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	if !isEditingEnabled() {
		writeJSONError(w, http.StatusForbidden, "Artwork creation is currently disabled")
		return
	}

	var req struct {
		ArtworkID int `json:"artwork_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("GenerateArtwork invalid body: %v", err)
		writeJSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.ArtworkID == 0 {
		writeJSONError(w, http.StatusBadRequest, "Artwork ID is required")
		return
	}

	artwork, err := h.db.GetArtwork(req.ArtworkID)
	if err != nil {
		log.Printf("Error getting artwork (id=%d): %v", req.ArtworkID, err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to get artwork")
		return
	}

	group, err := h.db.GetGroup(artwork.GroupID)
	if err != nil {
		log.Printf("Error getting group (id=%d for artwork=%d): %v", artwork.GroupID, req.ArtworkID, err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to get group")
		return
	}

	svg, err := h.generateSVG(group.Prompt, artwork.Model, artwork.Temperature, artwork.MaxTokens)
	if err != nil {
		log.Printf("Error generating SVG for artwork %d: %v", req.ArtworkID, err)
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Generated SVG for artwork %d: length=%d characters", req.ArtworkID, len(svg))

	if err := h.db.SaveArtworkSVG(req.ArtworkID, svg); err != nil {
		log.Printf("Error saving SVG (artwork=%d): %v", req.ArtworkID, err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to save SVG")
		return
	}

	log.Printf("Successfully saved SVG for artwork %d to database", req.ArtworkID)

	response := struct {
		ID  int    `json:"id"`
		SVG string `json:"svg"`
	}{
		ID:  req.ArtworkID,
		SVG: svg,
	}

	writeJSON(w, http.StatusOK, response)
}

// ListModelsHandler handles GET /api/models
func (h *Handler) ListModelsHandler(w http.ResponseWriter, r *http.Request) {
	models := config.GetAvailableModels()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"models": models,
	})
}

// UploadOriginalArtworkHandler handles POST /api/groups/{id}/original-artwork
func (h *Handler) UploadOriginalArtworkHandler(w http.ResponseWriter, r *http.Request, groupIDStr string) {
	if !isEditingEnabled() {
		writeJSONError(w, http.StatusForbidden, "Artwork editing is currently disabled")
		return
	}

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	// Parse multipart form with 10MB max memory
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Printf("Error parsing multipart form: %v", err)
		writeJSONError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	file, header, err := r.FormFile("artwork")
	if err != nil {
		log.Printf("Error getting file from form: %v", err)
		writeJSONError(w, http.StatusBadRequest, "No file provided")
		return
	}
	defer file.Close()

	// Validate file type (accept common image formats)
	contentType := header.Header.Get("Content-Type")
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	if !validTypes[contentType] {
		writeJSONError(w, http.StatusBadRequest, "Invalid file type. Only images (jpeg, png, gif, webp) are allowed")
		return
	}

	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to read file")
		return
	}

	// Get the group to update
	group, err := h.db.GetGroup(groupID)
	if err != nil {
		log.Printf("Error getting group %d: %v", groupID, err)
		writeJSONError(w, http.StatusNotFound, "Group not found")
		return
	}

	// Update the group with the original artwork
	group.OriginalArtwork = fileBytes
	group.UpdatedAt = time.Now()

	if err := h.db.UpdateGroup(*group); err != nil {
		log.Printf("Error updating group %d with original artwork: %v", groupID, err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to save original artwork")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Original artwork uploaded successfully",
	})
}

// GetOriginalArtworkHandler handles GET /api/groups/{id}/original-artwork
func (h *Handler) GetOriginalArtworkHandler(w http.ResponseWriter, r *http.Request, groupIDStr string) {
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid group ID")
		return
	}

	group, err := h.db.GetGroup(groupID)
	if err != nil {
		log.Printf("Error getting group %d: %v", groupID, err)
		writeJSONError(w, http.StatusNotFound, "Group not found")
		return
	}

	if group.OriginalArtwork == nil || len(group.OriginalArtwork) == 0 {
		writeJSONError(w, http.StatusNotFound, "No original artwork found for this group")
		return
	}

	// Detect content type from the first few bytes
	contentType := http.DetectContentType(group.OriginalArtwork)

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(group.OriginalArtwork)))
	w.WriteHeader(http.StatusOK)
	w.Write(group.OriginalArtwork)
}

// SetFeaturedArtworkHandler handles POST /api/artworks/{id}/featured
func (h *Handler) SetFeaturedArtworkHandler(w http.ResponseWriter, r *http.Request, artworkIDStr string) {
	if !isEditingEnabled() {
		writeJSONError(w, http.StatusForbidden, "Artwork editing is currently disabled")
		return
	}

	artworkID, err := strconv.Atoi(artworkIDStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid artwork ID")
		return
	}

	if err := h.db.SetFeaturedArtwork(artworkID); err != nil {
		log.Printf("Error setting featured artwork %d: %v", artworkID, err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to set featured artwork")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Artwork set as featured",
	})
}
