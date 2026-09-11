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
	found := map[string]bool{}
	for rows.Next() {
		var cid, notNull, primaryKey int
		var name, columnType string
		var defaultValue interface{}
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &primaryKey); err != nil {
			t.Fatal(err)
		}
		found[name] = true
	}
	for _, column := range []string{"model_name", "model_created_at", "model_metadata_json"} {
		if !found[column] {
			t.Fatalf("%s column was not added", column)
		}
	}
}

func TestCreateArtworkPersistsModelMetadata(t *testing.T) {
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
		GroupID: groupID, Model: "example/model-1", ModelName: "Example Model 1", ModelCreatedAt: modelCreatedAt,
		ModelMetadata: `{"id":"example/model-1"}`, CreatedAt: now, UpdatedAt: now,
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
	if artwork.ModelName != "Example Model 1" {
		t.Fatalf("ModelName = %q, want %q", artwork.ModelName, "Example Model 1")
	}
	if artwork.ModelMetadata != `{"id":"example/model-1"}` {
		t.Fatalf("ModelMetadata = %q", artwork.ModelMetadata)
	}
}

func TestBackfillArtworkModelMetadataPreservesExistingValues(t *testing.T) {
	db, err := New(filepath.Join(t.TempDir(), "gallery.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	now := time.Now().UTC().Truncate(time.Second)
	groupID, err := db.CreateGroup(models.ArtworkGroup{Title: "Test", Prompt: "Test", CreatedAt: now, UpdatedAt: now})
	if err != nil {
		t.Fatal(err)
	}
	artworkID, err := db.CreateArtwork(models.Artwork{GroupID: groupID, Model: "example/model:free", ModelName: "Stored name", CreatedAt: now, UpdatedAt: now})
	if err != nil {
		t.Fatal(err)
	}

	updated, err := db.BackfillArtworkModelMetadata([]models.ModelInfo{{ID: "example/model", Name: "Live name", Created: 1234, MetadataJSON: `{"id":"example/model"}`}})
	if err != nil {
		t.Fatal(err)
	}
	if updated != 1 {
		t.Fatalf("updated = %d, want 1", updated)
	}
	artwork, err := db.GetArtwork(artworkID)
	if err != nil {
		t.Fatal(err)
	}
	if artwork.ModelName != "Stored name" || artwork.ModelCreatedAt != 1234 || artwork.ModelMetadata == "" {
		t.Fatalf("metadata = (%q, %d, %q)", artwork.ModelName, artwork.ModelCreatedAt, artwork.ModelMetadata)
	}
}

func TestSaveArtworkGenerationPersistsUsageHistory(t *testing.T) {
	db, err := New(filepath.Join(t.TempDir(), "gallery.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	now := time.Now().UTC().Truncate(time.Second)
	groupID, err := db.CreateGroup(models.ArtworkGroup{Title: "Test", Prompt: "Test", CreatedAt: now, UpdatedAt: now})
	if err != nil {
		t.Fatal(err)
	}
	artworkID, err := db.CreateArtwork(models.Artwork{GroupID: groupID, Model: "example/model", CreatedAt: now, UpdatedAt: now})
	if err != nil {
		t.Fatal(err)
	}

	usage := models.GenerationUsage{
		PromptTokens: 120, CompletionTokens: 340, TotalTokens: 460,
		ReasoningTokens: 25, CachedTokens: 10, CostUSD: 0.0123,
		RawJSON: `{"prompt_tokens":120,"completion_tokens":340,"total_tokens":460,"cost":0.0123}`,
	}
	if err := db.SaveArtworkGeneration(artworkID, "<svg></svg>", usage); err != nil {
		t.Fatal(err)
	}

	var prompt, completion, total, reasoning, cached int
	var cost float64
	var raw string
	err = db.conn.QueryRow(`
		SELECT prompt_tokens, completion_tokens, total_tokens, reasoning_tokens, cached_tokens, cost_usd, usage_json
		FROM artwork_generations WHERE artwork_id = ?
	`, artworkID).Scan(&prompt, &completion, &total, &reasoning, &cached, &cost, &raw)
	if err != nil {
		t.Fatal(err)
	}
	if prompt != 120 || completion != 340 || total != 460 || reasoning != 25 || cached != 10 || cost != 0.0123 || raw != usage.RawJSON {
		t.Fatalf("stored usage = (%d, %d, %d, %d, %d, %f, %q)", prompt, completion, total, reasoning, cached, cost, raw)
	}
	artwork, err := db.GetArtwork(artworkID)
	if err != nil {
		t.Fatal(err)
	}
	if artwork.SVG != "<svg></svg>" {
		t.Fatalf("SVG = %q", artwork.SVG)
	}
	if !artwork.HasGenerationCost || artwork.GenerationCostUSD != usage.CostUSD {
		t.Fatalf("displayed generation cost = (%t, %f), want (true, %f)", artwork.HasGenerationCost, artwork.GenerationCostUSD, usage.CostUSD)
	}

	latestUsage := models.GenerationUsage{CostUSD: 0.0042}
	if err := db.SaveArtworkGeneration(artworkID, "<svg>new</svg>", latestUsage); err != nil {
		t.Fatal(err)
	}
	artwork, err = db.GetArtwork(artworkID)
	if err != nil {
		t.Fatal(err)
	}
	if artwork.GenerationCostUSD != latestUsage.CostUSD {
		t.Fatalf("latest generation cost = %f, want %f", artwork.GenerationCostUSD, latestUsage.CostUSD)
	}
}
