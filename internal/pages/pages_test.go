package pages

import (
	"sort"
	"testing"
	"time"
)

func TestModelVersionLess(t *testing.T) {
	models := []string{
		"google/gemini-3-pro-preview",
		"google/gemini-2.5-pro",
		"google/gemini-2.0-flash-001",
		"google/gemini-3.1-pro-preview",
	}

	sort.Slice(models, func(i, j int) bool {
		return modelVersionLess(models[i], models[j])
	})

	want := []string{
		"google/gemini-2.0-flash-001",
		"google/gemini-2.5-pro",
		"google/gemini-3-pro-preview",
		"google/gemini-3.1-pro-preview",
	}
	for i := range want {
		if models[i] != want[i] {
			t.Fatalf("models[%d] = %q, want %q", i, models[i], want[i])
		}
	}
}

func TestModelIDDisplayName(t *testing.T) {
	tests := map[string]string{
		"gemini-3-pro-preview": "Gemini 3 Pro Preview",
		"claude-3.7-sonnet":    "Claude 3.7 Sonnet",
		"gpt-oss-120b:free":    "GPT Oss 120b",
	}
	for modelID, want := range tests {
		if got := modelIDDisplayName(modelID); got != want {
			t.Errorf("modelIDDisplayName(%q) = %q, want %q", modelID, got, want)
		}
	}
}

func TestModelCapabilityRankPutsLesserModelsFirst(t *testing.T) {
	models := []string{"openai/gpt-5", "openai/gpt-5-mini", "openai/gpt-5-nano"}
	sort.SliceStable(models, func(i, j int) bool {
		return modelCapabilityRank(models[i]) < modelCapabilityRank(models[j])
	})

	want := []string{"openai/gpt-5-nano", "openai/gpt-5-mini", "openai/gpt-5"}
	for i := range want {
		if models[i] != want[i] {
			t.Fatalf("models[%d] = %q, want %q", i, models[i], want[i])
		}
	}
}

func TestFormatGenerationCost(t *testing.T) {
	tests := map[float64]string{
		0:         "Free",
		0.000001: "<$0.00001",
		0.000051: "$0.00005",
		0.000516: "$0.0005",
		0.006517: "$0.0065",
		0.0123:   "$0.01",
	}
	for cost, want := range tests {
		if got := formatGenerationCost(cost); got != want {
			t.Errorf("formatGenerationCost(%f) = %q, want %q", cost, got, want)
		}
	}
}

func TestAnthropicChronologyUsesOfficialDatesBeforeStoredOpenRouterDates(t *testing.T) {
	type artwork struct {
		model                string
		storedOpenRouterDate int64
	}

	artworks := []artwork{
		{model: "anthropic/claude-sonnet-4", storedOpenRouterDate: 1747930371},
		{model: "anthropic/claude-3.5-haiku"},
		{model: "anthropic/claude-3.7-sonnet"},
	}
	sort.SliceStable(artworks, func(i, j int) bool {
		leftDate, leftKnown := modelReleaseDate(artworks[i].model, artworks[i].storedOpenRouterDate)
		rightDate, rightKnown := modelReleaseDate(artworks[j].model, artworks[j].storedOpenRouterDate)
		if leftKnown && rightKnown && !leftDate.Equal(rightDate) {
			return leftDate.Before(rightDate)
		}
		if leftKnown != rightKnown {
			return leftKnown
		}
		return modelVersionLess(artworks[i].model, artworks[j].model)
	})

	want := []string{
		"anthropic/claude-3.5-haiku",
		"anthropic/claude-3.7-sonnet",
		"anthropic/claude-sonnet-4",
	}
	for i := range want {
		if artworks[i].model != want[i] {
			t.Fatalf("artworks[%d] = %q, want %q", i, artworks[i].model, want[i])
		}
	}

	releasedAt, ok := modelReleaseDate("anthropic/claude-3.7-sonnet", 0)
	if !ok || releasedAt.Year() != 2025 || releasedAt.Month() != time.February || releasedAt.Day() != 24 {
		t.Fatalf("Claude 3.7 release date = %v, %v; want 2025-02-24", releasedAt, ok)
	}
}
