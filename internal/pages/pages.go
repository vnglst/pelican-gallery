package pages

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"pelican-gallery/internal/config"
	"pelican-gallery/internal/database"
	"pelican-gallery/internal/models"
)

// PageHandler contains the page handlers
// PageHandler contains the page handlers
type PageHandler struct {
	db           *database.DB
	tmpl         *template.Template
	templateData models.TemplateData
}

// NewPageHandler creates a new page handler
// NewPageHandler creates a new page handler
func NewPageHandler(db *database.DB, tmpl *template.Template, templateData models.TemplateData) *PageHandler {
	return &PageHandler{
		db:           db,
		tmpl:         tmpl,
		templateData: templateData,
	}
}

// GalleryHandler handles requests to display the gallery of saved artworks
func (h *PageHandler) GalleryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	category := r.URL.Query().Get("category")

	var artworks []models.Artwork
	var err error

	if category != "" {
		artworks, err = h.db.GetArtworksByCategory(category)
		if err != nil {
			log.Printf("Error fetching artworks by category %s: %v", category, err)
			http.Error(w, "Failed to fetch artworks", http.StatusInternalServerError)
			return
		}
		log.Printf("Fetched %d artworks for category: %s", len(artworks), category)
	} else {
		artworks, err = h.db.GetAllArtworks()
		if err != nil {
			log.Printf("Error fetching artworks: %v", err)
			http.Error(w, "Failed to fetch artworks", http.StatusInternalServerError)
			return
		}
		log.Printf("Fetched %d artworks for gallery", len(artworks))
	}

	categories, err := h.db.GetAllCategories()
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}

	log.Printf("Fetched %d artworks and %d categories for gallery", len(artworks), len(categories))

	type GalleryArtwork struct {
		ID          int           `json:"id"`
		Title       string        `json:"title"`
		Slug        string        `json:"slug"`
		Category    string        `json:"category"`
		Prompt      string        `json:"prompt"`
		Model       string        `json:"model"`
		SVGContent  template.HTML `json:"svg_content"`
		Temperature float64       `json:"temperature"`
		MaxTokens   int           `json:"max_tokens"`
		CreatedAt   time.Time     `json:"created_at"`
	}

	type ArtworkGroup struct {
		Title     string           `json:"title"`
		Slug      string           `json:"slug"`
		Category  string           `json:"category"`
		Prompt    string           `json:"prompt"`
		CreatedAt time.Time        `json:"created_at"`
		Artworks  []GalleryArtwork `json:"artworks"`
	}

	galleryArtworks := make([]GalleryArtwork, len(artworks))
	for i, artwork := range artworks {
		galleryArtworks[i] = GalleryArtwork{
			ID:          artwork.ID,
			Title:       artwork.Title,
			Slug:        artwork.Slug,
			Category:    artwork.Category,
			Prompt:      artwork.Prompt,
			Model:       artwork.Model,
			SVGContent:  template.HTML(artwork.SVGContent),
			Temperature: artwork.Temperature,
			MaxTokens:   artwork.MaxTokens,
			CreatedAt:   artwork.CreatedAt,
		}
	}

	groupMap := make(map[string]*ArtworkGroup)
	for _, artwork := range galleryArtworks {
		if group, exists := groupMap[artwork.Slug]; exists {
			group.Artworks = append(group.Artworks, artwork)
		} else {
			groupMap[artwork.Slug] = &ArtworkGroup{
				Title:     artwork.Title,
				Slug:      artwork.Slug,
				Category:  artwork.Category,
				Prompt:    artwork.Prompt,
				CreatedAt: artwork.CreatedAt,
				Artworks:  []GalleryArtwork{artwork},
			}
		}
	}

	var groups []ArtworkGroup
	for _, group := range groupMap {
		groups = append(groups, *group)
	}

	data := struct {
		Title          string         `json:"title"`
		Groups         []ArtworkGroup `json:"groups"`
		Categories     []string       `json:"categories"`
		Category       string         `json:"category"`
		EditingEnabled bool           `json:"editing_enabled"`
	}{
		Title:          "Gallery - Pelican Art Gallery",
		Groups:         groups,
		Categories:     categories,
		Category:       category,
		EditingEnabled: isEditingEnabled(),
	}

	w.Header().Set("Content-Type", "text/html")
	if err := h.tmpl.ExecuteTemplate(w, "gallery.html", data); err != nil {
		log.Printf("Error executing gallery template: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

// isEditingEnabled checks if artwork editing/creating is enabled
func isEditingEnabled() bool {
	return config.IsEditingEnabled()
}

// HomepageHandler handles requests to the homepage
func (h *PageHandler) HomepageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	homepageData := struct {
		EditingEnabled bool `json:"editing_enabled"`
	}{
		EditingEnabled: config.IsEditingEnabled(),
	}
	if err := h.tmpl.ExecuteTemplate(w, "homepage.html", homepageData); err != nil {
		log.Printf("Failed to execute homepage template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// WorkshopHandler handles requests to the workshop page
func (h *PageHandler) WorkshopHandler(w http.ResponseWriter, r *http.Request) {
	// Check if editing is enabled
	if !isEditingEnabled() {
		log.Printf("Workshop access denied: editing is disabled")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// Check if we're editing an existing artwork
	editSlug := r.URL.Query().Get("edit")
	var editData *models.Artwork
	var editArtworks []models.Artwork
	if editSlug != "" {
		log.Printf("Edit mode requested for slug: %s", editSlug)
		// Get all artworks with this slug for editing
		artworks, err := h.db.GetArtworksBySlug(editSlug)
		if err != nil {
			log.Printf("Error fetching artwork for editing: %v", err)
		} else if len(artworks) > 0 {
			editData = &artworks[0] // Use first artwork for form data
			editArtworks = artworks // All artworks for display
			log.Printf("Found %d artwork(s) for editing: %s (category: %s)", len(artworks), editData.Slug, editData.Category)
		} else {
			log.Printf("No artwork found with slug: %s", editSlug)
		}
	}

	// Create safe HTML versions for display
	type SafeArtwork struct {
		ID          int
		Title       string
		Slug        string
		Category    string
		Prompt      string
		Model       string
		SVGContent  template.HTML
		Temperature float64
		MaxTokens   int
		CreatedAt   time.Time
	}

	var safeEditArtworks []SafeArtwork
	for _, artwork := range editArtworks {
		safeEditArtworks = append(safeEditArtworks, SafeArtwork{
			ID:          artwork.ID,
			Title:       artwork.Title,
			Slug:        artwork.Slug,
			Category:    artwork.Category,
			Prompt:      artwork.Prompt,
			Model:       artwork.Model,
			SVGContent:  template.HTML(artwork.SVGContent),
			Temperature: artwork.Temperature,
			MaxTokens:   artwork.MaxTokens,
			CreatedAt:   artwork.CreatedAt,
		})
	}

	// Prepare template data
	templateData := h.templateData

	// Create template data with edit information
	currentTemplateData := struct {
		Models           []models.ModelInfo     `json:"models"`
		Examples         []models.PromptExample `json:"examples"`
		DefaultTemp      float64                `json:"default_temp"`
		DefaultMaxTokens int                    `json:"default_max_tokens"`
		DefaultModels    []string               `json:"default_models"`
		ReasoningEnabled bool                   `json:"reasoning_enabled"`
		ReasoningEffort  string                 `json:"reasoning_effort"`
		EditingEnabled   bool                   `json:"editing_enabled"`
		EditData         *models.Artwork        `json:"edit_data,omitempty"`
		EditArtworks     []SafeArtwork          `json:"edit_artworks,omitempty"`
	}{
		Models:           templateData.Models,
		Examples:         templateData.Examples,
		DefaultTemp:      templateData.DefaultTemp,
		DefaultMaxTokens: templateData.DefaultMaxTokens,
		DefaultModels:    templateData.DefaultModels,
		ReasoningEnabled: templateData.ReasoningEnabled,
		ReasoningEffort:  templateData.ReasoningEffort,
		EditingEnabled:   config.IsEditingEnabled(),
		EditData:         editData,
		EditArtworks:     safeEditArtworks,
	}

	w.Header().Set("Content-Type", "text/html")
	if err := h.tmpl.ExecuteTemplate(w, "workshop.html", currentTemplateData); err != nil {
		log.Printf("Failed to execute template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
