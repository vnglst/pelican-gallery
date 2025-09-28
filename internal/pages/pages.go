package pages

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"pelican-gallery/internal/config"
	"pelican-gallery/internal/database"
	"pelican-gallery/internal/models"
)

// Filter constants for model providers
const (
	FilterOpenAI    = "openai"
	FilterAnthropic = "anthropic"
	FilterGoogle    = "google"
	FilterOther     = "other"
)

// TemplateParser is a function type for parsing templates
type TemplateParser func(*template.Template) (*template.Template, error)

// PageHandler contains the page handlers
type PageHandler struct {
	db             *database.DB
	tmpl           *template.Template
	templateData   models.TemplateData
	templateParser TemplateParser
}

// NewPageHandler creates a new page handler
func NewPageHandler(db *database.DB, tmpl *template.Template, templateData models.TemplateData, templateParser TemplateParser) *PageHandler {
	return &PageHandler{
		db:             db,
		tmpl:           tmpl,
		templateData:   templateData,
		templateParser: templateParser,
	}
}

// getTemplate returns the appropriate template (cached or re-parsed)
func (h *PageHandler) getTemplate() (*template.Template, error) {
	if h.templateParser != nil {
		return h.templateParser(h.tmpl)
	}
	return h.tmpl, nil
}

// GalleryHandler handles requests to display the gallery of saved artworks
func (h *PageHandler) GalleryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	category := r.URL.Query().Get("category")

	// No model filtering on gallery page — show all artworks for the selected category

	// If no category specified, redirect to first available category
	if category == "" {
		categories, err := h.db.GetDistinctCategories()
		if err != nil {
			log.Printf("Error fetching categories: %v", err)
			http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
			return
		}
		if len(categories) > 0 {
			http.Redirect(w, r, "/gallery/category/"+categories[0], http.StatusFound)
			return
		}
	}

	groups, artworkMap, err := h.db.ListGroupsWithArtworks(category)
	if err != nil {
		log.Printf("Error fetching groups with artworks: %v", err)
		http.Error(w, "Failed to fetch artworks", http.StatusInternalServerError)
		return
	}

	categories, err := h.db.GetDistinctCategories()
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}

	// Only include artworks from these three models (case-insensitive substring match)
	allowedModelSubs := []string{
		"anthropic/claude-sonnet-4",
		"google/gemini-2.5-pro",
		"openai/gpt-5",
	}
	allowedModelsContains := func(model string) bool {
		if model == "" {
			return false
		}
		low := strings.ToLower(model)
		for _, sub := range allowedModelSubs {
			if low == strings.ToLower(sub) {
				return true
			}
		}
		return false
	}

	type GalleryArtwork struct {
		models.Artwork
		Title      string        `json:"title"`
		Category   string        `json:"category"`
		Prompt     string        `json:"prompt"`
		SVGContent template.HTML `json:"svg_content"`
	}

	type GalleryGroup struct {
		models.ArtworkGroup
		Artworks []GalleryArtwork `json:"artworks"`
	}

	var galleryGroups []GalleryGroup
	var flatArtworks []GalleryArtwork
	for _, group := range groups {
		artworks := artworkMap[group.ID]
		var filteredArtworks []GalleryArtwork
		for _, artwork := range artworks {
			if allowedModelsContains(artwork.Model) {
				ga := GalleryArtwork{
					Artwork:    artwork,
					Title:      group.Title,
					Category:   group.Category,
					Prompt:     group.Prompt,
					SVGContent: template.HTML(artwork.SVG),
				}
				filteredArtworks = append(filteredArtworks, ga)
				// append to flat list as well
				flatArtworks = append(flatArtworks, ga)
			}
		}
		galleryGroups = append(galleryGroups, GalleryGroup{
			ArtworkGroup: group,
			Artworks:     filteredArtworks,
		})
	}

	log.Printf("Fetched %d groups with artworks and %d categories for gallery", len(galleryGroups), len(categories))

	data := struct {
		Title          string           `json:"title"`
		Groups         []GalleryGroup   `json:"groups"`
		Artworks       []GalleryArtwork `json:"artworks"`
		Categories     []string         `json:"categories"`
		Category       string           `json:"category"`
		EditingEnabled bool             `json:"editing_enabled"`
	}{
		Title:          "Gallery - Pelican Art Gallery",
		Groups:         galleryGroups,
		Artworks:       flatArtworks,
		Categories:     categories,
		Category:       category,
		EditingEnabled: isEditingEnabled(),
	}

	w.Header().Set("Content-Type", "text/html")

	tmpl, err := h.getTemplate()
	if err != nil {
		log.Printf("Error getting template: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "gallery.html", data); err != nil {
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

	// Get a random group with artworks from anthropic/claude-sonnet-4 and openai/gpt-5
	randomGroup, randomArtworks, err := h.db.GetRandomGroupWithModelArtworks("anthropic/claude-sonnet-4", "openai/gpt-5")
	var featuredGroup *models.ArtworkGroup
	var featuredArtworks []models.Artwork

	if err != nil {
		log.Printf("No random group found with both models, trying fallback: %v", err)
		// Fallback: try to get any random group with artworks from either model
		randomGroup, randomArtworks, err = h.db.GetRandomGroupWithModelArtworks("anthropic", "openai")
		if err != nil {
			log.Printf("No fallback group found either: %v", err)
			// If still no group found, just continue without featured content
		} else {
			featuredGroup = randomGroup
			featuredArtworks = randomArtworks
		}
	} else {
		featuredGroup = randomGroup
		featuredArtworks = randomArtworks
	}

	type HomepageArtwork struct {
		models.Artwork
		SVGContent template.HTML `json:"svg_content"`
	}

	var homepageArtworks []HomepageArtwork
	for _, artwork := range featuredArtworks {
		homepageArtworks = append(homepageArtworks, HomepageArtwork{
			Artwork:    artwork,
			SVGContent: template.HTML(artwork.SVG),
		})
	}

	w.Header().Set("Content-Type", "text/html")
	homepageData := struct {
		EditingEnabled   bool                 `json:"editing_enabled"`
		FeaturedGroup    *models.ArtworkGroup `json:"featured_group,omitempty"`
		FeaturedArtworks []HomepageArtwork    `json:"featured_artworks,omitempty"`
	}{
		EditingEnabled:   config.IsEditingEnabled(),
		FeaturedGroup:    featuredGroup,
		FeaturedArtworks: homepageArtworks,
	}

	tmpl, err := h.getTemplate()
	if err != nil {
		log.Printf("Error getting template: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "homepage.html", homepageData); err != nil {
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

	// Check if we're editing an existing artwork group
	editIDStr := r.URL.Query().Get("edit")
	var editGroup *models.ArtworkGroup
	var editArtworks []models.Artwork

	if editIDStr != "" {
		// Parse group ID
		var editID int
		if _, err := fmt.Sscanf(editIDStr, "%d", &editID); err != nil {
			log.Printf("Invalid edit ID: %s", editIDStr)
		} else {
			group, err := h.db.GetGroup(editID)
			if err != nil {
				log.Printf("Error fetching group for editing: %v", err)
			} else {
				editGroup = group
				editArtworks, err = h.db.ListArtworksByGroup(editID)
				if err != nil {
					log.Printf("Error fetching artworks for group %d: %v", editID, err)
				}
				log.Printf("Found group %d with %d artwork(s) for editing: %s", editID, len(editArtworks), group.Title)
			}
		}
	}

	// Prepare template data
	templateData := h.templateData

	// Create template data with edit information
	currentTemplateData := struct {
		Models       []models.ModelInfo   `json:"models"`
		EditGroup    *models.ArtworkGroup `json:"edit_group,omitempty"`
		EditArtworks []models.Artwork     `json:"edit_artworks,omitempty"`
	}{
		Models:       templateData.Models,
		EditGroup:    editGroup,
		EditArtworks: editArtworks,
	}

	w.Header().Set("Content-Type", "text/html")

	tmpl, err := h.getTemplate()
	if err != nil {
		log.Printf("Error getting template: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "workshop.html", currentTemplateData); err != nil {
		log.Printf("Failed to execute template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// ArtworkGroupHandler shows a page dedicated to a group and all its artworks
func (h *PageHandler) ArtworkGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Expect path like /group/123 or /group/123/
	raw := strings.TrimPrefix(r.URL.Path, "/group/")
	raw = strings.TrimSuffix(raw, "/")
	if raw == "" {
		log.Printf("ArtworkGroupHandler: empty group id in path: %q", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(raw)
	if err != nil {
		log.Printf("ArtworkGroupHandler: failed to parse group id from path %q: %v", r.URL.Path, err)
		http.NotFound(w, r)
		return
	}

	group, err := h.db.GetGroup(id)
	if err != nil {
		log.Printf("ArtworkGroupHandler: GetGroup(%d) error: %v", id, err)
		http.NotFound(w, r)
		return
	}

	// Parse model filters from query parameters (can be multiple)
	modelFilters := r.URL.Query()["model"]

	artworks, err := h.db.ListArtworksByGroup(id)
	if err != nil {
		log.Printf("Error fetching artworks for group %d: %v", id, err)
		http.Error(w, "Failed to load artworks", http.StatusInternalServerError)
		return
	}

	// If model filters are present, filter the artworks accordingly
	// Supported filters: "openai", "anthropic", "google", "other"
	var filtered []models.Artwork
	if len(modelFilters) == 0 {
		filtered = artworks
	} else {
		for _, a := range artworks {
			show := false
			lowModel := strings.ToLower(a.Model)
			for _, f := range modelFilters {
				ff := strings.ToLower(f)
				if ff == "other" {
					if !(strings.Contains(lowModel, "openai") || strings.Contains(lowModel, "anthropic") || strings.Contains(lowModel, "google")) {
						show = true
						break
					}
				} else {
					if strings.Contains(lowModel, ff) {
						show = true
						break
					}
				}
			}
			if show {
				filtered = append(filtered, a)
			}
		}
	}

	// Build template data using the filtered list
	type ArtworkWithHTML struct {
		models.Artwork
		SVGContent template.HTML
	}

	var artList []ArtworkWithHTML
	for _, a := range filtered {
		artList = append(artList, ArtworkWithHTML{Artwork: a, SVGContent: template.HTML(a.SVG)})
	}

	data := struct {
		Title          string
		Group          *models.ArtworkGroup
		Artworks       []ArtworkWithHTML
		EditingEnabled bool
		ModelFilters   []string
	}{
		Title:          "Artwork Group - Pelican Art Gallery",
		Group:          group,
		Artworks:       artList,
		EditingEnabled: isEditingEnabled(),
		ModelFilters:   modelFilters,
	}

	tmpl, err := h.getTemplate()
	if err != nil {
		log.Printf("Error getting template: %v", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.ExecuteTemplate(w, "artwork-group.html", data); err != nil {
		log.Printf("Failed to execute artwork-group template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
