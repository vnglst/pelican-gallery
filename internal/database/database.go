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
	// Enable foreign key enforcement
	_, err := db.conn.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS artwork_groups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		prompt TEXT NOT NULL,
		category TEXT NOT NULL DEFAULT '',
        original_url TEXT NOT NULL DEFAULT '',
        artist_name TEXT NOT NULL DEFAULT '',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS artworks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		group_id INTEGER NOT NULL,
		model TEXT NOT NULL,
		temperature REAL NOT NULL DEFAULT 0.0,
		max_tokens INTEGER NOT NULL DEFAULT 0,
		svg TEXT DEFAULT '',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (group_id) REFERENCES artwork_groups(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_artworks_group_id ON artworks(group_id);
	CREATE INDEX IF NOT EXISTS idx_artwork_groups_created_at ON artwork_groups(created_at);
	CREATE INDEX IF NOT EXISTS idx_artworks_created_at ON artworks(created_at);
	`

	_, err = db.conn.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

// CreateGroup creates a new artwork group
func (db *DB) CreateGroup(group models.ArtworkGroup) (int, error) {
	query := `
		INSERT INTO artwork_groups (title, prompt, category, original_url, artist_name, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		`

	result, err := db.conn.Exec(query, group.Title, group.Prompt, group.Category, group.OriginalURL, group.ArtistName, group.CreatedAt, group.UpdatedAt)
	if err != nil {
		return 0, fmt.Errorf("failed to create group: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return int(id), nil
}

// UpdateGroup updates an existing artwork group
func (db *DB) UpdateGroup(group models.ArtworkGroup) error {
	query := `
		UPDATE artwork_groups
		SET title = ?, prompt = ?, category = ?, original_url = ?, artist_name = ?, updated_at = ?
		WHERE id = ?
		`

	result, err := db.conn.Exec(query, group.Title, group.Prompt, group.Category, group.OriginalURL, group.ArtistName, group.UpdatedAt, group.ID)
	if err != nil {
		return fmt.Errorf("failed to update group: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("group with ID %d not found", group.ID)
	}

	return nil
}

// GetGroup retrieves an artwork group by ID
func (db *DB) GetGroup(id int) (*models.ArtworkGroup, error) {
	query := `
	   SELECT id, title, prompt, category, original_url, artist_name, created_at, updated_at
	   FROM artwork_groups
	   WHERE id = ?
	   `

	var group models.ArtworkGroup
	err := db.conn.QueryRow(query, id).Scan(
		&group.ID,
		&group.Title,
		&group.Prompt,
		&group.Category,
		&group.OriginalURL,
		&group.ArtistName,
		&group.CreatedAt,
		&group.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("group not found")
		}
		return nil, fmt.Errorf("failed to get group: %w", err)
	}

	return &group, nil
}

// ListGroups retrieves all artwork groups
func (db *DB) ListGroups() ([]models.ArtworkGroup, error) {
	query := `
	       SELECT id, title, prompt, category, original_url, artist_name, created_at, updated_at
	       FROM artwork_groups
	       ORDER BY created_at ASC
	       `

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query groups: %w", err)
	}
	defer rows.Close()

	var groups []models.ArtworkGroup
	for rows.Next() {
		var group models.ArtworkGroup
		err := rows.Scan(
			&group.ID,
			&group.Title,
			&group.Prompt,
			&group.Category,
			&group.CreatedAt,
			&group.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan group: %w", err)
		}
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return groups, nil
}

// CreateArtwork creates a new artwork
func (db *DB) CreateArtwork(artwork models.Artwork) (int, error) {
	query := `
	INSERT INTO artworks (group_id, model, temperature, max_tokens, svg, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := db.conn.Exec(query, artwork.GroupID, artwork.Model, artwork.Temperature, artwork.MaxTokens, artwork.SVG, artwork.CreatedAt, artwork.UpdatedAt)
	if err != nil {
		return 0, fmt.Errorf("failed to create artwork: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return int(id), nil
}

// GetArtwork retrieves an artwork by ID
func (db *DB) GetArtwork(id int) (*models.Artwork, error) {
	query := `
	SELECT id, group_id, model, temperature, max_tokens, svg, created_at, updated_at
	FROM artworks
	WHERE id = ?
	`

	var artwork models.Artwork
	err := db.conn.QueryRow(query, id).Scan(
		&artwork.ID,
		&artwork.GroupID,
		&artwork.Model,
		&artwork.Temperature,
		&artwork.MaxTokens,
		&artwork.SVG,
		&artwork.CreatedAt,
		&artwork.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("artwork not found")
		}
		return nil, fmt.Errorf("failed to get artwork: %w", err)
	}

	return &artwork, nil
}

// ListArtworksByGroup retrieves all artworks for a group
func (db *DB) ListArtworksByGroup(groupID int) ([]models.Artwork, error) {
	query := `
	SELECT id, group_id, model, temperature, max_tokens, svg, created_at, updated_at
	FROM artworks
	WHERE group_id = ?
	ORDER BY model ASC
	`

	rows, err := db.conn.Query(query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to query artworks: %w", err)
	}
	defer rows.Close()

	var artworks []models.Artwork
	for rows.Next() {
		var artwork models.Artwork
		err := rows.Scan(
			&artwork.ID,
			&artwork.GroupID,
			&artwork.Model,
			&artwork.Temperature,
			&artwork.MaxTokens,
			&artwork.SVG,
			&artwork.CreatedAt,
			&artwork.UpdatedAt,
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

// Artwork parameters are stored in `temperature` and `max_tokens` columns.

// SaveArtworkSVG saves the SVG content for an artwork
func (db *DB) SaveArtworkSVG(id int, svg string) error {
	query := `
	UPDATE artworks
	SET svg = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`

	result, err := db.conn.Exec(query, svg, id)
	if err != nil {
		return fmt.Errorf("failed to save artwork SVG: %w", err)
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

// DeleteGroup deletes a group by ID (cascades to delete all associated artworks)
func (db *DB) DeleteGroup(id int) error {
	query := `DELETE FROM artwork_groups WHERE id = ?`

	result, err := db.conn.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("group with ID %d not found", id)
	}

	return nil
}

// UpdateArtwork updates temperature and max_tokens for an artwork
func (db *DB) UpdateArtwork(id int, temperature float64, maxTokens int) error {
	query := `
	UPDATE artworks
	SET temperature = ?, max_tokens = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`

	result, err := db.conn.Exec(query, temperature, maxTokens, id)
	if err != nil {
		return fmt.Errorf("failed to update artwork: %w", err)
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

// ListGroupsWithArtworks retrieves groups with their associated artworks
// If category is not empty, filters groups by category
func (db *DB) ListGroupsWithArtworks(category string) ([]models.ArtworkGroup, map[int][]models.Artwork, error) {
	// Build query with optional category filter
	query := `
		SELECT id, title, prompt, category, original_url, artist_name, created_at, updated_at
		FROM artwork_groups`

	var args []interface{}
	if category != "" {
		query += ` WHERE category = ?`
		args = append(args, category)
	}

	query += ` ORDER BY created_at ASC`

	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query groups: %w", err)
	}
	defer rows.Close()

	var groups []models.ArtworkGroup
	var groupIDs []int
	for rows.Next() {
		var group models.ArtworkGroup
		err := rows.Scan(
			&group.ID,
			&group.Title,
			&group.Prompt,
			&group.Category,
			&group.OriginalURL,
			&group.ArtistName,
			&group.CreatedAt,
			&group.UpdatedAt,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan group: %w", err)
		}
		groups = append(groups, group)
		groupIDs = append(groupIDs, group.ID)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error iterating group rows: %w", err)
	}

	// If no groups found, return empty results
	if len(groups) == 0 {
		return groups, make(map[int][]models.Artwork), nil
	}

	// Fetch all artworks for these groups in one query
	artworkMap := make(map[int][]models.Artwork)

	// Build placeholders for IN clause
	placeholders := ""
	for i := range groupIDs {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
	}

	artworkQuery := fmt.Sprintf(`
	SELECT id, group_id, model, temperature, max_tokens, svg, created_at, updated_at
	FROM artworks
	WHERE group_id IN (%s)
	ORDER BY group_id, model ASC
	`, placeholders)

	// Convert groupIDs to interface{} slice for query
	artworkArgs := make([]interface{}, len(groupIDs))
	for i, id := range groupIDs {
		artworkArgs[i] = id
	}

	artworkRows, err := db.conn.Query(artworkQuery, artworkArgs...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query artworks: %w", err)
	}
	defer artworkRows.Close()

	for artworkRows.Next() {
		var artwork models.Artwork
		err := artworkRows.Scan(
			&artwork.ID,
			&artwork.GroupID,
			&artwork.Model,
			&artwork.Temperature,
			&artwork.MaxTokens,
			&artwork.SVG,
			&artwork.CreatedAt,
			&artwork.UpdatedAt,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan artwork: %w", err)
		}
		artworkMap[artwork.GroupID] = append(artworkMap[artwork.GroupID], artwork)
	}

	if err := artworkRows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error iterating artwork rows: %w", err)
	}

	return groups, artworkMap, nil
}

// GetDistinctCategories returns all distinct categories from artwork groups
func (db *DB) GetDistinctCategories() ([]string, error) {
	query := `
	SELECT DISTINCT category
	FROM artwork_groups
	WHERE category != ''
	ORDER BY category
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
		return nil, fmt.Errorf("error iterating category rows: %w", err)
	}

	return categories, nil
}

// GetRandomGroupWithModelArtworks returns a random group that has artworks from both specified models
func (db *DB) GetRandomGroupWithModelArtworks(model1, model2 string) (*models.ArtworkGroup, []models.Artwork, error) {
	// First, find groups that have artworks from both models
	query := `
		SELECT DISTINCT g.id, g.title, g.prompt, g.category, g.original_url, g.artist_name, g.created_at, g.updated_at
		FROM artwork_groups g
		WHERE EXISTS (
			SELECT 1 FROM artworks a WHERE a.group_id = g.id AND a.model LIKE ?
		)
		AND EXISTS (
			SELECT 1 FROM artworks a WHERE a.group_id = g.id AND a.model LIKE ?
		)
		ORDER BY RANDOM()
		LIMIT 1
	`

	var group models.ArtworkGroup
	err := db.conn.QueryRow(query, "%"+model1+"%", "%"+model2+"%").Scan(
		&group.ID,
		&group.Title,
		&group.Prompt,
		&group.Category,
		&group.OriginalURL,
		&group.ArtistName,
		&group.CreatedAt,
		&group.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, fmt.Errorf("no group found with artworks from both models")
		}
		return nil, nil, fmt.Errorf("failed to get random group: %w", err)
	}

	// Get artworks for this group, filtered by the two models
	artworkQuery := `
		SELECT id, group_id, model, temperature, max_tokens, svg, created_at, updated_at
		FROM artworks
		WHERE group_id = ? AND (model LIKE ? OR model LIKE ?)
		ORDER BY CASE
			WHEN model LIKE ? THEN 1
			WHEN model LIKE ? THEN 2
			ELSE 3
		END
		`

	rows, err := db.conn.Query(artworkQuery, group.ID, "%"+model1+"%", "%"+model2+"%", "%"+model1+"%", "%"+model2+"%")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query artworks: %w", err)
	}
	defer rows.Close()

	var artworks []models.Artwork
	for rows.Next() {
		var artwork models.Artwork
		err := rows.Scan(
			&artwork.ID,
			&artwork.GroupID,
			&artwork.Model,
			&artwork.Temperature,
			&artwork.MaxTokens,
			&artwork.SVG,
			&artwork.CreatedAt,
			&artwork.UpdatedAt,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan artwork: %w", err)
		}
		artworks = append(artworks, artwork)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("error iterating artwork rows: %w", err)
	}

	return &group, artworks, nil
}
