package database

import (
	"database/sql"
	"fmt"

	"pelican-gallery/internal/models"

	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
}

// New creates a new database connection and initializes the schema
func New(dbPath string) (*DB, error) {
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &DB{conn: conn}

	if err := db.CreateTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// CreateTables creates the necessary tables if they don't exist
func (db *DB) CreateTables() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS artworks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		slug TEXT NOT NULL,
		category TEXT NOT NULL,
		prompt TEXT NOT NULL,
		model TEXT NOT NULL,
		svg_content TEXT NOT NULL,
		temperature REAL NOT NULL,
		max_tokens INTEGER NOT NULL,
		created_at DATETIME NOT NULL,
		UNIQUE(slug, model) -- Ensure only one artwork per slug+model combination
	);
	
	CREATE INDEX IF NOT EXISTS idx_artworks_slug ON artworks(slug);
	CREATE INDEX IF NOT EXISTS idx_artworks_category ON artworks(category);
	CREATE INDEX IF NOT EXISTS idx_artworks_created_at ON artworks(created_at);
	`

	_, err := db.conn.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

// SaveArtwork saves an artwork to the database with upsert behavior (one per model+slug combination)
func (db *DB) SaveArtwork(artwork models.Artwork) error {
	// Use INSERT OR REPLACE to ensure only one artwork per model+slug combination
	// This leverages the UNIQUE(slug, model) constraint in the table
	query := `
	INSERT OR REPLACE INTO artworks (title, slug, category, prompt, model, svg_content, temperature, max_tokens, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.conn.Exec(query,
		artwork.Title,
		artwork.Slug,
		artwork.Category,
		artwork.Prompt,
		artwork.Model,
		artwork.SVGContent,
		artwork.Temperature,
		artwork.MaxTokens,
		artwork.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save artwork: %w", err)
	}

	return nil
}

// GetAllArtworks retrieves all artworks from the database, ordered by creation date (newest first)
func (db *DB) GetAllArtworks() ([]models.Artwork, error) {
	query := `
	SELECT id, title, slug, category, prompt, model, svg_content, temperature, max_tokens, created_at
	FROM artworks
	ORDER BY created_at DESC
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query artworks: %w", err)
	}
	defer rows.Close()

	var artworks []models.Artwork
	for rows.Next() {
		var artwork models.Artwork
		err := rows.Scan(
			&artwork.ID,
			&artwork.Title,
			&artwork.Slug,
			&artwork.Category,
			&artwork.Prompt,
			&artwork.Model,
			&artwork.SVGContent,
			&artwork.Temperature,
			&artwork.MaxTokens,
			&artwork.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan artwork: %w", err)
		}
		artworks = append(artworks, artwork)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return artworks, nil
}

// GetArtworkBySlug retrieves a specific artwork by its slug
func (db *DB) GetArtworkBySlug(slug string) (*models.Artwork, error) {
	query := `
	SELECT id, slug, category, prompt, model, svg_content, temperature, max_tokens, created_at
	FROM artworks
	WHERE slug = ?
	LIMIT 1
	`

	var artwork models.Artwork
	err := db.conn.QueryRow(query, slug).Scan(
		&artwork.ID,
		&artwork.Slug,
		&artwork.Category,
		&artwork.Prompt,
		&artwork.Model,
		&artwork.SVGContent,
		&artwork.Temperature,
		&artwork.MaxTokens,
		&artwork.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("artwork not found")
		}
		return nil, fmt.Errorf("failed to get artwork: %w", err)
	}

	return &artwork, nil
}

// GetArtworksByCategory retrieves all artworks in a specific category
func (db *DB) GetArtworksByCategory(category string) ([]models.Artwork, error) {
	query := `
	SELECT id, title, slug, category, prompt, model, svg_content, temperature, max_tokens, created_at
	FROM artworks
	WHERE category = ?
	ORDER BY created_at DESC
	`

	rows, err := db.conn.Query(query, category)
	if err != nil {
		return nil, fmt.Errorf("failed to query artworks: %w", err)
	}
	defer rows.Close()

	var artworks []models.Artwork
	for rows.Next() {
		var artwork models.Artwork
		err := rows.Scan(
			&artwork.ID,
			&artwork.Title,
			&artwork.Slug,
			&artwork.Category,
			&artwork.Prompt,
			&artwork.Model,
			&artwork.SVGContent,
			&artwork.Temperature,
			&artwork.MaxTokens,
			&artwork.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan artwork row: %w", err)
		}
		artworks = append(artworks, artwork)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return artworks, nil
}

// GetArtworksBySlug retrieves all artworks with a specific slug
func (db *DB) GetArtworksBySlug(slug string) ([]models.Artwork, error) {
	query := `
	SELECT id, title, slug, category, prompt, model, svg_content, temperature, max_tokens, created_at
	FROM artworks 
	WHERE slug = ?
	ORDER BY created_at DESC`

	rows, err := db.conn.Query(query, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to query artworks by slug: %w", err)
	}
	defer rows.Close()

	var artworks []models.Artwork
	for rows.Next() {
		var artwork models.Artwork
		err := rows.Scan(
			&artwork.ID,
			&artwork.Title,
			&artwork.Slug,
			&artwork.Category,
			&artwork.Prompt,
			&artwork.Model,
			&artwork.SVGContent,
			&artwork.Temperature,
			&artwork.MaxTokens,
			&artwork.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan artwork row: %w", err)
		}
		artworks = append(artworks, artwork)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating artwork rows: %w", err)
	}

	return artworks, nil
}

// GetAllCategories retrieves all unique categories from the database
func (db *DB) GetAllCategories() ([]string, error) {
	query := `
	SELECT DISTINCT category
	FROM artworks
	ORDER BY category ASC
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		err := rows.Scan(&category)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return categories, nil
}

// DeleteArtwork deletes an artwork by ID
func (db *DB) DeleteArtwork(id int) error {
	query := `DELETE FROM artworks WHERE id = ?`

	result, err := db.conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete artwork: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("artwork with ID %d not found", id)
	}

	return nil
}
