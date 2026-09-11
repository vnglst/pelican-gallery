package pages

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"pelican-gallery/internal/config"
	"pelican-gallery/internal/database"
	"pelican-gallery/internal/models"
)

// Filter constants for model providers
const (
	FilterOpenAI     = "openai"
	FilterAnthropic  = "anthropic"
	FilterGoogle     = "google"
	FilterOpenSource = "open-source"
)

var modelReleaseDates = map[string]time.Time{
	"anthropic/claude-3-haiku":    time.Date(2024, time.March, 13, 0, 0, 0, 0, time.UTC),
	"anthropic/claude-3.5-haiku":  time.Date(2024, time.October, 22, 0, 0, 0, 0, time.UTC),
	"anthropic/claude-3.7-sonnet": time.Date(2025, time.February, 24, 0, 0, 0, 0, time.UTC),
	"anthropic/claude-sonnet-4":   time.Date(2025, time.May, 22, 0, 0, 0, 0, time.UTC),
	"anthropic/claude-opus-4.1":   time.Date(2025, time.August, 5, 0, 0, 0, 0, time.UTC),
	"anthropic/claude-sonnet-4.5": time.Date(2025, time.September, 29, 0, 0, 0, 0, time.UTC),
	"anthropic/claude-opus-4.5":   time.Date(2025, time.November, 24, 0, 0, 0, 0, time.UTC),
	"anthropic/claude-opus-4.6":   time.Date(2026, time.February, 5, 0, 0, 0, 0, time.UTC),
	"openai/gpt-3.5-turbo":        time.Date(2023, time.March, 1, 0, 0, 0, 0, time.UTC),
	"openai/gpt-4":                time.Date(2023, time.March, 14, 0, 0, 0, 0, time.UTC),
	"openai/gpt-3.5-turbo-0613":   time.Date(2023, time.June, 13, 0, 0, 0, 0, time.UTC),
	"openai/gpt-4o":               time.Date(2024, time.May, 13, 0, 0, 0, 0, time.UTC),
	"openai/gpt-4o-mini":          time.Date(2024, time.July, 18, 0, 0, 0, 0, time.UTC),
	"openai/o1":                   time.Date(2024, time.December, 5, 0, 0, 0, 0, time.UTC),
	"openai/gpt-4.1":              time.Date(2025, time.April, 14, 0, 0, 0, 0, time.UTC),
	"openai/o3":                   time.Date(2025, time.April, 16, 0, 0, 0, 0, time.UTC),
	"openai/o4-mini":              time.Date(2025, time.April, 16, 0, 0, 0, 0, time.UTC),
	"openai/o4-mini-high":         time.Date(2025, time.April, 16, 0, 0, 0, 0, time.UTC),
	"openai/gpt-oss-120b":         time.Date(2025, time.August, 5, 0, 0, 0, 0, time.UTC),
	"openai/gpt-oss-20b":          time.Date(2025, time.August, 5, 0, 0, 0, 0, time.UTC),
	"openai/gpt-oss-20b:free":     time.Date(2025, time.August, 5, 0, 0, 0, 0, time.UTC),
	"openai/gpt-5":                time.Date(2025, time.August, 7, 0, 0, 0, 0, time.UTC),
	"openai/gpt-5-chat":           time.Date(2025, time.August, 7, 0, 0, 0, 0, time.UTC),
	"openai/gpt-5-mini":           time.Date(2025, time.August, 7, 0, 0, 0, 0, time.UTC),
	"openai/gpt-5-nano":           time.Date(2025, time.August, 7, 0, 0, 0, 0, time.UTC),
	"openai/gpt-5-codex":          time.Date(2025, time.September, 15, 0, 0, 0, 0, time.UTC),
	"openai/gpt-5.1":              time.Date(2025, time.November, 13, 0, 0, 0, 0, time.UTC),
	"openai/gpt-5.2":              time.Date(2025, time.December, 11, 0, 0, 0, 0, time.UTC),
	"openai/gpt-6-astra-pro":      time.Date(2026, time.September, 3, 0, 0, 0, 0, time.UTC),
}

var modelVersionPattern = regexp.MustCompile(`\d+`)

func modelCapabilityRank(model string) int {
	model = strings.ToLower(model)
	switch {
	case strings.Contains(model, "nano"), strings.Contains(model, "haiku"):
		return 0
	case strings.Contains(model, "mini"), strings.Contains(model, "lite"), strings.Contains(model, "flash"):
		return 1
	case strings.Contains(model, "opus"):
		return 4
	case strings.Contains(model, "pro"):
		return 3
	case strings.Contains(model, "sonnet"):
		return 2
	default:
		return 2
	}
}

