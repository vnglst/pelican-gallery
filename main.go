package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"pelican-gallery/internal/api"
	"pelican-gallery/internal/config"
	"pelican-gallery/internal/database"
	"pelican-gallery/internal/models"
	"pelican-gallery/internal/pages"

	"github.com/joho/godotenv"
)

// RateLimiter implements a simple in-memory rate limiter
// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	mu       sync.RWMutex
	requests map[string][]time.Time
	window   time.Duration
	limit    int
}

func NewRateLimiter(window time.Duration, limit int) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		window:   window,
		limit:    limit,
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	if requests, exists := rl.requests[key]; exists {
		validRequests := make([]time.Time, 0, len(requests))
		for _, req := range requests {
			if req.After(windowStart) {
				validRequests = append(validRequests, req)
			}
		}
		rl.requests[key] = validRequests
	}

	if len(rl.requests[key]) < rl.limit {
		rl.requests[key] = append(rl.requests[key], now)
		return true
	}

	return false
}

func (rl *RateLimiter) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientIP := getClientIP(r)
		if !rl.Allow(clientIP) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}

func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in case of multiple
		if idx := strings.Index(xff, ","); idx > 0 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}

	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}

	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}

	return r.RemoteAddr
}

//go:embed static/*
var staticFiles embed.FS

//go:embed templates/*
var templateFiles embed.FS

// isDevelopment checks if we're running in development mode
// isDevelopment checks if we're running in development mode
func isDevelopment() bool {
	return os.Getenv("GO_ENV") != "production"
}

// getStaticFS returns the appropriate file system for static files
// getStaticFS returns the appropriate file system for static files
func getStaticFS() http.FileSystem {
	if isDevelopment() {
		return http.Dir("static")
	}
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("Failed to create static file system: %v", err)
	}
	return http.FS(staticFS)
}

// parseTemplates returns the appropriate template for the environment
// parseTemplates returns the appropriate template for the environment
func parseTemplates() (*template.Template, error) {
	// Create template with custom functions
	funcMap := template.FuncMap{
		"modelName": getModelDisplayName,
		"json": func(v interface{}) (string, error) {
			b, err := json.Marshal(v)
			if err != nil {
				return "", err
			}
			return string(b), nil
		},
	}

	if isDevelopment() {
		tmpl := template.New("").Funcs(funcMap)
		return tmpl.ParseGlob("templates/*.html")
	}
	tmpl := template.New("").Funcs(funcMap)
	return tmpl.ParseFS(templateFiles, "templates/*.html")
}

// getModelDisplayName returns the display name for a model ID
// getModelDisplayName returns the display name for a model ID
func getModelDisplayName(modelID string) string {
	allModels := config.GetAvailableModels()
	for _, model := range allModels {
		if model.ID == modelID {
			return model.Name
		}
	}
	// Return the ID if no match found
	return modelID
}

// loggingMiddleware logs all HTTP requests
// loggingMiddleware logs all HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log the request
		log.Printf("Started %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

		// Create a response writer wrapper to capture status code
		wrapper := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(wrapper, r)

		// Log the response
		duration := time.Since(start)
		log.Printf("Completed %s %s with status %d in %v", r.Method, r.URL.Path, wrapper.statusCode, duration)
	})
}

// responseWriter wrapper to capture status code
// responseWriter wrapper to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	if apiKey := os.Getenv("OPENROUTER_API_KEY"); apiKey == "" {
		log.Println("WARNING: OPENROUTER_API_KEY environment variable not found - artwork generation will be disabled")
	} else {
		log.Println("INFO: OPENROUTER_API_KEY found - artwork generation is enabled")
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "artworks.db"
	}

	db, err := database.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	log.Printf("Database initialized at: %s", dbPath)

	promptConfig, err := config.LoadPromptConfig("config/prompt.yaml")
	if err != nil {
		log.Fatalf("Failed to load prompt config: %v", err)
	}

	tmpl, err := parseTemplates()
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	templateData := models.TemplateData{
		Models:           config.GetAvailableModels(),
		Examples:         config.GetExamplePrompts(),
		DefaultTemp:      config.GetDefaultTemperature(),
		DefaultMaxTokens: config.GetDefaultMaxTokens(),
		DefaultModels:    config.GetDefaultModels(),
		ReasoningEnabled: config.GetDefaultReasoningEnabled(),
		ReasoningEffort:  config.GetDefaultReasoningEffort(),
		EditingEnabled:   config.IsEditingEnabled(),
	}

	apiHandler := api.NewHandler(promptConfig, db, tmpl)

	pageHandler := pages.NewPageHandler(db, tmpl, templateData)

	rateLimiter := NewRateLimiter(time.Minute, 100)

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(getStaticFS())))

	mux.HandleFunc("/", pageHandler.HomepageHandler)
	mux.HandleFunc("/workshop", pageHandler.WorkshopHandler)
	mux.HandleFunc("/gallery/", func(w http.ResponseWriter, r *http.Request) {
		// Extract category from path: /gallery/category/nature -> "nature"
		path := r.URL.Path
		category := ""

		if path != "/gallery" && path != "/gallery/" {
			// Check if it's a category path
			if strings.HasPrefix(path, "/gallery/category/") {
				category = strings.TrimPrefix(path, "/gallery/category/")
				// URL decode the category
				if decoded, err := url.QueryUnescape(category); err == nil {
					category = decoded
				}
			} else {
				// Invalid path
				http.NotFound(w, r)
				return
			}
		}

		r.URL.RawQuery = r.URL.RawQuery + "&category=" + url.QueryEscape(category)
		pageHandler.GalleryHandler(w, r)
	})

	mux.HandleFunc("/api/generate", rateLimiter.Middleware(apiHandler.GenerateHandler))
	mux.HandleFunc("/api/save-artwork", rateLimiter.Middleware(apiHandler.SaveArtworkHandler))
	mux.HandleFunc("/api/regenerate-artwork", rateLimiter.Middleware(apiHandler.RegenerateArtworkHandler))
	mux.HandleFunc("/api/delete-artwork", rateLimiter.Middleware(apiHandler.DeleteArtworkHandler))

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Pelican Art Gallery starting on http://localhost:%s\n", port)
	fmt.Println("Press Ctrl+C to stop the server")

	loggedMux := loggingMiddleware(mux)

	if err := http.ListenAndServe(":"+port, loggedMux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
