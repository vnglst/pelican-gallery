package database

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"pelican-gallery/internal/models"
)

func TestNewMigratesExistingArtworksTable(t *testing.T) {
	path := filepath.Join(t.TempDir(), "legacy.db")
	legacy, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := legacy.Exec(`CREATE TABLE artworks (
		id INTEGER PRIMARY KEY, group_id INTEGER NOT NULL, model TEXT NOT NULL,
		temperature REAL NOT NULL DEFAULT 0, max_tokens INTEGER NOT NULL DEFAULT 0,
		svg TEXT DEFAULT '', featured BOOLEAN NOT NULL DEFAULT 0,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		t.Fatal(err)
	}
	if err := legacy.Close(); err != nil {
		t.Fatal(err)
	}

	db, err := New(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows, err := db.conn.Query("PRAGMA table_info(artworks)")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var cid, notNull, primaryKey int
		var name, columnType string
		var defaultValue interface{}
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &primaryKey); err != nil {
			t.Fatal(err)
		}
		if name == "model_created_at" {
			return
		}
	}
	t.Fatal("model_created_at column was not added")
}

func TestCreateArtworkPersistsModelCreatedAt(t *testing.T) {
	db, err := New(filepath.Join(t.TempDir(), "gallery.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	now := time.Now().UTC().Truncate(time.Second)
	groupID, err := db.CreateGroup(models.ArtworkGroup{
		Title: "Test", Prompt: "Test", CreatedAt: now, UpdatedAt: now,
	})
	if err != nil {
		t.Fatal(err)
	}

	const modelCreatedAt = int64(1_755_000_000)
	artworkID, err := db.CreateArtwork(models.Artwork{
		GroupID: groupID, Model: "example/model-1", ModelCreatedAt: modelCreatedAt,
		CreatedAt: now, UpdatedAt: now,
	})
	if err != nil {
		t.Fatal(err)
	}

	artwork, err := db.GetArtwork(artworkID)
	if err != nil {
		t.Fatal(err)
	}
	if artwork.ModelCreatedAt != modelCreatedAt {
		t.Fatalf("ModelCreatedAt = %d, want %d", artwork.ModelCreatedAt, modelCreatedAt)
	}
}