func modelProvider(model, storedMetadata string) string {
	if config.IsOpenSourceModel(model, storedMetadata) {
		return FilterOpenSource
	}

	prefix, _, _ := strings.Cut(strings.ToLower(model), "/")
	switch prefix {
	case FilterOpenAI, FilterAnthropic, FilterGoogle:
		return prefix
	}
	return ""
}

func modelReleaseDate(model string, storedOpenRouterCreated int64) (time.Time, bool) {
	model = strings.ToLower(model)
	if releasedAt, ok := modelReleaseDates[model]; ok {
		return releasedAt, true
	}

	lookupID := strings.TrimSuffix(model, ":free")
	if releasedAt, ok := modelReleaseDates[lookupID]; ok {
		return releasedAt, true
	}
	if storedOpenRouterCreated > 0 {
		return time.Unix(storedOpenRouterCreated, 0).UTC(), true
	}
	if info, ok := config.GetModelInfo(lookupID); ok && info.Created > 0 {
		return time.Unix(info.Created, 0).UTC(), true
	}
	return time.Time{}, false
}

func modelVersionLess(left, right string) bool {
	leftParts := modelVersionPattern.FindAllString(strings.ToLower(left), -1)
	rightParts := modelVersionPattern.FindAllString(strings.ToLower(right), -1)
	for i := 0; i < len(leftParts) && i < len(rightParts); i++ {
		leftNumber, _ := strconv.Atoi(leftParts[i])
		rightNumber, _ := strconv.Atoi(rightParts[i])
		if leftNumber != rightNumber {
			return leftNumber < rightNumber
		}
	}
	if len(leftParts) != len(rightParts) {
		return len(leftParts) < len(rightParts)
	}
	return strings.ToLower(left) < strings.ToLower(right)
}

func chronologyModelName(model string) string {
	lookupID := strings.TrimSuffix(strings.ToLower(model), ":free")
	if info, ok := config.GetModelInfo(lookupID); ok && info.Name != "" {
		return chronologyDisplayName(info.Name)
	}
	_, name, found := strings.Cut(model, "/")
	if found {
		return modelIDDisplayName(name)
	}
	return modelIDDisplayName(model)
}

func modelIDDisplayName(modelID string) string {
	words := strings.FieldsFunc(strings.TrimSuffix(modelID, ":free"), func(r rune) bool {
		return r == '-' || r == '_'
	})
	for i, word := range words {
		lower := strings.ToLower(word)
		switch lower {
		case "gpt":
			words[i] = "GPT"
		case "ai":
			words[i] = "AI"
		default:
			if len(word) > 0 {
				words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
			}
		}
	}
	return strings.Join(words, " ")
}

func chronologyDisplayName(name string) string {
	if _, displayName, found := strings.Cut(name, ": "); found {
		return displayName
	}
	return name
}

func formatGenerationCost(cost float64) string {
	if cost <= 0 {
		return "Free"
	}
	if cost >= 0.01 {
		return fmt.Sprintf("$%.2f", cost)
	}
	if cost >= 0.0001 {
		return fmt.Sprintf("$%.4f", cost)
	}
	if cost >= 0.00001 {
		return fmt.Sprintf("$%.5f", cost)
	}
	return "<$0.00001"
}

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

