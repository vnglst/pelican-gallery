package pages

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"pelican-gallery/internal/config"
	"pelican-gallery/internal/database"
	"pelican-gallery/internal/models"
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

	// If no category specified, redirect to first available category
	if category == "" {
		// Get all available categories for navigation
		categories, err := h.db.GetDistinctCategories()
		if err != nil {
			log.Printf("Error fetching categories: %v", err)
			http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
			return
		}

		// If we have categories, redirect to the first one
		if len(categories) > 0 {
			http.Redirect(w, r, "/gallery/category/"+categories[0], http.StatusFound)
			return
		}
		// If no categories, continue with empty category (will show empty gallery)
	}

	// Get groups with their artworks using the database helper
	groups, artworkMap, err := h.db.ListGroupsWithArtworks(category)
	if err != nil {
		log.Printf("Error fetching groups with artworks: %v", err)
		http.Error(w, "Failed to fetch artworks", http.StatusInternalServerError)
		return
	}

	// Get all available categories for navigation
	categories, err := h.db.GetDistinctCategories()
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}

	// Build gallery data structures
	type GalleryArtwork struct {
		models.Artwork
		Title      string        `json:"title"`       // From group
		Category   string        `json:"category"`    // From group
		Prompt     string        `json:"prompt"`      // From group
		SVGContent template.HTML `json:"svg_content"` // HTML-safe SVG
	}

	type GalleryGroup struct {
		models.ArtworkGroup
		Artworks []GalleryArtwork `json:"artworks"`
	}

	var galleryGroups []GalleryGroup
	for _, group := range groups {
		artworks := artworkMap[group.ID] // Get artworks for this group

		galleryArtworks := make([]GalleryArtwork, len(artworks))
		for i, artwork := range artworks {
			galleryArtworks[i] = GalleryArtwork{
				Artwork:    artwork,
				Title:      group.Title,
				Category:   group.Category,
				Prompt:     group.Prompt,
				SVGContent: template.HTML(artwork.SVG),
			}
		}

		galleryGroups = append(galleryGroups, GalleryGroup{
			ArtworkGroup: group,
			Artworks:     galleryArtworks,
		})
	}

	log.Printf("Fetched %d groups with artworks and %d categories for gallery", len(galleryGroups), len(categories))

	data := struct {
		Title          string         `json:"title"`
		Groups         []GalleryGroup `json:"groups"`
		Categories     []string       `json:"categories"`
		Category       string         `json:"category"`
		EditingEnabled bool           `json:"editing_enabled"`
	}{
		Title:          "Gallery - Pelican Art Gallery",
		Groups:         galleryGroups,
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

	w.Header().Set("Content-Type", "text/html")
	homepageData := struct {
		EditingEnabled bool `json:"editing_enabled"`
	}{
		EditingEnabled: config.IsEditingEnabled(),
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