// getCSSHash computes and returns the MD5 hash of the output.css file for cache busting
func (h *PageHandler) getCSSHash() string {
	cssPath := "static/css/output.css"
	content, err := os.ReadFile(cssPath)
	if err != nil {
		log.Printf("Error reading CSS file for hash: %v", err)
		return ""
	}
	hash := md5.Sum(content)
	return fmt.Sprintf("%x", hash)
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

	type GalleryArtwork struct {
		models.Artwork
		SVGContent     template.HTML `json:"svg_content"`
		DisplayName    string
		CostDisplay    string
		ReleasedAt     time.Time
		HasReleaseDate bool
	}

	type GalleryProvider struct {
		ID       string
		Label    string
		Artworks []GalleryArtwork
	}

	type GalleryGroup struct {
		models.ArtworkGroup
		Artworks           []GalleryArtwork `json:"artworks"`
		Providers          []GalleryProvider
		FeaturedSVG        template.HTML
		FeaturedModel      string
		ArtworkCount       int
		HasOriginalArtwork bool `json:"has_original_artwork"`
	}

	var galleryGroups []GalleryGroup
	var flatArtworks []GalleryArtwork
	for _, group := range groups {
		artworks := artworkMap[group.ID]
		var galleryArtworks []GalleryArtwork
		var featuredArtwork *models.Artwork
		var gpt5Artwork *models.Artwork

		for i, artwork := range artworks {
			if artwork.Featured {
				featuredArtwork = &artworks[i]
				break
			}
			if strings.ToLower(artwork.Model) == "openai/gpt-5" {
				gpt5Artwork = &artworks[i]
			}
		}

		// Use featured if available, otherwise fallback to GPT-5
		selectedArtwork := featuredArtwork
		if selectedArtwork == nil {
			selectedArtwork = gpt5Artwork
		}

		for _, artwork := range artworks {
			releasedAt, hasReleaseDate := modelReleaseDate(artwork.Model, artwork.ModelCreatedAt)
			displayName := chronologyDisplayName(artwork.ModelName)
			if displayName == "" {
				displayName = chronologyModelName(artwork.Model)
			}
			galleryArtworks = append(galleryArtworks, GalleryArtwork{
				Artwork: artwork, SVGContent: template.HTML(artwork.SVG), DisplayName: displayName,
				CostDisplay: formatGenerationCost(artwork.GenerationCostUSD), ReleasedAt: releasedAt, HasReleaseDate: hasReleaseDate,
			})
		}
		sort.SliceStable(galleryArtworks, func(i, j int) bool {
			if galleryArtworks[i].HasReleaseDate && galleryArtworks[j].HasReleaseDate {
				if !galleryArtworks[i].ReleasedAt.Equal(galleryArtworks[j].ReleasedAt) {
					return galleryArtworks[i].ReleasedAt.Before(galleryArtworks[j].ReleasedAt)
				}
				if leftRank, rightRank := modelCapabilityRank(galleryArtworks[i].Model), modelCapabilityRank(galleryArtworks[j].Model); leftRank != rightRank {
					return leftRank < rightRank
				}
			}
			if galleryArtworks[i].HasReleaseDate != galleryArtworks[j].HasReleaseDate {
				return galleryArtworks[i].HasReleaseDate
			}
			if galleryArtworks[i].Model != galleryArtworks[j].Model {
				return modelVersionLess(galleryArtworks[i].Model, galleryArtworks[j].Model)
			}
			return galleryArtworks[i].ID < galleryArtworks[j].ID
		})

		providerLabels := map[string]string{FilterOpenAI: "OpenAI", FilterGoogle: "Google", FilterAnthropic: "Anthropic", FilterOpenSource: "Open source"}
		var providers []GalleryProvider
		for _, providerID := range []string{FilterOpenAI, FilterGoogle, FilterAnthropic, FilterOpenSource} {
			provider := GalleryProvider{ID: providerID, Label: providerLabels[providerID]}
			for _, artwork := range galleryArtworks {
				if modelProvider(artwork.Model, artwork.ModelMetadata) == providerID {
					provider.Artworks = append(provider.Artworks, artwork)
				}
			}
			if len(provider.Artworks) > 0 {
				providers = append(providers, provider)
			}
		}

		var featuredSVG template.HTML
		var featuredModel string
		if selectedArtwork != nil {
			featuredSVG = template.HTML(selectedArtwork.SVG)
			featuredModel = chronologyDisplayName(selectedArtwork.ModelName)
			if featuredModel == "" {
				featuredModel = chronologyModelName(selectedArtwork.Model)
			}
		}

		hasOriginalArtwork := len(group.OriginalArtwork) > 0

		galleryGroups = append(galleryGroups, GalleryGroup{
			ArtworkGroup:       group,
			Artworks:           galleryArtworks,
			Providers:          providers,
			FeaturedSVG:        featuredSVG,
			FeaturedModel:      featuredModel,
			ArtworkCount:       len(galleryArtworks),
			HasOriginalArtwork: hasOriginalArtwork,
		})
		flatArtworks = append(flatArtworks, galleryArtworks...)
	}

	log.Printf("Fetched %d groups with artworks and %d categories for gallery", len(galleryGroups), len(categories))

	data := struct {
		Title          string           `json:"title"`
		Groups         []GalleryGroup   `json:"groups"`
		Artworks       []GalleryArtwork `json:"artworks"`
		Categories     []string         `json:"categories"`
		Category       string           `json:"category"`
		EditingEnabled bool             `json:"editing_enabled"`
		CSSHash        string           `json:"css_hash"`
	}{
		Title:          "Gallery - Pelican Art Gallery",
		Groups:         galleryGroups,
		Artworks:       flatArtworks,
		Categories:     categories,
		Category:       category,
		EditingEnabled: isEditingEnabled(),
		CSSHash:        h.getCSSHash(),
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

// AboutHandler handles requests to the about page
func (h *PageHandler) AboutHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := h.db.GetDistinctCategories()
	if err != nil {
		log.Printf("Error fetching categories for about page: %v", err)
	}

	w.Header().Set("Content-Type", "text/html")
	homepageData := struct {
		Categories []string `json:"categories"`
		CSSHash    string   `json:"css_hash"`
	}{
		Categories: categories,
		CSSHash:    h.getCSSHash(),
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
				} else {
					for i := range editArtworks {
						editArtworks[i].ModelProvider = modelProvider(editArtworks[i].Model, editArtworks[i].ModelMetadata)
						if releasedAt, ok := modelReleaseDate(editArtworks[i].Model, editArtworks[i].ModelCreatedAt); ok {
							editArtworks[i].ModelYear = releasedAt.Year()
							editArtworks[i].ModelSortTime = releasedAt.Unix()
						}
					}
				}
				log.Printf("Found group %d with %d artwork(s) for editing: %s", editID, len(editArtworks), group.Title)
			}
		}
	}

	// Prepare template data
	templateData := h.templateData

	// Create template data with edit information
	hasOriginalArtwork := false
	if editGroup != nil && editGroup.OriginalArtwork != nil && len(editGroup.OriginalArtwork) > 0 {
		hasOriginalArtwork = true
	}

	currentTemplateData := struct {
		Models             []models.ModelInfo   `json:"models"`
		EditGroup          *models.ArtworkGroup `json:"edit_group,omitempty"`
		EditArtworks       []models.Artwork     `json:"edit_artworks,omitempty"`
		HasOriginalArtwork bool                 `json:"has_original_artwork"`
		CSSHash            string               `json:"css_hash"`
	}{
		Models:             templateData.Models,
		EditGroup:          editGroup,
		EditArtworks:       editArtworks,
		HasOriginalArtwork: hasOriginalArtwork,
		CSSHash:            h.getCSSHash(),
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

	artworks, err := h.db.ListArtworksByGroup(id)
	if err != nil {
		log.Printf("Error fetching artworks for group %d: %v", id, err)
		http.Error(w, "Failed to load artworks", http.StatusInternalServerError)
		return
	}

	providerCounts := map[string]int{
		FilterOpenAI: 0, FilterAnthropic: 0, FilterGoogle: 0, FilterOpenSource: 0,
	}
	for _, artwork := range artworks {
		if category := modelProvider(artwork.Model, artwork.ModelMetadata); category != "" {
			providerCounts[category]++
		}
	}

	provider := strings.ToLower(r.URL.Query().Get("provider"))
	if providerCounts[provider] == 0 {
		provider = FilterOpenAI
		for _, candidate := range []string{FilterOpenAI, FilterGoogle, FilterAnthropic, FilterOpenSource} {
			if providerCounts[candidate] > 0 {
				provider = candidate
				break
			}
		}
	}

	type ArtworkWithHTML struct {
		models.Artwork
		SVGContent     template.HTML
		DisplayName    string
		CostDisplay    string
		ReleasedAt     time.Time
		HasReleaseDate bool
	}

	var artList []ArtworkWithHTML
	for _, artwork := range artworks {
		if modelProvider(artwork.Model, artwork.ModelMetadata) != provider {
			continue
		}
		releasedAt, hasReleaseDate := modelReleaseDate(artwork.Model, artwork.ModelCreatedAt)
		displayName := chronologyDisplayName(artwork.ModelName)
		if displayName == "" {
			displayName = chronologyModelName(artwork.Model)
		}
		artList = append(artList, ArtworkWithHTML{
			Artwork:        artwork,
			SVGContent:     template.HTML(artwork.SVG),
			DisplayName:    displayName,
			CostDisplay:    formatGenerationCost(artwork.GenerationCostUSD),
			ReleasedAt:     releasedAt,
			HasReleaseDate: hasReleaseDate,
		})
	}
	sort.SliceStable(artList, func(i, j int) bool {
		if artList[i].HasReleaseDate && artList[j].HasReleaseDate {
			if !artList[i].ReleasedAt.Equal(artList[j].ReleasedAt) {
				return artList[i].ReleasedAt.Before(artList[j].ReleasedAt)
			}
			if leftRank, rightRank := modelCapabilityRank(artList[i].Model), modelCapabilityRank(artList[j].Model); leftRank != rightRank {
				return leftRank < rightRank
			}
		}
		if artList[i].HasReleaseDate != artList[j].HasReleaseDate {
			return artList[i].HasReleaseDate
		}
		if artList[i].Model != artList[j].Model {
			return modelVersionLess(artList[i].Model, artList[j].Model)
		}
		return artList[i].ID < artList[j].ID
	})

	hasOriginalArtwork := len(group.OriginalArtwork) > 0

	data := struct {
		Title              string
		Group              *models.ArtworkGroup
		Artworks           []ArtworkWithHTML
		EditingEnabled     bool
		ActiveProvider     string
		ProviderCounts     map[string]int
		HasOriginalArtwork bool
		CSSHash            string
	}{
		Title:              "Artwork Group - Pelican Art Gallery",
		Group:              group,
		Artworks:           artList,
		EditingEnabled:     isEditingEnabled(),
		ActiveProvider:     provider,
		ProviderCounts:     providerCounts,
		HasOriginalArtwork: hasOriginalArtwork,
		CSSHash:            h.getCSSHash(),
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
